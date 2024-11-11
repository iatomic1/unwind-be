package router

import (
	"github.com/adeyemialameen04/unwind-be/core/server"
	"github.com/gin-gonic/gin"
)

func SetupRouter(srv *server.Server) {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// Create the base api group
	api := router.Group(srv.Config.ApiPrefixStr)
	{
		RegisterAuthRoutes(srv, api.Group("/auth"))
		RegisterDocsRoutes(api.Group("/docs"))
		RegisterProfileROutes(srv, api.Group("/user"))
	}

	srv.Router = router
}
