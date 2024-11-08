package auth

import (
	"fmt"

	"github.com/adeyemialameen04/unwind-be/core/server"
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
// @Router       /auth [post]
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
