package watchlist

import (
	"context"
	"errors"
	"fmt"

	"github.com/adeyemialameen04/unwind-be/core/server"
	"github.com/adeyemialameen04/unwind-be/internal/db/repository"
	"github.com/adeyemialameen04/unwind-be/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/golodash/galidator"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
//	@Success		200	{object}	server.SuccessResponse{data=[]repository.WatchList}
//	@Failure		401	{object}	server.UnauthorizedResponse			"Unauthorized"
//	@Failure		404	{object}	server.NotFoundResponse				"Profile not found"
//	@Failure		500	{object}	server.InternalServerErrorResponse	"Internal server error"
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

// UpdateWatchListStatus godoc
//
//	@Summary		Update WatchList status
//	@Description	Updates the status of a user watch list.
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Security		AccessTokenBearer
//	@Param			Anime	body		repository.UpdateWatchListStatusParams	true	"Updated status"
//	@Param			id	path	int	true	"Watch list ID"
//	@Success		200		{object}	server.SuccessResponse{data=repository.WatchList}
//	@Failure		401		{object}	server.UnauthorizedResponse			"Unauthorized"
//	@Failure		404		{object}	server.NotFoundResponse				"WatchList not found"
//	@Failure		500		{object}	server.InternalServerErrorResponse	"Internal server error"
//	@Router			/user/watch-list/{id} [patch]
func (h *Handler) UpdateWatchListStatus(c *gin.Context) {
	var (
		ctx = context.Background()
		req repository.UpdateWatchListStatusParams
	)

	if err := h.validateUpdateStatusListRequest(c, &req); err != nil {
		return
	}

	tx, err := h.srv.DB.Begin(ctx)
	if err != nil {
		server.SendInternalServerError(c, err, server.WithMessage("Error beginning transaction"))
		return
	}
	defer tx.Rollback(ctx)

	id := c.Param("id")
	parsedId, err := domain.ParseIDFromParams(c, id)
	if err != nil {
		return
	}

	repo := repository.New(tx)
	req.ID = parsedId
	watchList, err := repo.UpdateWatchListStatus(ctx, req)
	if err != nil {
		server.SendInternalServerError(c, err, server.WithMessage("Fail to update to watch list"))
		return
	}

	if err := tx.Commit(ctx); err != nil {
		server.SendInternalServerError(c, err, server.WithMessage("Failed to commit transaction"))
		return
	}

	server.SendSuccess(c, watchList, server.WithMessage("Updated watchList status successfully"))
}

// UpdateWatchListStatus godoc
//
//	@Summary		Delete Item from watchList.
//	@Description	Deletes an item from a user watch list.
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Security		AccessTokenBearer
//	@Param			id	path	int	true	"Watch list ID"
//	@Success		200		{object}	server.SuccessResponse{data=repository.WatchList}
//	@Failure		401		{object}	server.UnauthorizedResponse			"Unauthorized"
//	@Failure		404		{object}	server.NotFoundResponse				"WatchList not found"
//	@Failure		500		{object}	server.InternalServerErrorResponse	"Internal server error"
//	@Router			/user/watch-list/{id} [delete]
func (h *Handler) DeleteFromWatchList(c *gin.Context) {
	ctx := context.Background()

	tx, err := h.srv.DB.Begin(ctx)
	if err != nil {
		server.SendInternalServerError(c, err, server.WithMessage("Error beginning transaction"))
		return
	}
	defer tx.Rollback(ctx)

	id := c.Param("id")
	parsedId, err := domain.ParseIDFromParams(c, id)
	if err != nil {
		return
	}

	repo := repository.New(tx)
	watchList, err := repo.DeleteWatchList(ctx, parsedId)
	if err != nil {
		if err == pgx.ErrNoRows {
			server.SendNotFound(c, errors.New("Watch list item not found"))
			return
		}
		server.SendInternalServerError(c, err, server.WithMessage("Fail to delete item from watch list"))
		return
	}

	if err := tx.Commit(ctx); err != nil {
		server.SendInternalServerError(c, err, server.WithMessage("Failed to commit transaction"))
		return
	}

	server.SendSuccess(c, watchList, server.WithMessage("Deleted item from watchList successfully"))
}

// AddToList godoc
//
//	@Summary		Adds an anime to a user watch list
//	@Description	Retrieves a user watch list.
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Security		AccessTokenBearer
//
//	@Param			Anime	body		repository.AddToListParams	true	"Anime Data"
//
//	@Success		201		{object}	server.CreatedResponse{data=repository.WatchList}
//	@Failure		401		{object}	server.UnauthorizedResponse			"Unauthorized"
//	@Failure		404		{object}	server.NotFoundResponse				""
//	@Failure		500		{object}	server.InternalServerErrorResponse	"Internal server error"
//	@Router			/user/watch-list [post]
func (h *Handler) AddToList(c *gin.Context) {
	var (
		ctx = context.Background()
		req repository.AddToListParams
	)

	if err := h.validateAddtoListRequest(c, &req); err != nil {
		return
	}

	userId, err := domain.GetUserIDFromContext(c)
	if err != nil {
		return
	}

	tx, err := h.srv.DB.Begin(ctx)
	if err != nil {
		server.SendInternalServerError(c, err, server.WithMessage("Error beginning transaction"))
		return
	}
	defer tx.Rollback(ctx)

	req.UserID = userId
	repo := repository.New(tx)
	watchList, err := repo.AddToList(ctx, req)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == domain.UniqueViolation {
			server.SendConflict(c, err, server.WithMessage("This is already in your watch list"))
			return
		}

		server.SendInternalServerError(c, err, server.WithMessage("Fail to Add to watch list"))
		return
	}

	if err := tx.Commit(ctx); err != nil {
		server.SendInternalServerError(c, err, server.WithMessage("Failed to commit transaction"))
		return
	}

	server.SendCreated(c, watchList, server.WithMessage("Added to watchList successfully"))
}

