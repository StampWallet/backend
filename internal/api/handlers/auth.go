package api

import (
	"errors"
	"log"
	"strings"

	"github.com/gin-gonic/gin"

	api "github.com/StampWallet/backend/internal/api/models"
	"github.com/StampWallet/backend/internal/database"
	"github.com/StampWallet/backend/internal/managers"
	"github.com/StampWallet/backend/internal/middleware"
)

type AuthHandlers struct {
	authManager managers.AuthManager
	logger      *log.Logger
}

func CreateAuthHandlers(authManager managers.AuthManager, logger *log.Logger) *AuthHandlers {
	return &AuthHandlers{
		authManager: authManager,
		logger:      logger,
	}
}

// TODO share this with middleware
// Splits the token by :, expects exactly one :
func splitToken(token string) (string, string, error) {
	s := strings.Split(token, ":")
	if len(s) != 2 {
		return "", "", errors.New("invalid token format")
	}
	return s[0], s[1], nil
}

// Parses token from Authorization header value
func parseTokenFromHeader(header string) (string, string, error) {
	header_value_split := strings.Split(header, " ")
	if len(header_value_split) != 2 || header_value_split[0] != "Bearer" {
		return "", "", errors.New("invalid header formate")
	}

	return splitToken(header_value_split[1])
}

