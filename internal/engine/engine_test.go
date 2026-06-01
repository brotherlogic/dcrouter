package engine

import (
	"testing"
)

func TestResolveWorkspace(t *testing.T) {
	tests := []struct {
		name        string
		issueNum    string
		expected    string
		expectError bool
	}{
		{
			name:        "music",
			issueNum:    "",
			expected:    "music.devpod",
			expectError: false,
		},
		{
			name:        "music",
			issueNum:    "2",
			expected:    "music-2.devpod",
			expectError: false,
		},
		{
			name:        "music",
			issueNum:    "abc",
			expected:    "",
			expectError: true,
		},
	}

	for _, tc := range tests {
		resolved, err := ResolveWorkspace(tc.name, tc.issueNum)
		if tc.expectError {
			if err == nil {
				t.Errorf("Expected error for name=%q issueNum=%q, but got none", tc.name, tc.issueNum)
			}
		} else {
			if err != nil {
				t.Fatalf("Unexpected error for name=%q issueNum=%q: %v", tc.name, tc.issueNum, err)
			}
			if resolved != tc.expected {
				t.Errorf("Expected %q, got %q for name=%q issueNum=%q", tc.expected, resolved, tc.name, tc.issueNum)
			}
		}
	}
}
