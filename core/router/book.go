package router

import (
	"context"

	"github.com/adeyemialameen04/unwind-be/core/handlers"
	"github.com/adeyemialameen04/unwind-be/core/server"
	"github.com/gin-gonic/gin"
)

func RegisterBookRoutes(srv *server.Server, router *gin.RouterGroup) {
	router.GET("", handlers.GetBooks(context.Background(), srv.DB))
	router.POST("", handlers.CreateBook(context.Background(), srv.DB))
}
