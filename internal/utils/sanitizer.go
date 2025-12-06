package utils

import (
	"github.com/microcosm-cc/bluemonday"
)

// Global sanitizer instance (thread-safe)
var (
	// strictPolicy removes ALL HTML tags
	strictPolicy = bluemonday.StrictPolicy()

	// ugcPolicy allows safe HTML (links, formatting) - User Generated Content
	ugcPolicy = bluemonday.UGCPolicy()
)

// SanitizeStrict removes ALL HTML tags from input
// Use for: author names, emails, IDs
func SanitizeStrict(input string) string {
	return strictPolicy.Sanitize(input)
}

// SanitizeComment allows safe HTML in comments
// Allows: <b>, <i>, <a>, <p>, <br>, <ul>, <li>, etc.
// Removes: <script>, <style>, onclick, etc.
func SanitizeComment(input string) string {
	return ugcPolicy.Sanitize(input)
}

// Example of what gets sanitized:
//
// Input:  "<script>alert('xss')</script><b>Hello</b>"
// Strict: "Hello"
// UGC:    "<b>Hello</b>"
//
// Input:  "<a href='javascript:alert(1)'>Click</a>"
// Strict: "Click"
// UGC:    "Click" (dangerous href removed)

