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

	if !strings.Contains(err.Error(), "SSH tunnel failed") && !strings.Contains(err.Error(), "connection refused") && !strings.Contains(err.Error(), "dial tcp") && !strings.Contains(err.Error(), "timeout") {
		t.Errorf("error was not actionable, got: %v", err)
	}
}
