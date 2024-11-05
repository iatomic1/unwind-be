// core/handlers/books/book.go
package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/adeyemialameen04/unwind-be/internal/db/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

// GetBooks godoc
// @Summary      Get all books
// @Description  Get a list of all books
// @Tags         books
// @Accept       json
// @Produce      json
// @Success      200  {array}   repository.Book
// @Router       /books [get]
func GetBooks(ctx context.Context, conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		tx, err := conn.Begin(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
			return
		}
		defer tx.Rollback(ctx)

		repo := repository.New(tx)
		books, err := repo.FindAllBooks(ctx)
		if err != nil {
			log.Fatal(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch books"})
			return
		}

		if err := tx.Commit(ctx); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
			return
		}

		c.JSON(http.StatusOK, books)
	}
}

// CreateBook godoc
// @Summary      Create a new book
// @Description  Insert a new book into the database
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        book  body      repository.InsertBookParams  true  "Book data"
// @Success      201   {object}  repository.Book             "Book created successfully"
// @Failure      400   {object}  map[string]string            "Invalid request data"
// @Failure      500   {object}  map[string]string            "Failed to start transaction or insert book"
// @Router       /books [post]
func CreateBook(ctx context.Context, conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req repository.InsertBookParams

		// Bind the incoming JSON to the InsertBookParams struct
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
			return
		}

		// Start a transaction
		tx, err := conn.Begin(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
			return
		}
		defer tx.Rollback(ctx)

		// Insert book using the repository
		repo := repository.New(tx)
		book, err := repo.InsertBook(ctx, req)
		if err != nil {
			log.Fatal(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert book"})
			return
		}

		// Commit the transaction
		if err := tx.Commit(ctx); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Book created successfully", "data": book})
	}
}
