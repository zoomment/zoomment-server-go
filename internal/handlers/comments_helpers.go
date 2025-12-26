package handlers

import (
	"github.com/gin-gonic/gin"
	"zoomment-server/internal/models"
)

// commentToResponse converts a Comment model to a response map with _id
// Includes secret for newly created comments (used in AddComment)
// Note: Node.js doesn't have owner in schema, but some legacy data might have it
func commentToResponse(comment *models.Comment) gin.H {
	return gin.H{
		"_id":        comment.ID.Hex(),
		"parentId":   comment.ParentID,
		"author":     comment.Author,
		"email":      comment.Email,
		"gravatar":   comment.Gravatar,
		"body":       comment.Body,
		"domain":     comment.Domain,
		"pageUrl":    comment.PageURL,
		"pageId":     comment.PageID,
		"isVerified": comment.IsVerified,
		"secret":     comment.Secret,
		"createdAt":  comment.CreatedAt,
		"updatedAt":  comment.UpdatedAt,
		"isOwn":      comment.Secret != "",
	}
}

// commentsToResponse converts a slice of comments to response format
func commentsToResponse(comments []models.Comment) []gin.H {
	response := make([]gin.H, 0, len(comments))
	for _, comment := range comments {
		response = append(response, commentToResponse(&comment))
	}
	return response
}
