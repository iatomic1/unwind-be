package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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

// Base Response struct with common fields
type Response struct {
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
}

// Specific response types for different status codes
type SuccessResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Status  string      `json:"status" example:"OK"`
	Message string      `json:"message,omitempty"`
}

type CreatedResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Status  string      `json:"status" example:"Created"`
	Message string      `json:"message,omitempty"`
}

type ValidationErrorResponse struct {
	Errors  interface{} `json:"errors,omitempty"`
	Status  string      `json:"status" example:"Unprocessable Entity"`
	Message string      `json:"message,omitempty"`
}

type BadRequestResponse struct {
	Errors  interface{} `json:"errors,omitempty"`
	Status  string      `json:"status" example:"Bad Request"`
	Message string      `json:"message,omitempty"`
}

type UnauthorizedResponse struct {
	Errors  interface{} `json:"errors,omitempty"`
	Status  string      `json:"status" example:"Unauthorized"`
	Message string      `json:"message,omitempty"`
}

type ForbiddenResponse struct {
	Errors  interface{} `json:"errors,omitempty"`
	Status  string      `json:"status" example:"Forbidden"`
	Message string      `json:"message,omitempty"`
}

type NotFoundResponse struct {
	Errors  interface{} `json:"errors,omitempty"`
	Status  string      `json:"status" example:"Not Found"`
	Message string      `json:"message,omitempty"`
}

type MethodNotAllowedResponse struct {
	Errors  interface{} `json:"errors,omitempty"`
	Status  string      `json:"status" example:"Method Not Allowed"`
	Message string      `json:"message,omitempty"`
}

type ConflictResponse struct {
	Errors  interface{} `json:"errors,omitempty"`
	Status  string      `json:"status" example:"Conflict"`
	Message string      `json:"message,omitempty"`
}

type InternalServerErrorResponse struct {
	Errors  interface{} `json:"errors,omitempty"`
	Status  string      `json:"status" example:"Internal Server Error"`
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

func SendSuccess(c *gin.Context, data interface{}, opts ...ResponseOption) {
	resp := SuccessResponse{
		Status:  http.StatusText(http.StatusOK),
		Message: getDefaultMessage(http.StatusOK),
		Data:    data,
	}
	c.JSON(http.StatusOK, resp)
}

func SendCreated(c *gin.Context, data interface{}, opts ...ResponseOption) {
	resp := CreatedResponse{
		Status:  http.StatusText(http.StatusCreated),
		Message: getDefaultMessage(http.StatusCreated),
		Data:    data,
	}
	c.JSON(http.StatusCreated, resp)
}

func SendValidationError(c *gin.Context, errors interface{}, opts ...ResponseOption) {
	resp := ValidationErrorResponse{
		Status:  http.StatusText(http.StatusUnprocessableEntity),
		Message: getDefaultMessage(http.StatusUnprocessableEntity),
		Errors:  errors,
	}
	c.JSON(http.StatusUnprocessableEntity, resp)
}

func SendError(c *gin.Context, statusCode int, errs []error, opts ...ResponseOption) {
	var outputErrors []string
	if len(errs) > 0 {
		outputErrors = make([]string, 0, len(errs))
		for _, err := range errs {
			outputErrors = append(outputErrors, err.Error())
		}
	}

	switch statusCode {
	case http.StatusBadRequest:
		resp := BadRequestResponse{
			Status:  http.StatusText(statusCode),
			Message: getDefaultMessage(statusCode),
			Errors:  outputErrors,
		}
		c.JSON(statusCode, resp)
	case http.StatusUnauthorized:
		resp := UnauthorizedResponse{
			Status:  http.StatusText(statusCode),
			Message: getDefaultMessage(statusCode),
			Errors:  outputErrors,
		}
		c.JSON(statusCode, resp)
	case http.StatusForbidden:
		resp := ForbiddenResponse{
			Status:  http.StatusText(statusCode),
			Message: getDefaultMessage(statusCode),
			Errors:  outputErrors,
		}
		c.JSON(statusCode, resp)
	case http.StatusNotFound:
		resp := NotFoundResponse{
			Status:  http.StatusText(statusCode),
			Message: getDefaultMessage(statusCode),
			Errors:  outputErrors,
		}
		c.JSON(statusCode, resp)
	case http.StatusMethodNotAllowed:
		resp := MethodNotAllowedResponse{
			Status:  http.StatusText(statusCode),
			Message: getDefaultMessage(statusCode),
			Errors:  outputErrors,
		}
		c.JSON(statusCode, resp)
	case http.StatusConflict:
		resp := ConflictResponse{
			Status:  http.StatusText(statusCode),
			Message: getDefaultMessage(statusCode),
			Errors:  outputErrors,
		}
		c.JSON(statusCode, resp)
	case http.StatusInternalServerError:
		resp := InternalServerErrorResponse{
			Status:  http.StatusText(statusCode),
			Message: getDefaultMessage(statusCode),
			Errors:  outputErrors,
		}
		c.JSON(statusCode, resp)
	default:
		// Fallback to generic Response for unknown status codes
		resp := Response{
			Status:  http.StatusText(statusCode),
			Message: getDefaultMessage(statusCode),
			Errors:  outputErrors,
		}
		c.JSON(statusCode, resp)
	}
}

// Helper functions remain the same, they'll use SendError which now uses the appropriate response type
func wrapError(err error) []error {
	if err == nil {
		return nil
	}
	return []error{err}
}

func SendBadRequest(c *gin.Context, err error, opts ...ResponseOption) {
	SendError(c, http.StatusBadRequest, wrapError(err), opts...)
}

func SendUnauthorized(c *gin.Context, err error, opts ...ResponseOption) {
	SendError(c, http.StatusUnauthorized, wrapError(err), opts...)
}

func SendForbidden(c *gin.Context, err error, opts ...ResponseOption) {
	SendError(c, http.StatusForbidden, wrapError(err), opts...)
}

func SendNotFound(c *gin.Context, err error, opts ...ResponseOption) {
	SendError(c, http.StatusNotFound, wrapError(err), opts...)
}

func SendMethodNotAllowedError(c *gin.Context, err error, opts ...ResponseOption) {
	SendError(c, http.StatusMethodNotAllowed, wrapError(err), opts...)
}

func SendConflict(c *gin.Context, err error, opts ...ResponseOption) {
	SendError(c, http.StatusConflict, wrapError(err), opts...)
}

func SendInternalServerError(c *gin.Context, err error, opts ...ResponseOption) {
	SendError(c, http.StatusInternalServerError, wrapError(err), opts...)
}
