package api

import (
	. "github.com/StampWallet/backend/internal/services"
	"github.com/gin-gonic/gin"
	"log"
)

type FileHandlers struct {
	fileStorageService *FileStorageService
	logger             *log.Logger
	//this won't work
	//userAuthorizedAcessor *UserAuthorizedAccessor
}

func (handler *FileHandlers) getFile(c *gin.Context) {

}

func (handler *FileHandlers) uploadFile(c *gin.Context) {

}

func (handler *FileHandlers) Connect(rg *gin.RouterGroup) {

}
