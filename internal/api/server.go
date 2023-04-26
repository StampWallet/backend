package api

import (
	. "github.com/StampWallet/backend/internal/config"
	"github.com/StampWallet/backend/internal/middleware"
	"github.com/gin-gonic/gin"
	"log"
)

type APIServer struct {
	router         *gin.Engine
	authMiddleware *middleware.AuthMiddleware
	apiHandlers    *APIHandlers
	logger         *log.Logger
	config         Config
}

func (server *APIServer) Start() {
}
