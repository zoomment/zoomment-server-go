package utils

import "testing"

func TestCleanName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "keeps letters and spaces",
			input:    "John Doe",
			expected: "John Doe",
		},
		{
			name:     "keeps numbers",
			input:    "User123",
			expected: "User123",
		},
		{
			name:     "keeps dots and underscores",
			input:    "user.name_test",
			expected: "user.name_test",
		},
		{
			name:     "removes special characters",
			input:    "John<script>Doe",
			expected: "JohnscriptDoe",
		},
		{
			name:     "handles unicode letters",
			input:    "Тигран Симонян",
			expected: "Тигран Симонян",
		},
		{
			name:     "trims whitespace",
			input:    "  John  ",
			expected: "John",
		},
		{
			name:     "handles empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CleanName(tt.input)
			if result != tt.expected {
				t.Errorf("CleanName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCleanEmail(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "trims whitespace",
			input:    "  test@example.com  ",
			expected: "test@example.com",
		},
		{
			name:     "handles normal email",
			input:    "test@example.com",
			expected: "test@example.com",
		},
		{
			name:     "handles empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CleanEmail(tt.input)
			if result != tt.expected {
				t.Errorf("CleanEmail(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

