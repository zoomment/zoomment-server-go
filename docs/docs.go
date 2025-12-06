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
        }
    },
    "paths": {
        "/comments": {
            "get": {
                "summary": "List comments",
                "description": "Get comments for a page or domain",
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
                    }
                ],
                "responses": {
                    "200": {
                        "description": "List of comments with replies"
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
                    "201": {
                        "description": "Comment created"
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
                        "description": "Comment deleted"
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
                        "description": "User profile"
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
                        "description": "List of sites"
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
                        "description": "Site created"
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
                        "description": "Reactions data"
                    }
                }
            },
            "post": {
                "summary": "Add reaction",
                "description": "Add or toggle a reaction",
                "tags": ["Reactions"],
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
                        "description": "Reaction updated"
                    }
                }
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

