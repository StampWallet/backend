package managers

import (
	"errors"
	"fmt"
	"net/mail"
	"time"

	"github.com/lithammer/shortuuid/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	. "github.com/StampWallet/backend/internal/database"
	. "github.com/StampWallet/backend/internal/services"
	. "github.com/StampWallet/backend/internal/utils"
)

var (
	ErrInvalidEmail        = errors.New("Invalid email address")
	ErrInvalidLogin        = errors.New("Invalid login") // Invalid email/password
	ErrInvalidOldPassword  = errors.New("Invalid old password")
	ErrInvalidTokenPurpose = errors.New("Invalid token purpose") // Invalid token purpose
	ErrEmailExists         = errors.New("Email exists")          // Another user has the same email
	ErrInvalidToken        = errors.New("Invalid token")         // Invalid/unknown token
	ErrUnknownError        = errors.New("Unknown error")         // Unexpected error returned by external services
)

type AuthManager interface {
	// Creates a new user account from UserDetails struct. Returns a database object and, session token and
	// matching token secret for the user. Sends a verification email.
	Create(userDetails UserDetails) (*User, *Token, string, error)

	// Checks if email and password match any user. If yes, returns the database object and serssion token
	// for that user.
	Login(email string, password string) (*User, *Token, string, error)

	// Checks if token id and secret match any session token. If yes, invalidates the token.
	Logout(tokenId string, tokenSecret string) (*User, *Token, error)

	// Checks if token id and secret match any email token. If yes, invalidates the token and changes
	// EmailVerified in user's database object to true.
	ConfirmEmail(tokenId string, tokenSecret string) (*User, error)

	// Changes password of user, if oldPassword matches user.PasswordHash.
	ChangePassword(user *User, oldPassword string, newPassword string) (*User, error)

	// Changes email of user, if no other user has the same email. Changes user.EmailVerified to false,
	// sends a new verification email.
	ChangeEmail(user *User, newEmail string) (*User, error)
}

type UserDetails struct {
	//FirstName string
	//LastName  string
	Email    string
	Password string
}

type AuthManagerImpl struct {
	baseServices BaseServices
	emailService EmailService
	tokenService TokenService
}

func CreateAuthManagerImpl(baseServices BaseServices,
	emailService EmailService, tokenService TokenService) *AuthManagerImpl {
	return &AuthManagerImpl{
		baseServices: baseServices,
		emailService: emailService,
		tokenService: tokenService,
	}
}

func (manager *AuthManagerImpl) Create(userDetails UserDetails) (*User, *Token, string, error) {
	_, err := mail.ParseAddress(userDetails.Email)
	if err != nil {
		return nil, nil, "", ErrInvalidEmail
	}

	var existingUser User
	tx := manager.baseServices.Database.Begin()

	// Check if email is valid
	_, emailErr := mail.ParseAddress(userDetails.Email)
	if emailErr != nil {
		return nil, nil, "", ErrInvalidEmail
	}

	// Check if another user has this email
	txFirst := tx.First(&existingUser, &User{
		Email: userDetails.Email,
	})
	err = txFirst.GetError()
	if err == nil {
		tx.Rollback()
		return nil, nil, "", ErrEmailExists
	} else if err != nil && err != gorm.ErrRecordNotFound {
		tx.Rollback()
		return nil, nil, "", fmt.Errorf("failed to to find user, database error: %+v", err)
	}

	// Hash the password
	hash, bcryptErr := bcrypt.GenerateFromPassword([]byte(userDetails.Password), 10)
	if bcryptErr != nil {
		tx.Rollback()
		return nil, nil, "", bcryptErr
	}

	// Create the user in the database
	user := User{
		PublicId:     shortuuid.New(),
		Email:        userDetails.Email,
		PasswordHash: string(hash),
		//FirstName:     userDetails.FirstName,
		//LastName:      userDetails.LastName,
		EmailVerified: false,
	}
	txCreate := tx.Create(&user)
	if err := txCreate.GetError(); err != nil {
		tx.Rollback()
		return nil, nil, "", fmt.Errorf("%s failed to to create user, database error: %+v", CallerFilename(), err)
	}

	// Create token for email verification
	emailToken, emailSecret, err := manager.tokenService.Create(&user, TokenPurposeEmail, time.Now().Add(24*time.Hour))
	if err != nil {
		tx.Rollback()
		return nil, nil, "", fmt.Errorf("%s failed to to create email token, tokenservice error: %+v", CallerFilename(), err)
	}

	// Create token for session
	sessionToken, sessionSecret, err := manager.tokenService.Create(&user, TokenPurposeSession, time.Now().Add(time.Hour))
	if err != nil {
		tx.Rollback()
		return nil, nil, "", fmt.Errorf("%s failed to to create session token, tokenservice error: %+v", CallerFilename(), err)
	}

	// Send email verification token
	mailErr := manager.emailService.Send(userDetails.Email, "test", "test "+emailToken.TokenId+":"+emailSecret)
	if mailErr != nil {
		tx.Rollback()
		return nil, nil, "", fmt.Errorf("%s failed to to create session token, tokenservice error: %+v", CallerFilename(), err)
	}

	// Commit transaction
	if err := tx.Commit().GetError(); err != nil {
		return nil, nil, "", fmt.Errorf("%s failed to commit, database error: %+v", CallerFilename(), err)
	}
	return &user, sessionToken, sessionSecret, nil
}

