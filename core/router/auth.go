package router

import (
	"github.com/adeyemialameen04/unwind-be/core/handlers/auth"
	"github.com/adeyemialameen04/unwind-be/core/server"
	"github.com/adeyemialameen04/unwind-be/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(srv *server.Server, router *gin.RouterGroup) {
	authHandler := auth.NewAuthHandler(srv)
	authGroup := router
	authGroup.POST("/login", authHandler.LoginUser)
	authGroup.POST("/signup", authHandler.RegisterUser)

	authGroup.Use(middleware.RefreshTokenMiddleware(srv.Config))
	{
		authGroup.GET("/refresh", authHandler.RefreshToken)
	}
}
