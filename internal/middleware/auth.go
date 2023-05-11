package middleware

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"

	api "github.com/StampWallet/backend/internal/api/models"
	services "github.com/StampWallet/backend/internal/services"
)

type AuthMiddleware struct {
	logger       *log.Logger
	tokenService services.TokenService
}

func (middleware *AuthMiddleware) Handle(c *gin.Context) {
	auth_header := c.GetHeader("Authorization")
	header_value_split := strings.Split(auth_header, " ")
	if len(header_value_split) != 2 || header_value_split[0] != "Bearer" {
		c.JSON(401, api.DefaultResponse{
			Status: api.UNAUTHORIZED,
		})
		return
	}

	token_split := strings.Split(header_value_split[1], ":")
	if len(token_split) != 2 {
		c.JSON(401, api.DefaultResponse{
			Status: api.UNAUTHORIZED,
		})
		return
	}

	user, token, err := middleware.tokenService.Check(token_split[0], token_split[1])
	if err == services.UnknownToken {
		c.JSON(401, api.DefaultResponse{
			Status: api.UNAUTHORIZED,
		})
		return
	} else if err != nil || user == nil || token == nil {
		c.JSON(500, api.DefaultResponse{
			Status: api.UNKNOWN_ERROR,
		})
		middleware.logger.Printf("Error: in AuthMiddleware.Handle, middleware.tokenService.Check: %s", err)
		return
	}

	c.Set("user", user)
	c.Set("token", token)
}
