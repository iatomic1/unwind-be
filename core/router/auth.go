package router

import (
	"github.com/adeyemialameen04/unwind-be/core/handlers/auth"
	"github.com/adeyemialameen04/unwind-be/core/server"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(srv *server.Server, router *gin.RouterGroup) {
	authHandler := auth.NewAuthHandler(srv)
	bookGroup := router
	bookGroup.POST("", authHandler.LoginUser)
	// router.POST("", handlers.CreateBook(context.Background(), srv.DB))
}
