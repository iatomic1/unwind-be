package router

import (
	"github.com/adeyemialameen04/unwind-be/core/server"
	"github.com/gin-gonic/gin"
)

func SetupRouter(srv *server.Server) {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	RegisterBookRoutes(srv, router.Group("/books"))
	RegisterDocsRoutes(router.Group("/docs"))

	srv.Router = router
}
