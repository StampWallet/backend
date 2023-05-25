package middleware

import (
	"log"

	api "github.com/StampWallet/backend/internal/api/models"
	"github.com/StampWallet/backend/internal/database"
	"github.com/gin-gonic/gin"
)

type RequireValidEmailMiddleware struct {
	logger *log.Logger
}

func CreateRequireValidEmailMiddleware(logger *log.Logger) *RequireValidEmailMiddleware {
	return &RequireValidEmailMiddleware{
		logger: logger,
	}
}

func (middleware *RequireValidEmailMiddleware) Handle(c *gin.Context) {
	userAny, exists := c.Get("user")
	if !exists {
		middleware.logger.Printf("user not available in RequireValidEmailMiddleware context")
		c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
		return
	}
	user := userAny.(*database.User)
	if user.EmailVerified {
		c.Next()
	} else {
		c.JSON(403, api.DefaultResponse{Status: api.FORBIDDEN, Message: "EMAIL_NOT_VERIFIED"})
		c.Abort()
	}
}
