package router

import (
	"github.com/adeyemialameen04/unwind-be/core/handlers/books"
	"github.com/adeyemialameen04/unwind-be/core/server"
	"github.com/gin-gonic/gin"
)

func RegisterBookRoutes(srv *server.Server, router *gin.RouterGroup) {
	router.GET("", books.GetBooks)
}
