package handlers

import (
	"net/url"
	"time"

	"zoomment-server/internal/models"
)

// ========================================
// Response Types
// ========================================

// SiteResponse is the JSON response format for sites
type SiteResponse struct {
	ID        string    `json:"_id"`
	UserID    string    `json:"userId"`
	Domain    string    `json:"domain"`
	Verified  bool      `json:"verified"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// CommentResponse is the JSON response format for newly created comments
type CommentResponse struct {
	ID         string    `json:"_id"`
	ParentID   *string   `json:"parentId"`
	Author     string    `json:"author"`
	Email      string    `json:"email"`
	Gravatar   string    `json:"gravatar"`
	Body       string    `json:"body"`
	Domain     string    `json:"domain"`
	PageURL    string    `json:"pageUrl"`
	PageID     string    `json:"pageId"`
	IsVerified bool      `json:"isVerified"`
	Secret     string    `json:"secret"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	IsOwn      bool      `json:"isOwn"`
}

// DeletedResponse is the JSON response for deleted resources
type DeletedResponse struct {
	ID string `json:"_id"`
}

// VisitorCountResponse is the JSON response for visitor count
type VisitorCountResponse struct {
	PageID string `json:"pageId"`
	Count  int64  `json:"count"`
}

// PageViewCount represents visitor count for a page
type PageViewCount struct {
	PageID string `bson:"pageId" json:"pageId"`
	Count  int    `bson:"count" json:"count"`
}

// DomainPagesResponse is the JSON response for domain page views
type DomainPagesResponse struct {
	Domain string          `json:"domain"`
	Pages  []PageViewCount `json:"pages"`
}

// PaginatedCommentsResponse is the JSON response for paginated comments by site
type PaginatedCommentsResponse struct {
	Comments []CommentResponse `json:"comments"`
	Total    int64             `json:"total"`
	Limit    int               `json:"limit"`
	Skip     int               `json:"skip"`
	HasMore  bool              `json:"hasMore"`
}

// UserProfileResponse is the JSON response for user profile
type UserProfileResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// MessageResponse is a generic message response
type MessageResponse struct {
	Message string `json:"message"`
}

// ========================================
// Model to Response Converters
// ========================================

// SiteToResponse converts a Site model to response format
func SiteToResponse(site *models.Site) SiteResponse {
	return SiteResponse{
		ID:        site.ID.Hex(),
		UserID:    site.UserID.Hex(),
		Domain:    site.Domain,
		Verified:  site.Verified,
		CreatedAt: site.CreatedAt,
		UpdatedAt: site.UpdatedAt,
	}
}

// SitesToResponse converts a slice of sites to response format
func SitesToResponse(sites []models.Site) []SiteResponse {
	result := make([]SiteResponse, 0, len(sites))
	for i := range sites {
		result = append(result, SiteToResponse(&sites[i]))
	}
	return result
}

// CommentToResponse converts a Comment model to response format
func CommentToResponse(comment *models.Comment) CommentResponse {
	return CommentResponse{
		ID:         comment.ID.Hex(),
		ParentID:   comment.ParentID,
		Author:     comment.Author,
		Email:      comment.Email,
		Gravatar:   comment.Gravatar,
		Body:       comment.Body,
		Domain:     comment.Domain,
		PageURL:    comment.PageURL,
		PageID:     comment.PageID,
		IsVerified: comment.IsVerified,
		Secret:     comment.Secret,
		CreatedAt:  comment.CreatedAt,
		UpdatedAt:  comment.UpdatedAt,
		IsOwn:      true, // New comments are always "own"
	}
}

// CommentsToResponse converts a slice of comments to response format
func CommentsToResponse(comments []models.Comment) []CommentResponse {
	result := make([]CommentResponse, 0, len(comments))
	for i := range comments {
		result = append(result, CommentToResponse(&comments[i]))
	}
	return result
}

// NewDeletedResponse creates a deleted response
func NewDeletedResponse(id string) DeletedResponse {
	return DeletedResponse{ID: id}
}

// NewVisitorCountResponse creates a visitor count response
func NewVisitorCountResponse(pageID string, count int64) VisitorCountResponse {
	return VisitorCountResponse{PageID: pageID, Count: count}
}

// NewDomainPagesResponse creates a domain pages response
func NewDomainPagesResponse(domain string, pages []PageViewCount) DomainPagesResponse {
	return DomainPagesResponse{Domain: domain, Pages: pages}
}

// NewPaginatedCommentsResponse creates a paginated comments response
func NewPaginatedCommentsResponse(comments []models.Comment, total int64, limit, skip int) PaginatedCommentsResponse {
	return PaginatedCommentsResponse{
		Comments: CommentsToResponse(comments),
		Total:    total,
		Limit:    limit,
		Skip:     skip,
		HasMore:  int64(skip+len(comments)) < total,
	}
}

// ========================================
// Common Helpers
// ========================================

// ExtractDomainFromPageID parses a pageId and returns the domain
func ExtractDomainFromPageID(pageID string) (string, error) {
	parsedURL, err := url.Parse("https://" + pageID)
	if err != nil {
		return "", err
	}
	return parsedURL.Hostname(), nil
}
