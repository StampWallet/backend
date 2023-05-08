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
var EmailExists = errors.New("Email exists")
var UnknownError = errors.New("Unknown error")
var InvalidToken = errors.New("Invalid token")

type AuthManager interface {
	Create(userDetails UserDetails) (*User, *Token, error)
	Login(email string, password string) (*User, *Token, error)
	Logout(tokenId string, tokenSecret string) (*User, *Token, error)
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
	tx := manager.baseServices.Database.First(&existingUser, &User{
		Email: userDetails.Email,
	})
	err = tx.GetError()
	if err == nil {
		return nil, nil, "", EmailExists
	} else if err != nil && err != gorm.ErrRecordNotFound {
		return nil, nil, "", err
	}

	hash, bcryptErr := bcrypt.GenerateFromPassword([]byte(userDetails.Password), 10)
	if bcryptErr != nil {
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
	tx = manager.baseServices.Database.Create(&user)
	if err := tx.GetError(); err != nil {
		return nil, nil, "", err
	}
	tx = manager.baseServices.Database.Commit()
	if err := tx.GetError(); err != nil {
		return nil, nil, "", err
	}

	emailToken, emailSecret, err := manager.tokenService.Create(user, TokenPurposeEmail, time.Now().Add(24*time.Hour))
	if err != nil {
		return nil, nil, "", err
	}
	sessionToken, sessionSecret, err := manager.tokenService.Create(user, TokenPurposeSession, time.Now().Add(time.Hour))
	if err != nil {
		return nil, nil, "", err
	}
	manager.emailService.Send(userDetails.Email, "test", "test "+emailToken.TokenId+":"+emailSecret)
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

func (manager *AuthManagerImpl) Logout(tokenId string, tokenSecret string) (*User, *Token, error) {
	return nil, nil, nil
}

func (manager *AuthManagerImpl) ConfirmEmail(tokenId string, tokensECRET string) (*User, error) {
	return nil, nil
}

func (manager *AuthManagerImpl) ChangePassword(user User, oldPassword string, newPassword string) (*User, error) {
	return nil, nil
}

func (manager *AuthManagerImpl) ChangeEmail(user User, newEmail string) (*User, error) {
	return nil, nil
}
