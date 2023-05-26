package api

import (
	"log"

	. "github.com/StampWallet/backend/internal/database/accessors"
	. "github.com/StampWallet/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type FileHandlers struct {
	fileStorageService    FileStorageService
	logger                *log.Logger
	userAuthorizedAcessor UserAuthorizedAccessor
}

func (handler *FileHandlers) getFile(c *gin.Context) {

}

func (handler *FileHandlers) postFile(c *gin.Context) {

}

func (handler *FileHandlers) deleteFile(c *gin.Context) {

}

func (handler *FileHandlers) Connect(rg *gin.RouterGroup) {

}
