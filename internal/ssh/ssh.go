package ssh

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"syscall"

	"github.com/brotherlogic/dcrouter/internal/config"
)

// Executor defines the interface for executing commands
type Executor interface {
	Exec(argv0 string, argv []string, envv []string) error
	LookPath(file string) (string, error)
}

// SystemExecutor implements Executor using syscall and os/exec
type SystemExecutor struct{}

// Exec wraps syscall.Exec
func (s *SystemExecutor) Exec(argv0 string, argv []string, envv []string) error {
	return syscall.Exec(argv0, argv, envv)
}

// LookPath wraps exec.LookPath
func (s *SystemExecutor) LookPath(file string) (string, error) {
	return exec.LookPath(file)
}

// Execute replaces the current process with a nested SSH command
func Execute(e Executor, cfg *config.Config, port int) error {
	currentUser, err := user.Current()
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}
	username := currentUser.Username

	sshPath, err := e.LookPath("ssh")
	if err != nil {
		return fmt.Errorf("ssh executable not found: %w", err)
	}

	// Construct the inner SSH command
	innerCmd := fmt.Sprintf("ssh -t %s@%s -p %d", username, cfg.HostAddress, port)

	// Arguments for the outer SSH command
	// ssh -t <user>@<router> "ssh -t <user>@<host> -p <port>"
	args := []string{
		"ssh",
		"-t",
		fmt.Sprintf("%s@%s", username, cfg.RouterAddress),
		innerCmd,
	}

	// syscall.Exec requires the path to the executable, the arguments, and the environment
	err = e.Exec(sshPath, args, os.Environ())
	if err != nil {
		return fmt.Errorf("failed to execute ssh: %w", err)
	}

	return nil
}
