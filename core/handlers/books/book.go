package books

import (
	"context"
	"net/http"

	"github.com/adeyemialameen04/unwind-be/internal/db/repository"
	"github.com/adeyemialameen04/unwind-be/mock"
	"github.com/adeyemialameen04/unwind-be/types"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
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

func GetB(ctx context.Context, conn *pgx.Conn) (error, []repository.Book) {
	tx, _ := conn.Begin(ctx)
	defer tx.Rollback(ctx)

	repo := repository.New(tx)

	books, err := repo.FindAllBooks(ctx)
	if err != nil {
		return err, nil
	}

	return nil, books
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
