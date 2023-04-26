package middleware

import (
	"github.com/StampWallet/backend/internal/services"
	"github.com/gin-gonic/gin"
	"log"
)

type RequireValidEmailMiddleware struct {
	logger       *log.Logger
	tokenService *services.TokenService
}

func (middleware *RequireValidEmailMiddleware) Handle(c *gin.Context) {

}
