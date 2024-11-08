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

func SendError(c *gin.Context, statusCode int, errs ...error) {
	responseData := Response{
		Status:  http.StatusText(statusCode),
		Message: "error processing request",
	}

	outputErrors := make([]string, 0, len(errs))
	for _, err := range errs {
		outputErrors = append(outputErrors, err.Error())
	}
	responseData.Errors = outputErrors

	c.JSON(statusCode, responseData)
}

func SendBadRequest(c *gin.Context, err error) {
	SendError(c, http.StatusBadRequest, err)
}

// SendForbidden sends a JSON response with a status code of 403 (Forbidden).
func SendForbidden(c *gin.Context, err error) {
	SendError(c, http.StatusForbidden, err)
}

// SendInternalServerError sends a JSON response with a status code of 500 (Internal Server Error).
func SendInternalServerError(c *gin.Context, err error) {
	SendError(c, http.StatusInternalServerError, err)
}

// SendMethodNotAllowedError sends a JSON response with a status code of 405 (Method Not Allowed).
func SendMethodNotAllowedError(c *gin.Context, err error) {
	SendError(c, http.StatusMethodNotAllowed, err)
}

// SendNotFound sends a JSON response with a status code of 404 (Not Found).
func SendNotFound(c *gin.Context, err error) {
	SendError(c, http.StatusNotFound, err)
}

// SendConflict sends a JSON response with a status code of 409 (Unauthorized).
func SendConflict(c *gin.Context, err error) {
	SendError(c, http.StatusConflict, err)
}

// SendUnauthorized sends a JSON response with a status code of 401 (Unauthorized).
func SendUnauthorized(c *gin.Context, err error) {
	SendError(c, http.StatusUnauthorized, err)
}
