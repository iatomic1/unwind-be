package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/adeyemialameen04/unwind-be/core/server"
	"github.com/adeyemialameen04/unwind-be/internal/db/repository"
	utils "github.com/adeyemialameen04/unwind-be/internal/utils/auth"
	"github.com/gin-gonic/gin"
	"github.com/golodash/galidator"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	UniqueViolation           = "23505"
	UserCreated               = "User Created Successfully"
	ErrEmailAlreadyExist      = "User with this email already exists"
	ErrInvalidEmailOrPassword = "Invalid email or password"
	LoginSuccessful           = "Login Successful"
	ErrGeneratingTokens       = "Error generating token pair"
	TokensRefreshed           = "Tokens Refreshed successfully"
)

type RegisterResponse struct {
	utils.TokenPair
	User UserDetails `json:"user"`
}

type RefreshTokenResponse struct {
	utils.TokenPair
}

type UserDetails struct {
	Email string `json:"email"`
	ID    string `json:"id"`
} // @name EmailID

type Handler struct {
	srv *server.Server
}

func NewAuthHandler(srv *server.Server) *Handler {
	return &Handler{srv: srv}
}

// Login godoc
//
//	@Summary		Login to your account
//	@Description	Logs a user into his/her account
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			book	body		repository.RegisterUserParams	true	"Login data"
//	@Success		201		{object}	server.Response{data=RegisterResponse}	"Login success"
//	@Failure		400		{object}	map[string]string				"Invalid request data"
//	@Failure		500		{object}	map[string]string				"Failed to start transaction or insert book"
//	@Router			/auth/login [post]
func (h *Handler) LoginUser(c *gin.Context) {
	g := galidator.New().CustomMessages(galidator.Messages{
		"required": "$field is required",
		"min":      `$field can't be less than {min}`,
	})
	customizer := g.Validator(repository.RegisterUserParams{})

	ctx := context.Background()
	var req repository.RegisterUserParams
	if err := c.ShouldBindJSON(&req); err != nil {
		server.SendValidationError(c, customizer.DecryptErrors(err))
		return
	}

	tx, err := h.srv.DB.Begin(ctx)
	if err != nil {
		server.SendInternalServerError(c, err)
		return
	}
	defer tx.Rollback(ctx)

	repo := repository.New(tx)
	user, err := repo.GetUserByEmail(ctx, req.Email)
	verify := utils.VerifyPassword(req.Password, user.Password)

	if err != nil || !verify {
		server.SendUnauthorized(c, nil, server.WithMessage(ErrInvalidEmailOrPassword))
		return
	}

	tokens, err := utils.GenerateTokenPair(user.ID.String(), h.srv.Config)
	if err != nil {
		server.SendInternalServerError(c, err)
	}

	response := RegisterResponse{
		TokenPair: utils.TokenPair{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		},
		User: UserDetails{
			Email: user.Email,
			ID:    user.ID.String(),
		},
	}

	server.SendSuccess(c, response, server.WithMessage(LoginSuccessful))
}

// Signup godoc
//
//	@Summary		Create an account
//	@Description	Create an account on unwind
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			EmailAndPassword	body		repository.RegisterUserParams			true	"Signup data"
//	@Success		201		{object}	server.Response{data=RegisterResponse}	"User created successfully"
//	@Failure		400		{object}	map[string]string						"Invalid request data"
//	@Failure		500		{object}	map[string]string						"Failed to start transaction or insert book"
//	@Router			/auth/signup [post]
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
			server.SendConflict(c, err, server.WithMessage(ErrEmailAlreadyExist))
			return
		}

		server.SendInternalServerError(c, err)
		return
	}

	// Commit the transaction
	if err := tx.Commit(ctx); err != nil {
		server.SendInternalServerError(c, err, server.WithMessage("Failed to commit transaction"))
		return
	}

	tokens, err := utils.GenerateTokenPair(user.ID.String(), h.srv.Config)
	response := RegisterResponse{
		TokenPair: utils.TokenPair{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		},
		User: UserDetails{
			Email: user.Email,
			ID:    user.ID.String(),
		},
	}

	server.SendSuccess(c, response, server.WithMessage(UserCreated))
}

// @Summary Refreh Token
// @Description Refreshes token to get new token pair
// @Tags Auth
// @Accept json
// @Produce json
// @Router /auth/refresh [get]
func (h *Handler) RefreshToken(c *gin.Context) {
	ctx := context.Background()

	tx, err := h.srv.DB.Begin(ctx)
	if err != nil {
		server.SendInternalServerError(c, err)
		return
	}
	defer tx.Rollback(ctx)

	userId, ok := c.Get("userId")
	if !ok {
		server.SendUnauthorized(c, nil, server.WithMessage("UserId not found for some reason"))
		return
	}

	parsedUUID, err := uuid.Parse(userId.(string))
	if err != nil {
		server.SendInternalServerError(c, err, server.WithMessage("Error parsing uuid"))
		return
	}

	repo := repository.New(tx)
	user, err := repo.GetUserById(ctx, parsedUUID)
	fmt.Println(user)
	if err != nil {
		server.SendInternalServerError(c, err)
		return
	}

	tokens, err := utils.GenerateTokenPair(userId.(string), h.srv.Config)
	if err != nil {
		server.SendInternalServerError(c, err, server.WithMessage(ErrGeneratingTokens))
		return
	}

	response := utils.TokenPair{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}
	server.SendSuccess(c, response, server.WithMessage(TokensRefreshed))
}