// Handles registration request
func (handler *AuthHandlers) postAccount(c *gin.Context) {
	// Parse request body
	req := api.PostAccountRequest{}
	if err := c.BindJSON(&req); err != nil {
		handler.logger.Printf("failed to parse in postAccount %+v", err)
		c.JSON(400, api.DefaultResponse{Status: api.INVALID_REQUEST})
		return
	}

	// Pass data to authManager, handle errors
	_, token, secret, err := handler.authManager.Create(managers.UserDetails{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		handler.logger.Printf("failed to authManager.Create in postAccount %+v", err)
		if err == managers.ErrEmailExists {
			c.JSON(409, api.DefaultResponse{Status: api.CONFLICT})
		} else if err == managers.ErrInvalidEmail {
			c.JSON(400, api.DefaultResponse{Status: api.INVALID_REQUEST})
		} else {
			c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
		}
		return
	}

	c.JSON(201, api.PostAccountResponse{Token: token.TokenId + ":" + secret})
}

// Handles login request
func (handler *AuthHandlers) postSession(c *gin.Context) {
	// Parse request body
	req := api.PostAccountSessionRequest{}
	if err := c.BindJSON(&req); err != nil {
		handler.logger.Printf("failed to parse in postSession %+v", err)
		c.JSON(400, api.DefaultResponse{Status: api.INVALID_REQUEST})
		return
	}

	// Pass login data to authManager, handle errors
	_, token, tokenSecret, err := handler.authManager.Login(req.Email, req.Password)
	if err != nil {
		handler.logger.Printf("failed to authManager.Login in postSession %+v", err)
		if err == managers.ErrInvalidLogin {
			c.JSON(401, api.DefaultResponse{Status: api.UNAUTHORIZED})
		} else {
			c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
		}
		return
	}

	c.JSON(200, api.PostAccountSessionResponse{Token: token.TokenId + ":" + tokenSecret})
}

// Handles logout request
func (handler *AuthHandlers) deleteSession(c *gin.Context) {
	// Parse request body - get the session token
	id, secret, err := parseTokenFromHeader(c.GetHeader("Authorization"))
	if err != nil {
		handler.logger.Printf("failed to parseTokenFromHeader in deleteSession %+v", err)
		c.JSON(401, api.DefaultResponse{Status: api.UNAUTHORIZED})
		return
	}

	// Pass token to authManager, handle errors
	_, _, err = handler.authManager.Logout(id, secret)
	if err != nil {
		handler.logger.Printf("failed to authManager.Logout in deleteSession %+v", err)
		if err == managers.ErrInvalidToken || err == managers.ErrInvalidTokenPurpose {
			c.JSON(401, api.DefaultResponse{Status: api.UNAUTHORIZED})
		} else {
			c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
		}
		return
	}

	c.JSON(200, api.DefaultResponse{Status: api.OK})
}

// Handles email confirmation request
func (handler *AuthHandlers) postAccountEmailConfirmation(c *gin.Context) {
	// Parse request body
	req := api.PostAccountEmailConfirmationRequest{}
	if err := c.BindJSON(&req); err != nil {
		handler.logger.Printf("failed to parse in postAccountEmailConfirmation %+v", err)
		c.JSON(400, api.DefaultResponse{Status: api.INVALID_REQUEST})
		return
	}

	// Parse token from request
	tokenId, tokenSecret, err := splitToken(req.Token)
	if err != nil {
		handler.logger.Printf("failed to splitToken in postAccountEmailConfirmation %+v", err)
		c.JSON(400, api.DefaultResponse{Status: api.INVALID_REQUEST})
		return
	}

	// Pass data to authManager
	_, err = handler.authManager.ConfirmEmail(tokenId, tokenSecret)
	if err != nil {
		handler.logger.Printf("failed to authManager.ChangePassword in postAccountEmailConfirmation %+v", err)
		if err == managers.ErrInvalidToken || err == managers.ErrInvalidTokenPurpose {
			c.JSON(401, api.DefaultResponse{Status: api.UNAUTHORIZED})
		} else {
			c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
		}
		return
	}

	c.JSON(200, api.DefaultResponse{Status: api.OK})
}

// Handles email change request
func (handler *AuthHandlers) postAccountPassword(c *gin.Context) {
	// Parse request body
	req := api.PostAccountPasswordRequest{}
	if err := c.BindJSON(&req); err != nil {
		handler.logger.Printf("failed to parse in postAccountPasswordRequest %+v", err)
		c.JSON(400, api.DefaultResponse{Status: api.INVALID_REQUEST})
		return
	}

	// Get user from context (should be inserted by authMiddleware)
	userAny, exists := c.Get("user")
	if !exists {
		handler.logger.Printf("user not available context")
		c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
		return
	}
	user := userAny.(*database.User)

	// Pass data to authManager, handle errors
	_, err := handler.authManager.ChangePassword(user, req.OldPassword, req.Password)
	if err != nil {
		handler.logger.Printf("failed to authManager.ChangePassword in postAccountPassword %+v", err)
		if err == managers.ErrInvalidOldPassword {
			c.JSON(400, api.DefaultResponse{Status: api.INVALID_REQUEST, Message: "INVALID_PASSWORD"})
		} else {
			c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
		}
		return
	}

	c.JSON(200, api.DefaultResponse{Status: api.OK})
}

// Handles email change request
func (handler *AuthHandlers) postAccountEmail(c *gin.Context) {
	// Parse request body
	req := api.PostAccountEmailRequest{}
	if err := c.BindJSON(&req); err != nil {
		handler.logger.Printf("failed to parse in postAccountEmail %+v", err)
		c.JSON(400, api.DefaultResponse{Status: api.INVALID_REQUEST})
		return
	}

	// Get user from context
	userAny, exists := c.Get("user")
	if !exists {
		handler.logger.Printf("user not available context")
		c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
		return
	}
	user := userAny.(*database.User)

	// Pass data to authManager, handle errors
	_, err := handler.authManager.ChangeEmail(user, req.Email)
	if err != nil {
		handler.logger.Printf("failed to authManager.ChangeEmail in postAccountEmail %+v", err)
		if err == managers.ErrEmailExists {
			c.JSON(409, api.DefaultResponse{Status: api.CONFLICT})
		} else if err == managers.ErrInvalidEmail {
			c.JSON(400, api.DefaultResponse{Status: api.INVALID_REQUEST})
		} else {
			c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
		}
		return
	}

	c.JSON(200, api.DefaultResponse{Status: api.OK})
}

func (handler *AuthHandlers) Connect(rg *gin.RouterGroup, authMiddleware middleware.AuthMiddleware) {
	account := rg.Group("/account")
	{
		account.POST("", handler.postAccount)
		account.POST("/emailConfirmation", handler.postAccountEmailConfirmation)
		account.POST("/email", authMiddleware.Handle, handler.postAccountEmail)
		account.POST("/password", authMiddleware.Handle, handler.postAccountPassword)
	}
	rg.POST("/sessions", handler.postSession)
	rg.DELETE("/sessions", authMiddleware.Handle, handler.deleteSession)
}
