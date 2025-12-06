package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateSecret creates a random secret string (40 hex characters)
// Equivalent to: crypto.randomBytes(20).toString('hex')
func GenerateSecret() string {
	// Create 20 random bytes
	bytes := make([]byte, 20)

	// Fill with random data
	// crypto/rand is cryptographically secure
	_, err := rand.Read(bytes)
	if err != nil {
		// In production, handle this better
		// For now, return empty string on error
		return ""
	}

	// Convert to hex string (40 characters)
	return hex.EncodeToString(bytes)
}

