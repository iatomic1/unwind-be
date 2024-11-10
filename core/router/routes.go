package router

import (
	"github.com/adeyemialameen04/unwind-be/core/server"
	"github.com/gin-gonic/gin"
)

func SetupRouter(srv *server.Server) {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	RegisterAuthRoutes(srv, router.Group("/auth"))
	RegisterDocsRoutes(router.Group("/docs"))
	RegisterProfileROutes(srv, router.Group("/user"))

	srv.Router = router
}
