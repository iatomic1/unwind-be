package profile

import (
	"context"
	"fmt"

	"github.com/adeyemialameen04/unwind-be/core/server"
	"github.com/adeyemialameen04/unwind-be/internal/db/repository"
	"github.com/gin-gonic/gin"
	"github.com/golodash/galidator"
	"github.com/google/uuid"
)

type Handler struct {
	srv *server.Server
}

func NewProfileHandler(srv *server.Server) *Handler {
	return &Handler{
		srv: srv,
	}
}

// @Summary		Update Profile
// @Description	Updates a user profile
// @Tags			User
// @Accept			json
// @Produce		json
// @Security		AccessTokenBearer
// @Param			book	body		repository.UpdateProfileParams	true	"Profile Data"
// @Success		200		{object}	server.Response{data=repository.Profile}
// @Router			/user/profile [patch]
func (h *Handler) UpdateUserProfile(c *gin.Context) {
	g := galidator.New().CustomMessages(galidator.Messages{
		"required": "$field is required",
		"min":      `$field can't be less than {min}`,
	})
	customizer := g.Validator(repository.UpdateProfileParams{})

	ctx := context.Background()
	var req repository.UpdateProfileParams
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

	profileId, ok := c.Get("profileId")
	if !ok {
		server.SendUnauthorized(c, nil, server.WithMessage("profileId not found for some reason"))
		return
	}
	parsedProfileId, err := uuid.Parse(profileId.(string))
	if err != nil {
	}
	fmt.Println(profileId)

	req.ID = parsedProfileId
	repo := repository.New(tx)
	profile, err := repo.UpdateProfile(ctx, req)
	if err := tx.Commit(ctx); err != nil {
		server.SendInternalServerError(c, err, server.WithMessage("Failed to commit transaction"))
		return
	}

	if err != nil {
		server.SendInternalServerError(c, err, server.WithMessage("Here"))
		return
	}

	server.SendSuccess(c, profile, server.WithMessage("Profile updated successfully"))
}
