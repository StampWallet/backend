package middleware

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"

	api "github.com/StampWallet/backend/internal/api/models"
	services "github.com/StampWallet/backend/internal/services"
)

// Middleware that checks if the requests contains a valid (not expired, not recalled) session token.
// Session token is expected to be in the Authorization HTTP header. Expected HTTP authorization scheme is "Bearer".
// Token is expected to be in the following format: {{.TokenId}}:{{.TokenSecret}}
// If the token is valid, database.User object of the token owner is inserted into the context under "user" key and
// the request is passed to the next handler.
// If the token is not valid, the middleware returns 401 Unauthorized and the request is not passed to the next handler.
type AuthMiddleware struct {
	logger       *log.Logger
	tokenService services.TokenService
}

func CreateAuthMiddleware(logger *log.Logger, tokenService services.TokenService) *AuthMiddleware {
	return &AuthMiddleware{
		logger:       logger,
		tokenService: tokenService,
	}
}

// Gin handler function for the middleware
func (middleware *AuthMiddleware) Handle(c *gin.Context) {
	// Parse Authorization header value - divide on spaces
	auth_header := c.GetHeader("Authorization")
	header_value_split := strings.Split(auth_header, " ")

	// Make sure that Authorization scheme is Bearer and that the headers contains a single token
	if len(header_value_split) != 2 || header_value_split[0] != "Bearer" {
		c.AbortWithStatusJSON(401, api.DefaultResponse{
			Status: api.UNAUTHORIZED,
		})
		return
	}

	// Split the token by :
	token_split := strings.Split(header_value_split[1], ":")
	if len(token_split) != 2 {
		c.AbortWithStatusJSON(401, api.DefaultResponse{
			Status: api.UNAUTHORIZED,
		})
		return
	}

	// Call TokenService to make sure that the token is valid
	token, err := middleware.tokenService.Check(token_split[0], token_split[1])
	if err == services.ErrUnknownToken || err == services.ErrTokenExpired || err == services.ErrTokenUsed {
		c.AbortWithStatusJSON(401, api.DefaultResponse{
			Status: api.UNAUTHORIZED,
		})
		return
	} else if err != nil || token.User == nil || token == nil {
		c.AbortWithStatusJSON(500, api.DefaultResponse{
			Status: api.UNKNOWN_ERROR,
		})
		middleware.logger.Printf("Error: in AuthMiddleware.Handle, middleware.tokenService.Check: %s", err)
		return
	}

	c.Set("user", token.User)
	c.Set("token", token)
	c.Next()
}
