package managers

import (
    . "github.com/StampWallet/backend/internal/database"
    . "github.com/StampWallet/backend/internal/services"
)

type AuthManager interface {
    Create(userDetails UserDetails) (User, error)
    Login(email string, password string) (User, Token, error)
    Logout(tokenId string, token string) (User, Token, error)
    ConfirmEmail(user User, tokenId string, token string) error
    ChangePassword(user User, oldPassword string, newPassword string) error
    ChangeEmail(user User, newEmail string) (User, error)
}

type UserDetails struct {
    Email string
    Password string
}

type AuthManagerImpl struct {
    baseServices *BaseServices
    emailService *EmailService
    tokenService *TokenService
}

func (manager *AuthManagerImpl) Create(userDetails UserDetails) (*User, error) {
    return nil, nil
}              
               
func (manager *AuthManagerImpl) Login(email string, password string) (*User, *Token, error) {
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
