package api

import (
	"github.com/gin-gonic/gin"

	. "github.com/StampWallet/backend/internal/api/handlers"
	. "github.com/StampWallet/backend/internal/middleware"
)

type APIHandlers struct {
	AuthHandlers     *AuthHandlers
	BusinessHandlers *BusinessHandlers
	UserHandlers     *UserHandlers
	//FileHandlers     *FileHandlers
}

func (handlers *APIHandlers) Connect(rg *gin.RouterGroup, authMiddleware *AuthMiddleware,
	requireValidEmailMiddleware *RequireValidEmailMiddleware) {

	auth := rg.Group("/auth")
	handlers.AuthHandlers.Connect(auth, authMiddleware)

	business := rg.Group("/business", authMiddleware.Handle, requireValidEmailMiddleware.Handle)
	handlers.BusinessHandlers.Connect(business)

	user := rg.Group("/user", authMiddleware.Handle, requireValidEmailMiddleware.Handle)
	handlers.UserHandlers.Connect(user)

	//file := rg.Group("/file", authMiddleware.Handle, requireValidEmailMiddleware.Handle)
	//handlers.FileHandlers.Connect(file)
}
