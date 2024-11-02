package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"runtime"

	"github.com/adeyemialameen04/unwind-be/internal/config"
	"github.com/adeyemialameen04/unwind-be/mock"
	"github.com/adeyemialameen04/unwind-be/types"
	"github.com/gin-gonic/gin"
)

// getRootDir returns the absolute path to the project root directory
func getRootDir() string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(b), "../..")
}

func getBooks(c *gin.Context) {
	c.JSON(http.StatusOK, mock.Books)
}

func createBook(c *gin.Context) {
	var newBook types.Book

	if err := c.BindJSON(&newBook); err != nil {
		return
	}

	mock.Books = append(mock.Books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
}

func main() {
	rootDir := getRootDir()
	cfg, err := config.Load(rootDir)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	router := gin.Default()
	router.GET("/books", getBooks)
	router.POST("/books", createBook)

	serverAddr := fmt.Sprintf("localhost%s", cfg.HttpAddress)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
