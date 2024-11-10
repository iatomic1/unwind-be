package router

import (
	"github.com/adeyemialameen04/unwind-be/core/handlers/profile"
	"github.com/adeyemialameen04/unwind-be/core/server"
	"github.com/adeyemialameen04/unwind-be/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterProfileROutes(srv *server.Server, router *gin.RouterGroup) {
	profileHandler := profile.NewProfileHandler(srv)
	profileGroup := router
	profileGroup.Use(middleware.AccessTokenMiddleware(srv.Config))
	{
		profileGroup.PATCH("/profile", profileHandler.UpdateUserProfile)
	}
}
