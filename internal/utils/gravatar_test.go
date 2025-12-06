package utils

import "testing"

func TestGenerateGravatar(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected string
	}{
		{
			name:     "generates correct hash",
			email:    "test@example.com",
			expected: "55502f40dc8b7c769880b10874abc9d0", // MD5 of "test@example.com"
		},
		{
			name:     "lowercase email",
			email:    "TEST@EXAMPLE.COM",
			expected: "55502f40dc8b7c769880b10874abc9d0", // Same hash (case-insensitive)
		},
		{
			name:     "trims whitespace",
			email:    "  test@example.com  ",
			expected: "55502f40dc8b7c769880b10874abc9d0",
		},
		{
			name:     "empty email",
			email:    "",
			expected: "d41d8cd98f00b204e9800998ecf8427e", // MD5 of empty string
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateGravatar(tt.email)
			if result != tt.expected {
				t.Errorf("GenerateGravatar(%q) = %q, want %q", tt.email, result, tt.expected)
			}
		})
	}
}

