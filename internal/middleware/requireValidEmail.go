package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
)

type RequireValidEmailMiddleware struct {
	logger *log.Logger
}

func CreateRequireValidEmailMiddleware(logger *log.Logger) *RequireValidEmailMiddleware {
	return &RequireValidEmailMiddleware{
		logger: logger,
	}
}

func (middleware *RequireValidEmailMiddleware) Handle(c *gin.Context) {

}
