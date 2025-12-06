package utils

import "strings"

// CleanEmail trims whitespace from an email address
// Equivalent to your cleanEmail in Node.js
func CleanEmail(email string) string {
	return strings.TrimSpace(email)
}

