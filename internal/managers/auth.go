package managers

import (
	"errors"
	//"time"

	. "github.com/StampWallet/backend/internal/database"
	. "github.com/StampWallet/backend/internal/services"
	//"github.com/lithammer/shortuuid/v4"
	//"golang.org/x/crypto/bcrypt"
)

var InvalidEmail = errors.New("Invalid email")
var InvalidLogin = errors.New("Invalid login")
var EmailExists = errors.New("Email exists")
var UnknownError = errors.New("Email exists")

type AuthManager interface {
    Create(userDetails UserDetails) (User, Token, error)
    Login(email string, password string) (User, Token, error)
    Logout(tokenId string, token string) (User, Token, error)
    ConfirmEmail(user User, tokenId string, token string) error
    ChangePassword(user User, oldPassword string, newPassword string) error
    ChangeEmail(user User, newEmail string) (User, error)
}

type UserDetails struct {
    FirstName string
    LastName string
    Email string
    Password string
}

type AuthManagerImpl struct {
    baseServices *BaseServices
    emailService EmailService
    tokenService TokenService
}

func (manager *AuthManagerImpl) Create(userDetails UserDetails) (*User, *Token, error) {
    //var existingUser *User
    //manager.baseServices.Database.Find(&existingUser, &User{ 
    //    Email: userDetails.Email,
    //})
    //if existingUser != nil {
    //    return nil, nil, EmailExists
    //}
    //hash, bcryptErr := bcrypt.GenerateFromPassword([]byte(userDetails.Password), 10)
    //if bcryptErr != nil {
    //    return nil, nil, bcryptErr
    //}
    //user := User{ 
    //    PublicId: shortuuid.New(),
    //    Email: userDetails.Email, 
    //    PasswordHash: string(hash),
    //    FirstName: userDetails.FirstName,
    //    LastName: userDetails.LastName,
    //    EmailVerified: false,
    //}
    //manager.baseServices.Database.Create(&user)
    //manager.baseServices.Database.Commit()
    //_, err := manager.tokenService.Create(user, EmailTokenPurpose, time.Now().Add(24*time.Hour))
    //if err != nil {
    //    return nil, nil, err
    //}
    //sessionToken, err := manager.tokenService.Create(user, SessionTokenPurpose, time.Now().Add(time.Hour))
    //if err != nil {
    //    return nil, nil, err
    //}
    //manager.emailService.Send(userDetails.Email, "test", "test")
    //return &user, &sessionToken, nil
    return nil, nil, nil
}              
               
func (manager *AuthManagerImpl) Login(email string, password string) (*User, *Token, error) {
    var user User
    manager.baseServices.Database.First(&user, User{ Email: email })
    return nil, nil, nil
}              
               
func (manager *AuthManagerImpl) Logout(tokenId string, token string) (*User, *Token, error) {
    return nil, nil, nil
}              
               
func (manager *AuthManagerImpl) ConfirmEmail(user User, tokenId string, token string) error {
    return nil
}              
               
func (manager *AuthManagerImpl) ChangePassword(user User, oldPassword string, newPassword string) error {
    return nil
}              
               
func (manager *AuthManagerImpl) ChangeEmail(user User, newEmail string) (*User, error) {
    return nil, nil
}
