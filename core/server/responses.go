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

// ResponseOption is a function type that modifies a Response
type ResponseOption func(*Response)

// WithMessage sets a custom message for the response
func WithMessage(message string) ResponseOption {
	return func(r *Response) {
		r.Message = message
	}
}

// getDefaultMessage returns the default message for a given status code
func getDefaultMessage(statusCode int) string {
	switch statusCode {
	case http.StatusOK:
		return "Operation successful"
	case http.StatusBadRequest:
		return "Invalid request"
	case http.StatusUnauthorized:
		return "Unauthorized access"
	case http.StatusForbidden:
		return "Access forbidden"
	case http.StatusNotFound:
		return "Resource not found"
	case http.StatusMethodNotAllowed:
		return "Method not allowed"
	case http.StatusConflict:
		return "Resource conflict"
	case http.StatusUnprocessableEntity:
		return "Validation failed"
	case http.StatusInternalServerError:
		return "Internal server error"
	default:
		return http.StatusText(statusCode)
	}
}

func SendSuccess(c *gin.Context, data interface{}, opts ...ResponseOption) {
	resp := Response{
		Status:  http.StatusText(http.StatusOK),
		Message: getDefaultMessage(http.StatusOK),
		Data:    data,
	}

	// Apply any provided options
	for _, opt := range opts {
		opt(&resp)
	}

	c.JSON(http.StatusOK, resp)
}

func SendCreated(c *gin.Context, data interface{}, opts ...ResponseOption) {
	resp := Response{
		Status:  http.StatusText(http.StatusCreated),
		Message: getDefaultMessage(http.StatusCreated),
		Data:    data,
	}

	// Apply any provided options
	for _, opt := range opts {
		opt(&resp)
	}

	c.JSON(http.StatusCreated, resp)
}

func SendValidationError(c *gin.Context, errors interface{}, opts ...ResponseOption) {
	resp := Response{
		Status:  http.StatusText(http.StatusUnprocessableEntity),
		Message: getDefaultMessage(http.StatusUnprocessableEntity),
		Errors:  errors,
	}

	for _, opt := range opts {
		opt(&resp)
	}

	c.JSON(http.StatusUnprocessableEntity, resp)
}

func SendError(c *gin.Context, statusCode int, errs []error, opts ...ResponseOption) {
	var outputErrors []string
	if errs != nil && len(errs) > 0 {
		outputErrors = make([]string, 0, len(errs))
		for _, err := range errs {
			outputErrors = append(outputErrors, err.Error())
		}
	}

	resp := Response{
		Status:  http.StatusText(statusCode),
		Message: getDefaultMessage(statusCode),
	}

	// Only add Errors field if there are actual errors to include
	if outputErrors != nil && len(outputErrors) > 0 {
		resp.Errors = outputErrors
	}

	for _, opt := range opts {
		opt(&resp)
	}

	c.JSON(statusCode, resp)
}

// Helper function to convert a single error to an error slice
func wrapError(err error) []error {
	if err == nil {
		return nil
	}
	return []error{err}
}

// SendBadRequest sends a JSON response with a status code of 400 (Bad Request)
func SendBadRequest(c *gin.Context, err error, opts ...ResponseOption) {
	SendError(c, http.StatusBadRequest, wrapError(err), opts...)
}

// SendForbidden sends a JSON response with a status code of 403 (Forbidden)
func SendForbidden(c *gin.Context, err error, opts ...ResponseOption) {
	SendError(c, http.StatusForbidden, wrapError(err), opts...)
}

// SendInternalServerError sends a JSON response with a status code of 500 (Internal Server Error)
func SendInternalServerError(c *gin.Context, err error, opts ...ResponseOption) {
	SendError(c, http.StatusInternalServerError, wrapError(err), opts...)
}

// SendMethodNotAllowedError sends a JSON response with a status code of 405 (Method Not Allowed)
func SendMethodNotAllowedError(c *gin.Context, err error, opts ...ResponseOption) {
	SendError(c, http.StatusMethodNotAllowed, wrapError(err), opts...)
}

// SendNotFound sends a JSON response with a status code of 404 (Not Found)
func SendNotFound(c *gin.Context, err error, opts ...ResponseOption) {
	SendError(c, http.StatusNotFound, wrapError(err), opts...)
}

// SendConflict sends a JSON response with a status code of 409 (Conflict)
func SendConflict(c *gin.Context, err error, opts ...ResponseOption) {
	SendError(c, http.StatusConflict, wrapError(err), opts...)
}

// SendUnauthorized sends a JSON response with a status code of 401 (Unauthorized)
func SendUnauthorized(c *gin.Context, err error, opts ...ResponseOption) {
	SendError(c, http.StatusUnauthorized, wrapError(err), opts...)
}
