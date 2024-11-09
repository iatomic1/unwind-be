package router

import (
	"github.com/adeyemialameen04/unwind-be/core/handlers/books"
	"github.com/adeyemialameen04/unwind-be/core/server"
	"github.com/adeyemialameen04/unwind-be/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterBookRoutes(srv *server.Server, router *gin.RouterGroup) {
	bookHandler := books.NewBookHandler(srv)
	bookGroup := router
	bookGroup.GET("", bookHandler.GetBooks)
	bookGroup.Use(middleware.AccessTokenMiddleware(srv.Config))
	{
		bookGroup.POST("", bookHandler.CreateBook)
	}
}
