package utils

import (
	"regexp"
	"strings"
)

// CleanName removes special characters from a name
// Allows letters (any language), spaces, digits, dots, and underscores
// Equivalent to your cleanName in Node.js
func CleanName(input string) string {
	// In Go, we use regexp package for regular expressions
	// This pattern matches letters (Unicode), digits, spaces, dots, underscores
	// (?i) makes it case-insensitive
	re := regexp.MustCompile(`[^\p{L}\s\d._]`)

	// Replace all non-matching characters with empty string
	cleaned := re.ReplaceAllString(input, "")

	return strings.TrimSpace(cleaned)
}

