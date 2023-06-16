package api

import (
	"log"

	. "github.com/StampWallet/backend/internal/config"
	"github.com/StampWallet/backend/internal/database"
	"github.com/StampWallet/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

type APIServer struct {
	router                      *gin.Engine
	authMiddleware              *middleware.AuthMiddleware
	requireValidEmailMiddleware *middleware.RequireValidEmailMiddleware
	apiHandlers                 *APIHandlers
	logger                      *log.Logger
	config                      Config
}

func CreateAPIServer(
	authMiddleware *middleware.AuthMiddleware,
	requireValidEmailMiddleware *middleware.RequireValidEmailMiddleware,
	apiHandlers *APIHandlers,
	logger *log.Logger,
	config Config) *APIServer {

	server := &APIServer{
		router:                      gin.New(),
		authMiddleware:              authMiddleware,
		requireValidEmailMiddleware: requireValidEmailMiddleware,
		apiHandlers:                 apiHandlers,
		logger:                      logger,
		config:                      config,
	}

	server.router.Use(func(c *gin.Context) {
		c.Next()
		var uid uint
		user, exists := c.Get("user")
		if exists {
			uid = user.(*database.User).ID
		}
		logger.Printf("%s %s %d %d %d", c.Request.Method, c.Request.RequestURI,
			c.Writer.Status(), c.Writer.Size(), uid)
	})

	server.router.Use(gin.Recovery())

	server.apiHandlers.Connect(&server.router.RouterGroup,
		server.authMiddleware,
		server.requireValidEmailMiddleware)

	server.router.Static("/static", config.StaticPath)

	return server
}

func (server *APIServer) Start() error {
	return server.router.Run(server.config.ListenIP)
}
