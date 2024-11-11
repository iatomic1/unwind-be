package router

import (
	"net/http"
	"path/filepath"

	"github.com/adeyemialameen04/unwind-be/internal/projectpath"
	scalargo "github.com/bdpiprava/scalar-go"
	"github.com/gin-gonic/gin"
)

func RegisterDocsRoutes(router *gin.RouterGroup) {
	specUrl := filepath.Join(projectpath.Root, "/internal/docs/openapi.json")
	// spec2Url := filepath.Join(projectpath.Root, "/oas.json")

	router.GET("/reference", func(c *gin.Context) {
		content, err := scalargo.NewV2(
			scalargo.WithSpecURL("/api/v1/docs/swagger.json"),
			scalargo.WithMetaDataOpts(
				scalargo.WithTitle("Unwind"),
			),
			scalargo.WithTheme(scalargo.ThemeDeepSpace),
			scalargo.WithLayout(scalargo.LayoutClassic),
			scalargo.WithBaseServerURL("http://localhost:2020/"),
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, content)
	})

	router.GET("/swagger.json", func(c *gin.Context) {
		// c.File(spec2Url)
		c.File(specUrl)
	})
}
