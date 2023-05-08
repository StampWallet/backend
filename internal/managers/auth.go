package managers

import (
	"errors"
	"net/mail"
	"time"

	. "github.com/StampWallet/backend/internal/database"
	. "github.com/StampWallet/backend/internal/services"
	"github.com/lithammer/shortuuid/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var InvalidEmail = errors.New("Invalid email")
var InvalidLogin = errors.New("Invalid login")
var InvalidOldPassword = errors.New("Invalid old password")
var InvalidTokenPurpose = errors.New("Invalid login")
var EmailExists = errors.New("Email exists")
var UnknownError = errors.New("Unknown error")
var InvalidToken = errors.New("Invalid token")

type AuthManager interface {
	Create(userDetails UserDetails) (*User, *Token, error)
	Login(email string, password string) (*User, *Token, error)
	ConfirmEmail(tokenId string, tokenSecret string) (*User, error)
	ChangePassword(user User, oldPassword string, newPassword string) (*User, error)
	ChangeEmail(user User, newEmail string) (*User, error)
}

type UserDetails struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
}

type AuthManagerImpl struct {
	baseServices *BaseServices
	emailService EmailService
	tokenService TokenService
}

func (manager *AuthManagerImpl) Create(userDetails UserDetails) (*User, *Token, string, error) {
	_, err := mail.ParseAddress(userDetails.Email)
	if err != nil {
		return nil, nil, "", InvalidEmail
	}

	var existingUser User
	tx := manager.baseServices.Database.Begin()
	tx = tx.First(&existingUser, &User{
		Email: userDetails.Email,
	})
	err = tx.GetError()
	if err == nil {
		tx.Rollback()
		return nil, nil, "", EmailExists
	} else if err != nil && err != gorm.ErrRecordNotFound {
		tx.Rollback()
		return nil, nil, "", err
	}

	hash, bcryptErr := bcrypt.GenerateFromPassword([]byte(userDetails.Password), 10)
	if bcryptErr != nil {
		tx.Rollback()
		return nil, nil, "", bcryptErr
	}
	user := User{
		PublicId:      shortuuid.New(),
		Email:         userDetails.Email,
		PasswordHash:  string(hash),
		FirstName:     userDetails.FirstName,
		LastName:      userDetails.LastName,
		EmailVerified: false,
	}
	tx = tx.Create(&user)
	if err := tx.GetError(); err != nil {
		tx.Rollback()
		return nil, nil, "", err
	}

	emailToken, emailSecret, err := manager.tokenService.Create(user, TokenPurposeEmail, time.Now().Add(24*time.Hour))
	if err != nil {
		tx.Rollback()
		return nil, nil, "", err
	}
	sessionToken, sessionSecret, err := manager.tokenService.Create(user, TokenPurposeSession, time.Now().Add(time.Hour))
	if err != nil {
		tx.Rollback()
		return nil, nil, "", err
	}
	mailErr := manager.emailService.Send(userDetails.Email, "test", "test "+emailToken.TokenId+":"+emailSecret)
	if mailErr != nil {
		tx.Rollback()
		return nil, nil, "", err
	}

	if err := tx.Commit().GetError(); err != nil {
		return nil, nil, "", err
	}
	return &user, sessionToken, sessionSecret, nil
}

func (manager *AuthManagerImpl) Login(email string, password string) (*User, *Token, string, error) {
	var user User
	tx := manager.baseServices.Database.First(&user, User{Email: email})
	err := tx.GetError()
	if err == gorm.ErrRecordNotFound {
		return nil, nil, "", InvalidLogin
	} else if err != nil {
		return nil, nil, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return nil, nil, "", InvalidLogin
	} else if err != nil {
		return nil, nil, "", err
	}

	sessionToken, sessionSecret, err := manager.tokenService.Create(user, TokenPurposeSession, time.Now().Add(time.Hour))
	if err != nil {
		return nil, nil, "", err
	}

	return &user, sessionToken, sessionSecret, nil
}

func (manager *AuthManagerImpl) ConfirmEmail(tokenId string, tokenSecret string) (*User, error) {
	tx := manager.baseServices.Database.Begin()

	user, token, err := manager.tokenService.Check(tokenId, tokenSecret)
	if err == UnknownToken {
		tx.Rollback()
		return nil, InvalidToken
	} else if err != nil {
		tx.Rollback()
		return nil, err
	}

	if token.TokenPurpose != TokenPurposeEmail {
		tx.Rollback()
		return nil, InvalidTokenPurpose
	}
	if token.Recalled || token.Used {
		tx.Rollback()
		return nil, InvalidToken
	}

	user.EmailVerified = true
	tx = tx.Save(user)
	if err = tx.GetError(); err != nil {
		tx.Rollback()
		return nil, err
	}

	user, token, err = manager.tokenService.Invalidate(*token)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().GetError(); err != nil {
		return nil, err
	}
	return user, nil
}

func (manager *AuthManagerImpl) ChangePassword(user User, oldPassword string, newPassword string) (*User, error) {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return nil, InvalidOldPassword
	} else if err != nil {
		return nil, err
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), 10)
	if err != nil {
		return nil, err
	}

	user.PasswordHash = string(newHash)
	tx := manager.baseServices.Database.Save(&user)
	if err := tx.GetError(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (manager *AuthManagerImpl) ChangeEmail(user User, newEmail string) (*User, error) {
	_, err := mail.ParseAddress(newEmail)
	if err != nil {
		return nil, InvalidEmail
	}

	user.Email = newEmail
	user.EmailVerified = false

	tx := manager.baseServices.Database.Begin()
	tx = tx.Save(&user)
	dbErr := tx.GetError()
	if dbErr == gorm.ErrDuplicatedKey {
		tx.Rollback()
		manager.baseServices.Logger.Printf("gorm error when changing email %s\n", dbErr)
		return nil, EmailExists
	} else if dbErr != nil {
		return nil, dbErr
	}

	emailToken, emailSecret, err := manager.tokenService.Create(user,
		TokenPurposeEmail, time.Now().Add(24*time.Hour))
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	mailError := manager.emailService.Send(user.Email, "test", "test "+emailToken.TokenId+":"+emailSecret)
	if mailError != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().GetError(); err != nil {
		return nil, err
	}

	return &user, nil
}
