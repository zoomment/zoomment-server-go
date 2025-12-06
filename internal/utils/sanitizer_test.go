package utils

import "testing"

// Test functions must start with "Test" and take *testing.T
func TestSanitizeStrict(t *testing.T) {
	// Table-driven tests - a common Go pattern
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "removes script tags",
			input:    "<script>alert('xss')</script>Hello",
			expected: "Hello",
		},
		{
			name:     "removes all HTML",
			input:    "<b>Bold</b> and <i>italic</i>",
			expected: "Bold and italic",
		},
		{
			name:     "handles plain text",
			input:    "Just plain text",
			expected: "Just plain text",
		},
		{
			name:     "removes dangerous attributes",
			input:    `<div onclick="evil()">Content</div>`,
			expected: "Content",
		},
		{
			name:     "handles empty string",
			input:    "",
			expected: "",
		},
	}

	// Run each test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeStrict(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeStrict(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSanitizeComment(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains string // Check if result contains this
		notContains string // Check if result does NOT contain this
	}{
		{
			name:        "allows bold tags",
			input:       "<b>Bold text</b>",
			contains:    "<b>",
			notContains: "",
		},
		{
			name:        "removes script tags",
			input:       "<script>evil()</script>Safe",
			contains:    "Safe",
			notContains: "script",
		},
		{
			name:        "removes onclick",
			input:       `<a onclick="evil()" href="http://example.com">Link</a>`,
			contains:    "href",
			notContains: "onclick",
		},
		{
			name:        "removes javascript href",
			input:       `<a href="javascript:alert(1)">Click</a>`,
			contains:    "Click",
			notContains: "javascript",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeComment(tt.input)

			if tt.contains != "" && !contains(result, tt.contains) {
				t.Errorf("SanitizeComment(%q) = %q, should contain %q", tt.input, result, tt.contains)
			}

			if tt.notContains != "" && contains(result, tt.notContains) {
				t.Errorf("SanitizeComment(%q) = %q, should NOT contain %q", tt.input, result, tt.notContains)
			}
		})
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

