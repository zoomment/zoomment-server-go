// Package docs provides Swagger documentation for Zoomment API
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "swagger": "2.0",
    "info": {
        "title": "Zoomment API",
        "description": "Open Source Self-Hosted Comment System API",
        "version": "1.0.0",
        "contact": {
            "name": "Zoomment",
            "url": "https://github.com/zoomment"
        },
        "license": {
            "name": "MIT"
        }
    },
    "host": "localhost:8080",
    "basePath": "/api",
    "schemes": ["http", "https"],
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "in": "header",
            "name": "token"
        },
        "FingerprintAuth": {
            "type": "apiKey",
            "in": "header",
            "name": "fingerprint",
            "description": "Browser fingerprint for anonymous tracking"
        }
    },
    "paths": {
        "/comments": {
            "get": {
                "summary": "List comments",
                "description": "Get comments for a page or domain with pagination",
                "tags": ["Comments"],
                "parameters": [
                    {
                        "name": "pageId",
                        "in": "query",
                        "type": "string",
                        "description": "Page identifier"
                    },
                    {
                        "name": "domain",
                        "in": "query",
                        "type": "string",
                        "description": "Domain name"
                    },
                    {
                        "name": "limit",
                        "in": "query",
                        "type": "integer",
                        "default": 10,
                        "description": "Number of comments to return (max 50)"
                    },
                    {
                        "name": "skip",
                        "in": "query",
                        "type": "integer",
                        "default": 0,
                        "description": "Number of comments to skip"
                    },
                    {
                        "name": "sort",
                        "in": "query",
                        "type": "string",
                        "enum": ["asc", "desc"],
                        "default": "asc",
                        "description": "Sort order by date: asc (oldest first) or desc (newest first)"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Paginated list of comments",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "comments": {
                                    "type": "array",
                                    "items": {
                                        "$ref": "#/definitions/Comment"
                                    }
                                },
                                "total": {
                                    "type": "integer"
                                },
                                "limit": {
                                    "type": "integer"
                                },
                                "skip": {
                                    "type": "integer"
                                },
                                "hasMore": {
                                    "type": "boolean"
                                }
                            }
                        }
                    },
                    "400": {
                        "description": "Bad request - pageId or domain required"
                    }
                }
            },
            "post": {
                "summary": "Add comment",
                "description": "Create a new comment",
                "tags": ["Comments"],
                "parameters": [
                    {
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "required": ["pageUrl", "pageId", "body", "author", "email"],
                            "properties": {
                                "pageUrl": {"type": "string"},
                                "pageId": {"type": "string"},
                                "body": {"type": "string"},
                                "author": {"type": "string"},
                                "email": {"type": "string"},
                                "parentId": {"type": "string"}
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Comment created",
                        "schema": {
                            "$ref": "#/definitions/Comment"
                        }
                    },
                    "400": {
                        "description": "Validation error"
                    }
                }
            }
        },
        "/comments/{id}": {
            "delete": {
                "summary": "Delete comment",
                "description": "Delete a comment by ID",
                "tags": ["Comments"],
                "security": [{"ApiKeyAuth": []}],
                "parameters": [
                    {
                        "name": "id",
                        "in": "path",
                        "required": true,
                        "type": "string"
                    },
                    {
                        "name": "secret",
                        "in": "query",
                        "type": "string",
                        "description": "Secret for guest deletion"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Comment deleted",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "_id": {
                                    "type": "string",
                                    "description": "Deleted comment ID"
                                }
                            }
                        }
                    },
                    "403": {
                        "description": "Forbidden"
                    },
                    "404": {
                        "description": "Comment not found"
                    }
                }
            }
        },
        "/comments/{commentId}/replies": {
            "get": {
                "summary": "Get replies",
                "description": "Get replies for a specific comment with pagination",
                "tags": ["Comments"],
                "parameters": [
                    {
                        "name": "commentId",
                        "in": "path",
                        "required": true,
                        "type": "string"
                    },
                    {
                        "name": "limit",
                        "in": "query",
                        "type": "integer",
                        "default": 10,
                        "description": "Number of replies to return"
                    },
                    {
                        "name": "skip",
                        "in": "query",
                        "type": "integer",
                        "default": 0,
                        "description": "Number of replies to skip"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Paginated list of replies",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "replies": {
                                    "type": "array",
                                    "items": {
                                        "$ref": "#/definitions/Comment"
                                    }
                                },
                                "total": {
                                    "type": "integer"
                                },
                                "limit": {
                                    "type": "integer"
                                },
                                "skip": {
                                    "type": "integer"
                                },
                                "hasMore": {
                                    "type": "boolean"
                                }
                            }
                        }
                    }
                }
            }
        },
        "/users/auth": {
            "post": {
                "summary": "Request magic link",
                "description": "Send magic link email for authentication",
                "tags": ["Users"],
                "parameters": [
                    {
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "required": ["email"],
                            "properties": {
                                "email": {"type": "string", "format": "email"}
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Magic link sent"
                    },
                    "400": {
                        "description": "Invalid email"
                    }
                }
            }
        },
        "/users/profile": {
            "get": {
                "summary": "Get profile",
                "description": "Get current user's profile",
                "tags": ["Users"],
                "security": [{"ApiKeyAuth": []}],
                "responses": {
                    "200": {
                        "description": "User profile",
                        "schema": {
                            "$ref": "#/definitions/User"
                        }
                    },
                    "403": {
                        "description": "Forbidden"
                    }
                }
            }
        },
        "/users": {
            "delete": {
                "summary": "Delete account",
                "description": "Delete current user's account and their sites",
                "tags": ["Users"],
                "security": [{"ApiKeyAuth": []}],
                "responses": {
                    "200": {
                        "description": "Account deleted"
                    },
                    "403": {
                        "description": "Forbidden"
                    }
                }
            }
        },
        "/sites": {
            "get": {
                "summary": "List sites",
                "description": "Get all sites for current user",
                "tags": ["Sites"],
                "security": [{"ApiKeyAuth": []}],
                "responses": {
                    "200": {
                        "description": "List of sites",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/Site"
                            }
                        }
                    }
                }
            },
            "post": {
                "summary": "Add site",
                "description": "Register a new site",
                "tags": ["Sites"],
                "security": [{"ApiKeyAuth": []}],
                "parameters": [
                    {
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "required": ["url"],
                            "properties": {
                                "url": {"type": "string", "format": "url"}
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Site created",
                        "schema": {
                            "$ref": "#/definitions/Site"
                        }
                    },
                    "404": {
                        "description": "Meta tag not found"
                    },
                    "409": {
                        "description": "Site already exists"
                    }
                }
            }
        },
        "/sites/{id}": {
            "delete": {
                "summary": "Delete site",
                "description": "Delete a site by ID",
                "tags": ["Sites"],
                "security": [{"ApiKeyAuth": []}],
                "parameters": [
                    {
                        "name": "id",
                        "in": "path",
                        "required": true,
                        "type": "string"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Site deleted"
                    },
                    "404": {
                        "description": "Site not found"
                    }
                }
            }
        },
        "/reactions": {
            "get": {
                "summary": "Get reactions",
                "description": "Get reactions for a page",
                "tags": ["Reactions"],
                "parameters": [
                    {
                        "name": "pageId",
                        "in": "query",
                        "required": true,
                        "type": "string"
                    },
                    {
                        "name": "fingerprint",
                        "in": "header",
                        "type": "string"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Reactions data",
                        "schema": {
                            "$ref": "#/definitions/Reaction"
                        }
                    }
                }
            },
            "post": {
                "summary": "Add reaction",
                "description": "Add or toggle a reaction",
                "tags": ["Reactions"],
                "security": [{"FingerprintAuth": []}],
                "parameters": [
                    {
                        "name": "fingerprint",
                        "in": "header",
                        "required": true,
                        "type": "string"
                    },
                    {
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "required": ["pageId", "reaction"],
                            "properties": {
                                "pageId": {"type": "string"},
                                "reaction": {"type": "string"}
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Reaction updated",
                        "schema": {
                            "$ref": "#/definitions/Reaction"
                        }
                    }
                }
            }
        },
        "/visitors": {
            "get": {
                "summary": "Get visitor count",
                "description": "Get visitor count for a page",
                "tags": ["Visitors"],
                "parameters": [
                    {
                        "name": "pageId",
                        "in": "query",
                        "required": true,
                        "type": "string",
                        "description": "Page identifier"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Visitor count",
                        "schema": {
                            "$ref": "#/definitions/Visitor"
                        }
                    },
                    "400": {
                        "description": "Bad request - pageId is required"
                    }
                }
            },
            "post": {
                "summary": "Track visitor",
                "description": "Track a page visitor (requires fingerprint)",
                "tags": ["Visitors"],
                "security": [{"FingerprintAuth": []}],
                "parameters": [
                    {
                        "name": "fingerprint",
                        "in": "header",
                        "required": true,
                        "type": "string"
                    },
                    {
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "required": ["pageId"],
                            "properties": {
                                "pageId": {"type": "string"}
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Visitor tracked",
                        "schema": {
                            "$ref": "#/definitions/Visitor"
                        }
                    },
                    "400": {
                        "description": "Fingerprint required"
                    }
                }
            }
        },
        "/visitors/domain": {
            "get": {
                "summary": "Get visitors by domain",
                "description": "Get page view counts grouped by pageId for a domain",
                "tags": ["Visitors"],
                "parameters": [
                    {
                        "name": "domain",
                        "in": "query",
                        "required": true,
                        "type": "string"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Page view counts by pageId",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "domain": {"type": "string"},
                                "pages": {
                                    "type": "array",
                                    "items": {
                                        "type": "object",
                                        "properties": {
                                            "pageId": {"type": "string"},
                                            "count": {"type": "integer"}
                                        }
                                    }
                                }
                            }
                        }
                    }
                }
            }
        },
        "/votes": {
            "post": {
                "summary": "Vote on comment",
                "description": "Vote on a comment (upvote/downvote). Toggle off if same vote, update if opposite.",
                "tags": ["Votes"],
                "security": [{"FingerprintAuth": []}],
                "parameters": [
                    {
                        "name": "fingerprint",
                        "in": "header",
                        "required": true,
                        "type": "string"
                    },
                    {
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/VoteInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Vote recorded",
                        "schema": {
                            "$ref": "#/definitions/Vote"
                        }
                    },
                    "400": {
                        "description": "Validation error"
                    },
                    "404": {
                        "description": "Comment not found"
                    }
                }
            },
            "get": {
                "summary": "Get votes bulk",
                "description": "Get votes for multiple comments",
                "tags": ["Votes"],
                "parameters": [
                    {
                        "name": "commentIds",
                        "in": "query",
                        "required": true,
                        "type": "string",
                        "description": "Comma-separated comment IDs"
                    },
                    {
                        "name": "fingerprint",
                        "in": "header",
                        "type": "string"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Vote counts per comment",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "$ref": "#/definitions/Vote"
                            }
                        }
                    }
                }
            }
        },
        "/votes/{commentId}": {
            "get": {
                "summary": "Get votes for comment",
                "description": "Get vote counts for a single comment",
                "tags": ["Votes"],
                "parameters": [
                    {
                        "name": "commentId",
                        "in": "path",
                        "required": true,
                        "type": "string"
                    },
                    {
                        "name": "fingerprint",
                        "in": "header",
                        "type": "string"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Vote counts",
                        "schema": {
                            "$ref": "#/definitions/Vote"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "Comment": {
            "type": "object",
            "properties": {
                "_id": {
                    "type": "string",
                    "description": "Comment ID"
                },
                "parentId": {
                    "type": "string",
                    "description": "Parent comment ID for replies"
                },
                "author": {
                    "type": "string",
                    "description": "Comment author name"
                },
                "gravatar": {
                    "type": "string",
                    "description": "Gravatar hash"
                },
                "body": {
                    "type": "string",
                    "description": "Comment body (HTML)"
                },
                "isVerified": {
                    "type": "boolean",
                    "description": "Whether comment author is verified"
                },
                "isOwn": {
                    "type": "boolean",
                    "description": "Whether comment belongs to current user"
                },
                "owner": {
                    "type": "object",
                    "properties": {
                        "name": {"type": "string"},
                        "gravatar": {"type": "string"}
                    }
                },
                "createdAt": {
                    "type": "string",
                    "format": "date-time"
                },
                "repliesCount": {
                    "type": "integer",
                    "description": "Number of replies to this comment"
                }
            }
        },
        "User": {
            "type": "object",
            "properties": {
                "id": {"type": "string"},
                "name": {"type": "string"},
                "email": {"type": "string", "format": "email"}
            }
        },
        "Site": {
            "type": "object",
            "properties": {
                "_id": {"type": "string"},
                "domain": {"type": "string"},
                "verified": {"type": "boolean"},
                "createdAt": {"type": "string", "format": "date-time"}
            }
        },
        "Reaction": {
            "type": "object",
            "properties": {
                "aggregation": {
                    "type": "array",
                    "items": {
                        "type": "object",
                        "properties": {
                            "_id": {"type": "string", "description": "Emoji"},
                            "count": {"type": "integer"}
                        }
                    }
                },
                "userReaction": {
                    "type": "string",
                    "description": "Current user's reaction emoji"
                }
            }
        },
        "Visitor": {
            "type": "object",
            "properties": {
                "pageId": {"type": "string"},
                "count": {"type": "integer"}
            }
        },
        "Vote": {
            "type": "object",
            "properties": {
                "commentId": {"type": "string"},
                "upvotes": {"type": "integer"},
                "downvotes": {"type": "integer"},
                "score": {"type": "integer", "description": "upvotes - downvotes"},
                "userVote": {"type": "integer", "enum": [-1, 0, 1], "description": "Current user's vote"}
            }
        },
        "VoteInput": {
            "type": "object",
            "required": ["commentId", "value"],
            "properties": {
                "commentId": {"type": "string"},
                "value": {"type": "integer", "enum": [1, -1], "description": "1 = upvote, -1 = downvote"}
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0.0",
	Host:             "localhost:8080",
	BasePath:         "/api",
	Schemes:          []string{"http", "https"},
	Title:            "Zoomment API",
	Description:      "Open Source Self-Hosted Comment System API",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
