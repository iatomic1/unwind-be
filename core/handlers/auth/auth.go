package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/adeyemialameen04/unwind-be/core/server"
	"github.com/adeyemialameen04/unwind-be/internal/db/repository"
	"github.com/adeyemialameen04/unwind-be/internal/domain"
	"github.com/adeyemialameen04/unwind-be/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/golodash/galidator"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
)

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
//	@Param			EmailAndPassword	body		repository.RegisterUserParams				true	"Login data"
//	@Success		201					{object}	server.Response{data=domain.AuthResponse}	"Login success"
//	@Failure		400					{object}	map[string]string							"Invalid request data"
//	@Failure		500					{object}	server.InternalServerErrorResponse			"Internal server error"
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
		server.SendUnauthorized(c, nil, server.WithMessage(domain.ErrInvalidEmailOrPassword))
		return
	}

	profile, err := repo.GetProfileByUserId(ctx, user.ID)
	if err != nil {
		server.SendBadRequest(c, err, server.WithMessage("Error retrieveing user profile"))
		return
	}

	tokens, err := utils.GenerateTokenPair(utils.EmailID{
		Email:     user.Email,
		ID:        user.ID.String(),
		ProfileID: profile.ID.String(),
	}, h.srv.Config)
	if err != nil {
		server.SendInternalServerError(c, err)
	}

	response := domain.AuthResponse{
		TokenPair: utils.TokenPair{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		},
		User: domain.EmailID{
			Email:     user.Email,
			ID:        user.ID.String(),
			ProfileID: profile.ID.String(),
		},
	}

	server.SendSuccess(c, response, server.WithMessage(domain.LoginSuccessful))
}

// Signup godoc
//
//	@Summary		Create an account
//	@Description	Create an account on unwind
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			EmailAndPassword	body		domain.RegisterRequest						true	"Signup data"
//	@Success		201					{object}	server.Response{data=domain.AuthResponse}	"User created successfully"
//	@Failure		500					{object}	server.InternalServerErrorResponse			"Internal server error"
//	@Router			/auth/signup [post]
func (h *Handler) RegisterUser(c *gin.Context) {
	g := galidator.New().CustomMessages(galidator.Messages{
		"required": "$field is required",
	})
	ctx := context.Background()
	customizer := g.Validator(domain.RegisterRequest{})

	var req domain.RegisterRequest
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
		if errors.As(err, &pgErr) && pgErr.Code == domain.UniqueViolation {
			server.SendConflict(c, err, server.WithMessage(domain.ErrEmailAlreadyExist))
			return
		}

		server.SendInternalServerError(c, err)
		return
	}

	parsedUUID, err := uuid.Parse(user.ID.String())
	if err != nil {
		server.SendInternalServerError(c, err, server.WithMessage("Error parsing uuid"))
		return
	}

	profile, err := repo.InsertProfile(ctx, repository.InsertProfileParams{
		Username: req.Username,
		UserID:   parsedUUID,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == domain.UniqueViolation {
			server.SendConflict(c, err, server.WithMessage(domain.ErrUsernameAlreadyExist))
			return
		}

		server.SendInternalServerError(c, err)
		return
	}

	if err := tx.Commit(ctx); err != nil {
		server.SendInternalServerError(c, err, server.WithMessage("Failed to commit transaction"))
		return
	}

	tokens, err := utils.GenerateTokenPair(utils.EmailID{
		Email:     user.Email,
		ID:        user.ID.String(),
		ProfileID: profile.ID.String(),
	}, h.srv.Config)
	response := domain.AuthResponse{
		TokenPair: utils.TokenPair{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		},
		User: domain.EmailID{
			Email:     user.Email,
			ID:        user.ID.String(),
			ProfileID: profile.ID.String(),
		},
	}

	server.SendCreated(c, response, server.WithMessage(domain.UserCreated))
}

