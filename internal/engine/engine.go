package engine

import (
	"fmt"
	"strconv"
)

// ResolveWorkspace resolves a container name and optional issue number to a DevPod workspace name
func ResolveWorkspace(name string, issueNum string) (string, error) {
	if issueNum != "" {
		if _, err := strconv.Atoi(issueNum); err != nil {
			return "", fmt.Errorf("issue number must be strictly an integer: %w", err)
		}
		return fmt.Sprintf("%s-%s.devpod", name, issueNum), nil
	}
	return fmt.Sprintf("%s.devpod", name), nil
}
