package main

import (
	"net/http"

	"github.com/adeyemialameen04/unwind-be/mock"
	"github.com/adeyemialameen04/unwind-be/types"
	"github.com/gin-gonic/gin"
)

func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, mock.Books)
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
	router := gin.Default()
	router.GET("/books", getBooks)

	router.Run("localhost:8080")
}
