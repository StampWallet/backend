package api

import (
	"io"
	"log"

	"github.com/gin-gonic/gin"

	api "github.com/StampWallet/backend/internal/api/models"
	"github.com/StampWallet/backend/internal/database"
	. "github.com/StampWallet/backend/internal/database/accessors"
	. "github.com/StampWallet/backend/internal/services"
	. "github.com/StampWallet/backend/internal/utils"
)

// FileHandlers is a struct that implements handlers for all operations under "/file/" URL path.
// Operations include: uploading files to existing FileMetadata, deleting files from FileMetadata and
// downloading files (does not require the requester to be the owner of the file).
type FileHandlers struct {
	fileStorageService    FileStorageService
	logger                *log.Logger
	userAuthorizedAcessor UserAuthorizedAccessor
}

// Retuns FileMetadata with fileId PublicId if owned by user
// On error, returns nil and responds with appropriate HTTP data
func (handler *FileHandlers) getFileById(c *gin.Context, user *database.User, fileId string) *database.FileMetadata {
	fileMetadataTmp, err := handler.userAuthorizedAcessor.Get(user, &database.FileMetadata{PublicId: fileId})
	if err == ErrNoAccess {
		c.JSON(403, api.DefaultResponse{Status: api.FORBIDDEN})
		return nil
	} else if err == ErrNotFound {
		c.JSON(404, api.DefaultResponse{Status: api.NOT_FOUND})
		return nil
	} else if err != nil {
		handler.logger.Printf("%s unknown error after userAuthorizedAcessor.Get: %+v", CallerFilename(), err)
		c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
		return nil
	}

	return fileMetadataTmp.(*database.FileMetadata)
}

// Handles get file request
// Requires fileId path parameter which will contain FileMetadata.PublicId
func (handler *FileHandlers) getFile(c *gin.Context) {
	fileId := c.Param("fileId")

	// Get file by id, handle errors
	file, mimetype, err := handler.fileStorageService.GetData(fileId)
	if err == ErrNoSuchFile || err == ErrFileNotUploaded {
		c.JSON(404, api.DefaultResponse{Status: api.NOT_FOUND})
		return
	} else if err != nil {
		handler.logger.Printf("%s unknown error after fileStorageService.GetData: %+v", CallerFilename(), err)
		c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
		return
	}

	// Write Content-Type and file
	c.Header("Content-Type", mimetype)
	_, err = io.Copy(c.Writer, file)
	//content-type already set, cant just return a json now
	if err != nil {
		handler.logger.Printf("%s unknown error after io.Copy: %+v", CallerFilename(), err)
		return
	}
}

// Handles upload file request
// Requires fileId path parameter which will contain FileMetadata.PublicId
// Requires logged in user to be inserted into context under "user"
func (handler *FileHandlers) postFile(c *gin.Context) {
	fileId := c.Param("fileId")

	// Get user from context (should be inserted by authMiddleware)
	user := getUserFromContext(handler.logger, c)
	if user == nil {
		return
	}

	// Get FileMetadata
	fileMetadata := handler.getFileById(c, user, fileId)
	if fileMetadata == nil {
		return
	}

	// Get Content-Type from headers, return error if missing
	mimetype := c.GetHeader("Content-Type")
	if mimetype == "" {
		c.JSON(400, api.DefaultResponse{Status: api.INVALID_REQUEST, Message: "INVALID_CONTENT_TYPE"})
		return
	}

	// Pass request bdoy to fileStorageService
	_, err := handler.fileStorageService.Upload(*fileMetadata, c.Request.Body, mimetype)
	if err == ErrInvalidMimeType {
		c.JSON(400, api.DefaultResponse{Status: api.INVALID_REQUEST, Message: "INVALID_CONTENT_TYPE"})
		return
	} else if err != nil {
		handler.logger.Printf("%s unknown error after fileStorageService.Upload: %+v", CallerFilename(), err)
		c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
		return
	}

	c.JSON(200, api.DefaultResponse{Status: api.OK})
	return
}

// Handles delete file request
// Requires fileId path parameter which will contain FileMetadata.PublicId
// Requires logged in user to be inserted into context under "user"
func (handler *FileHandlers) deleteFile(c *gin.Context) {
	fileId := c.Param("fileId")

	// Get user from context (should be inserted by authMiddleware)
	user := getUserFromContext(handler.logger, c)
	if user == nil {
		return
	}

	// Get FileMetadata
	fileMetadata := handler.getFileById(c, user, fileId)
	if fileMetadata == nil {
		return
	}

	// Pass removal request to FileStorageService, handle errors
	err := handler.fileStorageService.RemoveFile(*fileMetadata)
	if err == ErrNoSuchFile || err == ErrFileNotUploaded {
		c.JSON(404, api.DefaultResponse{Status: api.NOT_FOUND})
		return
	} else if err != nil {
		handler.logger.Printf("%s unknown error after fileStorageService.Remove: %+v", CallerFilename(), err)
		c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
		return
	}

	c.JSON(200, api.DefaultResponse{Status: api.OK})
	return
}

func CreateFileHandlers(fileStorageService FileStorageService, logger *log.Logger,
	userAuthorizedAcessor UserAuthorizedAccessor) *FileHandlers {
	return &FileHandlers{
		fileStorageService:    fileStorageService,
		logger:                logger,
		userAuthorizedAcessor: userAuthorizedAcessor,
	}
}

func (handler *FileHandlers) Connect(rg *gin.RouterGroup) {
	rg.GET("/:fileId", handler.getFile)
	rg.POST("/:fileId", handler.postFile)
	rg.DELETE("/:fileId", handler.deleteFile)
}
