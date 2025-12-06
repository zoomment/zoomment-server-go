package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

// GenerateGravatar creates an MD5 hash of an email for Gravatar
// Gravatar uses MD5 hash of lowercase email as identifier
func GenerateGravatar(email string) string {
	// Lowercase and trim the email
	email = strings.ToLower(strings.TrimSpace(email))

	// Create MD5 hash
	hash := md5.Sum([]byte(email))

	// Convert to hex string
	return hex.EncodeToString(hash[:])
}

