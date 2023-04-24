package middleware

import (
    "log"
    "github.com/gin-gonic/gin"
    "github.com/StampWallet/backend/internal/services"
)

type RequireValidEmailMiddleware struct {
    logger *log.Logger
    tokenService *services.TokenService
}

func (middleware *RequireValidEmailMiddleware) Handle(c *gin.Context) {

}
