package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
}

func SendSuccess(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Status:  http.StatusText(http.StatusOK),
		Message: message,
		Data:    data,
	})
}

func SendValidationError(c *gin.Context, errors interface{}) {
	c.JSON(http.StatusUnprocessableEntity, Response{
		Status: http.StatusText(http.StatusUnprocessableEntity),
		Errors: errors,
	})
}
