package router

import (
	"github.com/adeyemialameen04/unwind-be/core/handlers/profile"
	watchlist "github.com/adeyemialameen04/unwind-be/core/handlers/watch-list"
	"github.com/adeyemialameen04/unwind-be/core/server"
	"github.com/adeyemialameen04/unwind-be/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(srv *server.Server, router *gin.RouterGroup) {
	profileHandler := profile.NewProfileHandler(srv)
	watchListHandler := watchlist.NewWatchListHandler(srv)

	watchListGroup := router.Group("/watch-list")
	profileGroup := router.Group("/profile")

	profileGroup.Use(middleware.AccessTokenMiddleware(srv.Config))
	{
		profileGroup.PATCH("", profileHandler.UpdateUserProfile)
		profileGroup.GET("", profileHandler.GetUserProfileByID)
	}

	watchListGroup.Use(middleware.AccessTokenMiddleware(srv.Config))
	{
		watchListGroup.GET("", watchListHandler.GetWatchList)
		watchListGroup.POST("", watchListHandler.AddToList)
		watchListGroup.GET("/:id", watchListHandler.GetWatchListByMediaID)
		watchListGroup.PATCH("/:id", watchListHandler.UpdateWatchListStatus)
		watchListGroup.DELETE("/:id", watchListHandler.DeleteFromWatchList)
	}
}