// GetWatchListByMediaID godoc
//
//	@Summary		Get WatchList using media id
//	@Description	Retrieves a user watch list by it's media ID.
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"Media ID"
//	@Security		AccessTokenBearer
//	@Success		200	{object}	server.SuccessResponse{data=repository.WatchList}
//	@Failure		401	{object}	server.UnauthorizedResponse			"Unauthorized"
//	@Failure		404	{object}	server.NotFoundResponse				"Watchlist not found"
//	@Failure		500	{object}	server.InternalServerErrorResponse	"Internal server error"
//	@Router			/user/watch-list/{id} [get]
func (h *Handler) GetWatchListByMediaID(c *gin.Context) {
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

	mediaId := c.Param("id")

	repo := repository.New(tx)
	watchListItem, err := repo.GetWatchListByMediaID(ctx, repository.GetWatchListByMediaIDParams{
		MediaID: &mediaId,
		UserID:  userId,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			server.SendNotFound(c, errors.New("Not found"))
			return
		}
		server.SendInternalServerError(c, err, server.WithMessage("Failed to retrieve watch list"))
		return
	}

	if err := tx.Commit(ctx); err != nil {
		server.SendInternalServerError(c, err, server.WithMessage("Failed to commit transaction"))
		return
	}

	fmt.Println(watchListItem)
	fmt.Println(watchListItem.UserID.String(), "The one gotten")
	fmt.Println(userId.String(), "The one in context")

	server.SendSuccess(c, watchListItem, server.WithMessage("WatchList retrieved successfully"))
}

func (h *Handler) validateAddtoListRequest(c *gin.Context, req *repository.AddToListParams) error {
	g := galidator.New().CustomMessages(galidator.Messages{
		"required": "$field is required",
	})
	customizer := g.Validator(repository.AddToListParams{})

	if err := c.ShouldBindJSON(req); err != nil {
		server.SendValidationError(c, customizer.DecryptErrors(err))
		return err
	}
	return nil
}

func (h *Handler) validateUpdateStatusListRequest(c *gin.Context, req *repository.UpdateWatchListStatusParams) error {
	g := galidator.New().CustomMessages(galidator.Messages{
		"required": "$field is required",
	})
	customizer := g.Validator(repository.UpdateWatchListStatusParams{})

	if err := c.ShouldBindJSON(req); err != nil {
		server.SendValidationError(c, customizer.DecryptErrors(err))
		return err
	}
	return nil
}
