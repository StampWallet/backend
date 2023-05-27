package api

import (
	"log"

	"github.com/gin-gonic/gin"

	api "github.com/StampWallet/backend/internal/api/models"
	"github.com/StampWallet/backend/internal/database"
)

func getUserFromContext(logger *log.Logger, c *gin.Context) *database.User {
	userAny, exists := c.Get("user")
	if !exists {
		log.Printf("user not available context")
		c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
		return nil
	}
	return userAny.(*database.User)
}