// handleUserCreation manages user creation or retrieval logic
func (h *Handler) handleUserCreation(ctx context.Context, repo *repository.Queries, gothUser goth.User) (*utils.EmailID, error) {
	existingUser, err := repo.GetUserByEmail(ctx, gothUser.Email)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	if existingUser != nil {

		parsedUUID, err := uuid.Parse(existingUser.ID.String())
		if err != nil {
			log.Println("Error parsing uuid")
			return nil, err
		}
		existingProfile, err := repo.GetProfileByUserId(ctx, parsedUUID)
		if err != nil {
			if err == pgx.ErrNoRows {
				log.Println("Error retrieving usern profile")
				return nil, err
			}
		}

		return &utils.EmailID{
			Email:     existingUser.Email,
			ID:        existingUser.ID.String(),
			ProfileID: existingProfile.ID.String(),
		}, nil
	}

	// Create new user
	userParams := repository.RegisterUserParams{
		Email:    gothUser.Email,
		Password: uuid.New().String(),
	}

	user, err := repo.RegisterUser(ctx, userParams)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == domain.UniqueViolation {
			return nil, fmt.Errorf("email already exists: %w", err)
		}
		return nil, err
	}

	// Create user profile
	username := generateUsername(gothUser.Name, gothUser.Email)
	profile, err := repo.InsertProfile(ctx, repository.InsertProfileParams{
		Username:   username,
		UserID:     user.ID,
		Name:       &gothUser.Name,
		ProfilePic: &gothUser.AvatarURL,
	})
	if err != nil {
		return nil, err
	}

	return &utils.EmailID{
		Email:     user.Email,
		ID:        user.ID.String(),
		ProfileID: profile.ID.String(),
	}, nil
}

func (h *Handler) OAuthCallback(c *gin.Context) {
	provider := c.Param("provider")
	if !isValidProvider(provider) {
		server.SendBadRequest(c, fmt.Errorf("invalid provider: %s", provider))
		return
	}

	req := c.Request.WithContext(context.WithValue(c.Request.Context(), "provider", provider))
	c.Request = req

	gothUser, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		server.SendInternalServerError(c, err)
		return
	}
	fmt.Println(gothUser)
	c.SetCookie("test", "LFG", 3600, "/", "", false, true)
}

func (h *Handler) OauthBegin(c *gin.Context) {
	provider := c.Param("provider")
	req := c.Request.WithContext(context.WithValue(c.Request.Context(), "provider", provider))
	c.Request = req

	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func isValidProvider(provider string) bool {
	validProviders := map[string]bool{
		"google":   true,
		"github":   true,
		"facebook": true,
		// Add more providers as needed
	}
	return validProviders[provider]
}

func generateUsername(name, email string) string {
	if name == "" {
		return strings.Split(email, "@")[0]
	}
	// Remove spaces and special characters, convert to lowercase
	username := strings.ToLower(strings.Join(strings.Fields(name), ""))
	username = regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(username, "")

	// Append random string to ensure uniqueness
	randomStr := uuid.New().String()[:6]
	return fmt.Sprintf("%s_%s", username, randomStr)
}

// @Summary		Refresh Token
// @Description	Refreshes token to get new token pair
// @Security		RefreshTokenBearer
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Success		200	{object}	server.SuccessResponse{data=utils.TokenPair}	"TokenPair"
// @Failure		401	{object}	server.UnauthorizedResponse						"Unauthorized"
// @Failure		404	{object}	server.NotFoundResponse							"Profile not found"
// @Failure		500	{object}	server.InternalServerErrorResponse				"Internal server error"
// @Router			/auth/refresh [get]
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

	profileId, ok := c.Get("profileId")
	if !ok {
		server.SendUnauthorized(c, nil, server.WithMessage("profileId not found for some reason"))
		return
	}

	email, ok := c.Get("email")
	if !ok {
		server.SendUnauthorized(c, nil, server.WithMessage("profileId not found for some reason"))
		return
	}

	parsedUserId, parsedProfileId, err := domain.ParseIDs(userId.(string), profileId.(string))
	if err != nil {
		server.SendInternalServerError(c, err, server.WithMessage("Error parsing uuid"))
		return
	}

	repo := repository.New(tx)
	user, err := repo.GetUserById(ctx, parsedUserId)
	fmt.Println(user)
	if err != nil {
		server.SendInternalServerError(c, err)
		return
	}

	tokens, err := utils.GenerateTokenPair(utils.EmailID{
		ID:        parsedUserId.String(),
		Email:     email.(string),
		ProfileID: parsedProfileId.String(),
	}, h.srv.Config)
	if err != nil {
		server.SendInternalServerError(c, err, server.WithMessage(domain.ErrGeneratingTokens))
		return
	}

	response := utils.TokenPair{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}
	server.SendSuccess(c, response, server.WithMessage(domain.TokensRefreshed))
}
