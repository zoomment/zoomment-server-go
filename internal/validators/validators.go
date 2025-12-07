// Package validators contains request validation structs for API endpoints.
// All request types are defined here to ensure consistent validation across handlers.
package validators

// AddCommentRequest validates POST /api/comments
type AddCommentRequest struct {
	PageURL  string  `json:"pageUrl" binding:"required,url,max=2000"`
	PageID   string  `json:"pageId" binding:"required,max=500"`
	Body     string  `json:"body" binding:"required,min=1,max=10000"`
	Author   string  `json:"author" binding:"required,min=1,max=100"`
	Email    string  `json:"email" binding:"required,email,max=254"`
	ParentID *string `json:"parentId" binding:"omitempty,len=24"`
}

// AuthRequest validates POST /api/users/auth
type AuthRequest struct {
	Email string `json:"email" binding:"required,email,max=254"`
}

// AddSiteRequest validates POST /api/sites
type AddSiteRequest struct {
	URL string `json:"url" binding:"required,url,max=2000"`
}

// AddReactionRequest validates POST /api/reactions
type AddReactionRequest struct {
	PageID   string `json:"pageId" binding:"required,max=500"`
	Reaction string `json:"reaction" binding:"required,min=1,max=20"`
}

