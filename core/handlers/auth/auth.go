package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/adeyemialameen04/unwind-be/core/server"
	"github.com/adeyemialameen04/unwind-be/internal/db/repository"
	"github.com/adeyemialameen04/unwind-be/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/golodash/galidator"
	"github.com/jackc/pgx/v5/pgconn"
)

type RegisterResponse struct {
	AccessToken  string      `json:"accessToken"`
	RefreshToken string      `json:"refreshToken"`
	User         UserDetails `json:"user"`
}

const (
	UniqueViolation = "23505"
)

type UserDetails struct {
	Email string `json:"email"`
	ID    string `json:"id"`
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
// @Param        book  body      repository.RegisterUserParams  true  "Login data"
// @Failure      400   {object}  map[string]string            "Invalid request data"
// @Failure      500   {object}  map[string]string            "Failed to start transaction or insert book"
// @Router       /auth/login [post]
func (h *Handler) LoginUser(c *gin.Context) {
	g := galidator.New().CustomMessages(galidator.Messages{
		"required": "$field is required",
	})
	customizer := g.Validator(repository.RegisterUserParams{})

	var req repository.RegisterUserParams
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println(err)
		server.SendValidationError(c, customizer.DecryptErrors(err))
		return
	}

	server.SendSuccess(c, "got it", "lol")
}

// Signup godoc
// @Summary      Create an account
// @Description  Create an account on unwind
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        book  body      repository.RegisterUserParams  true  "Signup data"
// @Success      201   {object}  server.Response{data=RegisterResponse} "User created successfully"
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

	tx, err := h.srv.DB.Begin(ctx)
	if err != nil {
		server.SendInternalServerError(c, err)
		return
	}
	defer tx.Rollback(ctx)

	hashed_password, err := utils.HashPassword(req.Password)
	if err != nil {
		fmt.Println(err)
		return
	}
	repo := repository.New(tx)
	user, err := repo.RegisterUser(ctx, repository.RegisterUserParams{
		Email:    req.Email,
		Password: hashed_password,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == UniqueViolation {
			server.SendConflict(c, err)
			return
		}

		server.SendInternalServerError(c, err)
		return
	}

	// Commit the transaction
	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	tokens, err := h.generateTokens(user.ID.String())
	response := RegisterResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		User: UserDetails{
			Email: user.Email,
			ID:    user.ID.String(),
		},
	}

	server.SendSuccess(c, "User Created Successfully", response)
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

func (h *Handler) generateTokens(userID string) (Tokens, error) {
	accessToken, err := utils.CreateJWT(userID, false, h.srv.Config)
	if err != nil {
		return Tokens{}, err
	}

	refreshToken, err := utils.CreateJWT(userID, true, h.srv.Config)
	if err != nil {
		return Tokens{}, err
	}

	return Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
