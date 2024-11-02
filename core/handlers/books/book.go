package books

import (
	"net/http"

	"github.com/adeyemialameen04/unwind-be/mock"
	"github.com/adeyemialameen04/unwind-be/types"
	"github.com/gin-gonic/gin"
)

// GetBooks godoc
// @Summary      Get all books
// @Description  Get a list of all books
// @Tags         books
// @Accept       json
// @Produce      json
// @Success      200  {array}   types.Book
// @Router       /books [get]
func GetBooks(c *gin.Context) {
	c.JSON(http.StatusOK, mock.Books)
}

// CreateBook godoc
// @Summary      Create a new book
// @Description  Create a new book with the provided data
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        book  body      types.Book  true  "Book data"
// @Success      201   {object}  types.Book
// @Router       /books [post]
func CreateBook(c *gin.Context) {
	var newBook types.Book
	if err := c.BindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	mock.Books = append(mock.Books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
}
