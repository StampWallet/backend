package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
)

type RequireValidEmailMiddleware struct {
	logger *log.Logger
}

func (middleware *RequireValidEmailMiddleware) Handle(c *gin.Context) {

}
