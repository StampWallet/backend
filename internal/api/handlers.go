package api

import (
    "github.com/gin-gonic/gin"
    . "github.com/StampWallet/backend/internal/api/handlers"
    . "github.com/StampWallet/backend/internal/middleware"
)

type APIHandlers struct {
    authHandlers *AuthHandlers
    businessHandlers *BusinessHandlers
    userHandlers *UserHandlers
    fileHandlers *FileHandlers
}

func (handlers *APIHandlers) Connect(rg *gin.RouterGroup, authMiddleware *AuthMiddleware) {
}
