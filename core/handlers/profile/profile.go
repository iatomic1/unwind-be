package profile

import (
	"context"

	"github.com/adeyemialameen04/unwind-be/core/server"
	"github.com/adeyemialameen04/unwind-be/internal/db/repository"
	"github.com/adeyemialameen04/unwind-be/internal/domain"
	"github.com/adeyemialameen04/unwind-be/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/golodash/galidator"
)

type Handler struct {
	srv *server.Server
}

func NewProfileHandler(srv *server.Server) *Handler {
	if srv == nil {
		panic("server cannot be nil")
	}
	return &Handler{srv: srv}
}

// UpdateUserProfile godoc
//
//	@Summary		Update Profile
//	@Description	Updates a user profile including optional profile and cover pictures
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Security		AccessTokenBearer
//	@Param			book	body		repository.UpdateProfileParams	true	"Profile Data"
//	@Success		200		{object}	server.SuccessResponse{data=repository.Profile}
//	@Failure		401		{object}	server.UnauthorizedResponse			"Unauthorized"
//	@Failure		404		{object}	server.NotFoundResponse				"Profile not found"
//	@Failure		500		{object}	server.InternalServerErrorResponse	"Internal server error"
//	@Router			/user/profile [patch]
func (h *Handler) UpdateUserProfile(c *gin.Context) {
	var (
		ctx = context.Background()
		req repository.UpdateProfileParams
	)

	// Validate request body
	if err := h.validateUpdateProfileRequest(c, &req); err != nil {
		return // Error response already sent in validateRequest
	}

	// Get and validate profile ID from context
	profileID, err := domain.GetProfileIDFromContext(c)
	if err != nil {
		return // Error response already sent in getProfileIDFromContext
	}

	// Start transaction
	tx, err := h.srv.DB.Begin(ctx)
	if err != nil {
		server.SendInternalServerError(c, err, server.WithMessage("Failed to start transaction"))
		return
	}
	defer tx.Rollback(ctx)

	// Get existing profile
	repo := repository.New(tx)
	existingProfile, err := repo.GetProfileById(ctx, profileID)
	if err != nil {
		server.SendInternalServerError(c, err, server.WithMessage("Failed to fetch profile"))
		return
	}
	if existingProfile == nil {
		server.SendNotFound(c, nil, server.WithMessage("Profile not found"))
		return
	}

	// Handle image uploads and update request
	if err := h.handleImageUploads(c, &req); err != nil {
		return // Error response already sent in handleImageUploads
	}

	// Update profile
	req.ID = profileID
	profile, err := repo.UpdateProfile(ctx, req)
	if err != nil {
		server.SendInternalServerError(c, err, server.WithMessage("Failed to update profile"))
		return
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		server.SendInternalServerError(c, err, server.WithMessage("Failed to commit transaction"))
		return
	}

	server.SendSuccess(c, profile, server.WithMessage("Profile updated successfully"))
}

// @Summary		Get Profile
// @Description	Retrieves a user profile
// @Tags			User
// @Accept			json
// @Produce		json
// @Security		AccessTokenBearer
// @Param			book	body		repository.UpdateProfileParams	true	"Profile Data"
// @Success		200		{object}	server.SuccessResponse{data=repository.Profile}
// @Failure		401		{object}	server.UnauthorizedResponse			"Unauthorized"
// @Failure		404		{object}	server.NotFoundResponse				"Profile not found"
// @Failure		500		{object}	server.InternalServerErrorResponse	"Internal server error"
// @Router			/user/profile [get]
func (h *Handler) GetUserProfileByID(c *gin.Context) {
	ctx := context.Background()
	tx, err := h.srv.DB.Begin(ctx)
	if err != nil {
		server.SendInternalServerError(c, err)
		return
	}
	defer tx.Rollback(ctx)

	profileID, err := domain.GetProfileIDFromContext(c)
	if err != nil {
		return
	}

	repo := repository.New(tx)

	profile, err := repo.GetProfileById(ctx, profileID)
	if err := tx.Commit(ctx); err != nil {
		server.SendInternalServerError(c, err, server.WithMessage("Failed to commit transaction"))
		return
	}

	server.SendSuccess(c, profile, server.WithMessage("Profile retrieved successfully"))
}

func (h *Handler) validateUpdateProfileRequest(c *gin.Context, req *repository.UpdateProfileParams) error {
	g := galidator.New().CustomMessages(galidator.Messages{
		"required": "$field is required",
		"min":      `$field can't be less than {min}`,
	})
	customizer := g.Validator(repository.UpdateProfileParams{})

	if err := c.ShouldBindJSON(req); err != nil {
		server.SendValidationError(c, customizer.DecryptErrors(err))
		return err
	}
	return nil
}

func (h *Handler) handleImageUploads(c *gin.Context, req *repository.UpdateProfileParams) error {
	if req.CoverPic != nil {
		coverPicURL, err := utils.UploadImage(h.srv.Cld, *req.CoverPic, req.Username+"cover")
		if err != nil {
			server.SendInternalServerError(c, err, server.WithMessage("Failed to upload cover picture"))
			return err
		}
		req.CoverPic = &coverPicURL
	}

	if req.ProfilePic != nil {
		profilePicURL, err := utils.UploadImage(h.srv.Cld, *req.ProfilePic, req.Username+"profile")
		if err != nil {
			server.SendInternalServerError(c, err, server.WithMessage("Failed to upload profile picture"))
			return err
		}
		req.ProfilePic = &profilePicURL
	}

	return nil
}
