package api

import (
    "log"
    "github.com/gin-gonic/gin"
    . "github.com/StampWallet/backend/internal/services"
)

type FileHandlers struct {
    fileStorageService *FileStorageService
    logger *log.Logger
    //this won't work
    //userAuthorizedAcessor *UserAuthorizedAccessor
}

func (handler *FileHandlers) getFile(c *gin.Context) {
               
}              
               
func (handler *FileHandlers) uploadFile(c *gin.Context) {
               
}              
               
func (handler *FileHandlers) Connect(rg *gin.RouterGroup) {

}
