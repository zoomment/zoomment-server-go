package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AppError represents a structured application error
type AppError struct {
	Code       string `json:"code"`             // Machine-readable error code
	Message    string `json:"message"`          // Human-readable message
	StatusCode int    `json:"-"`                // HTTP status code
	Details    any    `json:"details,omitempty"` // Additional details (validation errors, etc.)
}

// Error implements the error interface
func (e *AppError) Error() string {
	return e.Message
}

// Common error codes
const (
	ErrCodeValidation   = "validation_error"
	ErrCodeNotFound     = "not_found"
	ErrCodeForbidden    = "forbidden"
	ErrCodeUnauthorized = "unauthorized"
	ErrCodeDatabase     = "database_error"
	ErrCodeInternal     = "internal_error"
	ErrCodeConflict     = "conflict"
	ErrCodeBadRequest   = "bad_request"
)

// Pre-defined common errors
var (
	ErrNotFound = &AppError{
		Code:       ErrCodeNotFound,
		Message:    "Resource not found",
		StatusCode: http.StatusNotFound,
	}

	ErrForbidden = &AppError{
		Code:       ErrCodeForbidden,
		Message:    "You don't have permission to access this resource",
		StatusCode: http.StatusForbidden,
	}

	ErrUnauthorized = &AppError{
		Code:       ErrCodeUnauthorized,
		Message:    "Authentication required",
		StatusCode: http.StatusUnauthorized,
	}

	ErrDatabaseError = &AppError{
		Code:       ErrCodeDatabase,
		Message:    "A database error occurred",
		StatusCode: http.StatusInternalServerError,
	}

	ErrInternalError = &AppError{
		Code:       ErrCodeInternal,
		Message:    "An internal error occurred",
		StatusCode: http.StatusInternalServerError,
	}
)

// New creates a new AppError
func New(code, message string, statusCode int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}

// ValidationError creates a validation error with details
func ValidationError(message string, details any) *AppError {
	return &AppError{
		Code:       ErrCodeValidation,
		Message:    message,
		StatusCode: http.StatusBadRequest,
		Details:    details,
	}
}

// NotFound creates a not found error with custom message
func NotFound(resource string) *AppError {
	return &AppError{
		Code:       ErrCodeNotFound,
		Message:    resource + " not found",
		StatusCode: http.StatusNotFound,
	}
}

// Conflict creates a conflict error (e.g., duplicate entry)
func Conflict(message string) *AppError {
	return &AppError{
		Code:       ErrCodeConflict,
		Message:    message,
		StatusCode: http.StatusConflict,
	}
}

// BadRequest creates a bad request error
func BadRequest(message string) *AppError {
	return &AppError{
		Code:       ErrCodeBadRequest,
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}

// Response sends an error response to the client
// Format: {"message": "...", "code": "...", "details": ...}
// "message" is always present (Node.js compatible), others are optional
func (e *AppError) Response(c *gin.Context) {
	response := gin.H{"message": e.Message}

	// Add optional fields only if they have values
	if e.Code != "" {
		response["code"] = e.Code
	}
	if e.Details != nil {
		response["details"] = e.Details
	}

	c.JSON(e.StatusCode, response)
}

// Abort sends an error response and aborts the request
func (e *AppError) Abort(c *gin.Context) {
	response := gin.H{"message": e.Message}

	if e.Code != "" {
		response["code"] = e.Code
	}
	if e.Details != nil {
		response["details"] = e.Details
	}

	c.AbortWithStatusJSON(e.StatusCode, response)
}

// HandleError is a helper to respond with an error
// Usage: errors.HandleError(c, err)
func HandleError(c *gin.Context, err error) {
	if appErr, ok := err.(*AppError); ok {
		appErr.Response(c)
		return
	}

	// Unknown error - return internal error
	ErrInternalError.Response(c)
}

