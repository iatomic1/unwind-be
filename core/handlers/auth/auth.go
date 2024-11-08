package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/adeyemialameen04/unwind-be/core/server"
	"github.com/adeyemialameen04/unwind-be/internal/db/repository"
	"github.com/gin-gonic/gin"
	"github.com/golodash/galidator"
)

type AuthRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password"  binding:"required"`
}

type Handler struct {
	srv *server.Server
}

func NewAuthHandler(srv *server.Server) *Handler {
	return &Handler{srv: srv}
}

// Login godoc
// @Summary      Login to your account
// @Description  Logs a user into his/her account
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        book  body      AuthRequest  true  "Login data"
// @Failure      400   {object}  map[string]string            "Invalid request data"
// @Failure      500   {object}  map[string]string            "Failed to start transaction or insert book"
// @Router       /auth/login [post]
func (h *Handler) LoginUser(c *gin.Context) {
	g := galidator.New().CustomMessages(galidator.Messages{
		"required": "$field is required",
	})
	customizer := g.Validator(AuthRequest{})

	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println(err)
		server.SendValidationError(c, customizer.DecryptErrors(err))
		return
	}

	server.SendSuccess(c, "got it", "lol")
}

// Signup godoc
// @Summary      Login to your account
// @Description  Logs a user into his/her account
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        book  body      repository.RegisterUserParams  true  "Login data"
// @Failure      400   {object}  map[string]string            "Invalid request data"
// @Failure      500   {object}  map[string]string            "Failed to start transaction or insert book"
// @Router       /auth/signup [post]
func (h *Handler) RegisterUser(c *gin.Context) {
	g := galidator.New().CustomMessages(galidator.Messages{
		"required": "$field is required",
	})
	ctx := context.Background()
	customizer := g.Validator(repository.RegisterUserParams{})

	var req repository.RegisterUserParams
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println(err)
		server.SendValidationError(c, customizer.DecryptErrors(err))
		return
	}

	// Start a transaction
	tx, err := h.srv.DB.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback(ctx)

	// Insert book using the repository
	repo := repository.New(tx)
	user, err := repo.RegisterUser(ctx, req)
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

	server.SendSuccess(c, "got it", user)
}
