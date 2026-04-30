package engine

import (
	"testing"
)

func TestResolveWorkspace(t *testing.T) {
	name := "music"
	expected := "music.devpod"
	
	resolved, err := ResolveWorkspace(name)
	if err != nil {
		t.Fatalf("ResolveWorkspace failed: %v", err)
	}

	if resolved != expected {
		t.Errorf("Expected %q, got %q", expected, resolved)
	}
}
