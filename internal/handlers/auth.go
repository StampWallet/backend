package handlers

import (
    "log"
    "github.com/StampWallet/backend/internal/managers"
)

type AuthHandlers struct {
    authManager AuthManager
    logger log.Logger
}

func (handler *AuthHandlers) postAccount(c *gin.Context) {
               
}              
               
func (handler *AuthHandlers) postAccountEmail(c *gin.Context) {
               
}              
               
func (handler *AuthHandlers) postAccountPassword(c *gin.Context) {
               
}              
               
func (handler *AuthHandlers) postAccountEmailConfirmation(c *gin.Context) { 
               
}              
               
func (handler *AuthHandlers) postSession(c *gin.Context) {
               
}              
               
func (handler *AuthHandlers) deleteSession(c *gin.Context) {
               
}              
               
func (handler *AuthHandlers) Connect(rg *gin.RouterGroup, authMiddleware AuthMiddleware) {

}
