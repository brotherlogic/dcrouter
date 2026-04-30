package ssh

import (
	"strings"
	"testing"

	"github.com/brotherlogic/dcrouter/internal/config"
)

type MockExecutor struct {
	Argv0 string
	Argv  []string
	Envv  []string
	Err   error
}

func (m *MockExecutor) Exec(argv0 string, argv []string, envv []string) error {
	m.Argv0 = argv0
	m.Argv = argv
	m.Envv = envv
	return m.Err
}

func (m *MockExecutor) LookPath(file string) (string, error) {
	return "/usr/bin/" + file, nil
}

func TestExecute(t *testing.T) {
	mock := &MockExecutor{}
	cfg := &config.Config{
		RouterAddress: "router.example.com",
		HostAddress:   "host.example.com",
	}
	port := 1234

	err := Execute(mock, cfg, port)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// We can't easily mock user.Current() without more abstraction, 
	// but we can check if the arguments contain what we expect.
	// Since user.Current() will return the actual user running the test,
	// we'll check for the presence of the host and router addresses.

	expectedArgv0 := "/usr/bin/ssh"
	if mock.Argv0 != expectedArgv0 {
		t.Errorf("Expected Argv0 %q, got %q", expectedArgv0, mock.Argv0)
	}

	if len(mock.Argv) != 4 {
		t.Errorf("Expected 4 arguments, got %d: %v", len(mock.Argv), mock.Argv)
	}

	// Check router address in Argv[2]
	expectedRouterPart := "@router.example.com"
	foundRouter := false
	if mock.Argv[1] == "-t" {
		if len(mock.Argv) > 2 && strings.Contains(mock.Argv[2], expectedRouterPart) {
			foundRouter = true
		}
	}
	if !foundRouter {
		t.Errorf("Could not find router address in Argv: %v", mock.Argv)
	}

	// Check inner command in Argv[3]
	expectedInnerPart := "@host.example.com -p 1234"
	if !strings.Contains(mock.Argv[3], expectedInnerPart) {
		t.Errorf("Could not find host and port in inner command: %q", mock.Argv[3])
	}
}
