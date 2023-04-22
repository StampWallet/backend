package api

import (
    "log"
    "github.com/gin-gonic/gin"
    "github.com/StampWallet/backend/internal/middleware"
    . "github.com/StampWallet/backend/internal/config"
)

type APIServer struct {
    router *gin.Engine
    authMiddleware *middleware.AuthMiddleware
    apiHandlers *APIHandlers
    logger *log.Logger
    config Config
}

func (server *APIServer) Start() {
}