func (manager *AuthManagerImpl) Login(email string, password string) (*User, *Token, string, error) {
	var user User

	// Find user
	tx := manager.baseServices.Database.First(&user, User{Email: email})
	err := tx.GetError()
	if err == gorm.ErrRecordNotFound {
		return nil, nil, "", ErrInvalidLogin
	} else if err != nil {
		return nil, nil, "", err
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return nil, nil, "", ErrInvalidLogin
	} else if err != nil {
		return nil, nil, "", err
	}

	// Create session token
	sessionToken, sessionSecret, err := manager.tokenService.Create(&user, TokenPurposeSession, time.Now().Add(time.Hour))
	if err != nil {
		return nil, nil, "", err
	}

	return &user, sessionToken, sessionSecret, nil
}

func (manager *AuthManagerImpl) Logout(tokenId string, tokenSecret string) (*User, *Token, error) {
	// Find token
	token, err := manager.tokenService.Check(tokenId, tokenSecret)
	if err != nil {
		manager.baseServices.Logger.Printf("tokenService.Check returned an error: %s", err)
		return nil, nil, ErrInvalidToken
	}
	if token.TokenPurpose != TokenPurposeSession {
		return nil, nil, ErrInvalidTokenPurpose
	}

	// Invalidate token
	invalidatedToken, err := manager.tokenService.Invalidate(token)
	if err != nil {
		manager.baseServices.Logger.Printf("tokenService.Invalidate returned an error: %s", err)
		return nil, nil, ErrInvalidToken
	}
	return invalidatedToken.User, invalidatedToken, nil
}

func (manager *AuthManagerImpl) ConfirmEmail(tokenId string, tokenSecret string) (*User, error) {
	tx := manager.baseServices.Database.Begin()

	// Find email token
	token, err := manager.tokenService.Check(tokenId, tokenSecret)
	if err == ErrUnknownToken {
		tx.Rollback()
		return nil, ErrInvalidToken
	} else if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Check if token is valid
	if token.TokenPurpose != TokenPurposeEmail {
		tx.Rollback()
		return nil, ErrInvalidTokenPurpose
	}
	if token.Recalled || token.Used {
		tx.Rollback()
		return nil, ErrInvalidToken
	}

	// Change email verification status
	token.User.EmailVerified = true
	txSave := tx.Save(token.User)
	if err = txSave.GetError(); err != nil {
		tx.Rollback()
		return nil, err
	}

	// Invalidate token
	token, err = manager.tokenService.Invalidate(token)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit().GetError(); err != nil {
		// TODO make sure that user cannot get locked up here if transaction fails and token is invalidated
		return nil, err
	}
	return token.User, nil
}

func (manager *AuthManagerImpl) ChangePassword(user *User, oldPassword string, newPassword string) (*User, error) {
	// Check if old password matches
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return nil, ErrInvalidOldPassword
	} else if err != nil {
		return nil, err
	}

	// Generate new password
	newHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), 10)
	if err != nil {
		return nil, err
	}

	// Update password in the database
	user.PasswordHash = string(newHash)
	tx := manager.baseServices.Database.Save(user)
	if err := tx.GetError(); err != nil {
		return nil, err
	}

	return user, nil
}

func (manager *AuthManagerImpl) ChangeEmail(user *User, newEmail string) (*User, error) {
	// Check if email is valid
	_, err := mail.ParseAddress(newEmail)
	if err != nil {
		return nil, ErrInvalidEmail
	}

	// Update user in the db
	user.Email = newEmail
	user.EmailVerified = false

	tx := manager.baseServices.Database.Begin()
	txSave := tx.Save(&user)
	dbErr := txSave.GetError()
	// Another user has the same email
	if dbErr == gorm.ErrDuplicatedKey {
		tx.Rollback()
		manager.baseServices.Logger.Printf("gorm error when changing email %s\n", dbErr)
		return nil, ErrEmailExists
	} else if dbErr != nil {
		return nil, dbErr
	}

	// Create email verification token
	emailToken, emailSecret, err := manager.tokenService.Create(user,
		TokenPurposeEmail, time.Now().Add(24*time.Hour))
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Send email verification token
	mailError := manager.emailService.Send(user.Email, "test", "test "+emailToken.TokenId+":"+emailSecret)
	if mailError != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit().GetError(); err != nil {
		return nil, err
	}

	return user, nil
}
