package watchlist

import (
	"context"
	"errors"

	"github.com/adeyemialameen04/unwind-be/core/server"
	"github.com/adeyemialameen04/unwind-be/internal/db/repository"
	"github.com/adeyemialameen04/unwind-be/internal/domain"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	srv *server.Server
}

func NewWatchListHandler(srv *server.Server) *Handler {
	if srv == nil {
		panic("server cannot be nil")
	}
	return &Handler{srv: srv}
}

// GetWatchList godoc
//
//	@Summary		Get User WatchList
//	@Description	Retrieves a user watch list.
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Security		AccessTokenBearer
//	@Success		200		{object}	server.SuccessResponse{data=repository.WatchList}
//	@Failure		401		{object}	server.UnauthorizedResponse			"Unauthorized"
//	@Failure		404		{object}	server.NotFoundResponse				"Profile not found"
//	@Failure		500		{object}	server.InternalServerErrorResponse	"Internal server error"
//	@Router			/user/watch-list [get]
func (h *Handler) GetWatchList(c *gin.Context) {
	ctx := context.Background()
	tx, err := h.srv.DB.Begin(ctx)
	if err != nil {
		server.SendInternalServerError(c, err, server.WithMessage("Error beginning transaction"))
		return
	}
	defer tx.Rollback(ctx)

	userId, err := domain.GetUserIDFromContext(c)
	if err != nil {
		return
	}

	repo := repository.New(tx)
	watchList, err := repo.GetUserWatchList(ctx, userId)
	if err != nil {
		server.SendInternalServerError(c, err, server.WithMessage("Failed to retrieve watch list"))
		return
	}

	if err := tx.Commit(ctx); err != nil {
		server.SendInternalServerError(c, err, server.WithMessage("Failed to commit transaction"))
		return
	}

	if watchList == nil {
		server.SendNotFound(c, errors.New("You don't have anything in your watch list"))
		return
	}

	server.SendSuccess(c, watchList, server.WithMessage("WatchList retrieved successfully"))
}
