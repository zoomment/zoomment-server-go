package validators

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// RegisterCustomValidators registers custom validation rules
// Call this once at startup
func RegisterCustomValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// Register custom validators here if needed
		// Example: v.RegisterValidation("customrule", customRuleFunc)
		_ = v
	}
}

// AddCommentRequest with comprehensive validation
type AddCommentRequest struct {
	// URL must be valid and required
	PageURL string `json:"pageUrl" binding:"required,url,max=2000"`

	// PageID is required, max 500 chars
	PageID string `json:"pageId" binding:"required,max=500"`

	// Body: required, 1-10000 characters
	Body string `json:"body" binding:"required,min=1,max=10000"`

	// Author: required, 1-100 characters
	Author string `json:"author" binding:"required,min=1,max=100"`

	// Email: required, must be valid email
	Email string `json:"email" binding:"required,email,max=254"`

	// ParentID: optional, if provided must be 24 chars (MongoDB ObjectID)
	ParentID *string `json:"parentId" binding:"omitempty,len=24"`
}

// AuthRequest with email validation
type AuthRequest struct {
	Email string `json:"email" binding:"required,email,max=254"`
}

// AddSiteRequest with URL validation
type AddSiteRequest struct {
	URL string `json:"url" binding:"required,url,max=2000"`
}

// AddReactionRequest with validation
type AddReactionRequest struct {
	PageID   string `json:"pageId" binding:"required,max=500"`
	Reaction string `json:"reaction" binding:"required,min=1,max=20"`
}

// ValidationErrorResponse formats validation errors nicely
type ValidationErrorResponse struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// FormatValidationErrors converts validator errors to user-friendly messages
func FormatValidationErrors(err error) []ValidationErrorResponse {
	var errors []ValidationErrorResponse

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			var message string
			switch e.Tag() {
			case "required":
				message = "This field is required"
			case "email":
				message = "Invalid email format"
			case "url":
				message = "Invalid URL format"
			case "min":
				message = "Value is too short"
			case "max":
				message = "Value is too long"
			case "len":
				message = "Invalid length"
			default:
				message = "Invalid value"
			}

			errors = append(errors, ValidationErrorResponse{
				Field:   e.Field(),
				Message: message,
			})
		}
	}

	return errors
}

