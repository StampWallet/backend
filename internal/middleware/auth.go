package middleware

import (
    "log"
    "github.com/gin-gonic/gin"
    "github.com/StampWallet/backend/internal/services"
)

type AuthMiddleware struct {
    logger *log.Logger
    tokenService *services.TokenService
}

func (middleware *AuthMiddleware) Handle(c *gin.Context) {

}
