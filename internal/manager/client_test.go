package manager

import (
	"strings"
	"testing"

	"github.com/brotherlogic/dcrouter/internal/config"
)

func TestFetchContainers_SSHFailure(t *testing.T) {
	cfg := &config.Config{
		RouterAddress: "invalid-router.local",
		HostAddress:   "invalid-host.local",
	}

	_, err := FetchContainers(cfg)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	// The error should now include stderr from the SSH command which fails to resolve the hostname
	if !strings.Contains(err.Error(), "SSH process exited unexpectedly") && !strings.Contains(err.Error(), "stderr:") {
		t.Errorf("error was not actionable, expected SSH failure details, got: %v", err)
	}
}
