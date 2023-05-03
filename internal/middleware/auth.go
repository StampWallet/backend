package middleware

import (
	"log"

	"github.com/StampWallet/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	logger       *log.Logger
	tokenService services.TokenService
}

func (middleware *AuthMiddleware) Handle(c *gin.Context) {

}
