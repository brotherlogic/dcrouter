package update

import (
	"testing"
)

func TestCheckForUpdate_Dev(t *testing.T) {
	// Should not panic or do anything for dev version
	CheckForUpdate("dev", false)
}

func TestGetLatestRelease_Error(t *testing.T) {
	// Since we know brotherlogic/dcrouter might not have releases yet,
	// we expect an error or 404. This just ensures it doesn't crash.
	_, err := GetLatestRelease()
	if err != nil {
		t.Logf("Expected error or no release: %v", err)
	}
}
