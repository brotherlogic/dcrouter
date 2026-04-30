package engine

import (
	"fmt"
)

// ResolveWorkspace resolves a container name to a DevPod workspace name
func ResolveWorkspace(name string) (string, error) {
	return fmt.Sprintf("%s.devpod", name), nil
}
