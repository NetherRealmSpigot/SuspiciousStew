package types

import (
	"github.com/gin-gonic/gin"
)

type APIRoute struct {
	Path     string
	Method   string
	Handlers []gin.HandlerFunc
}
