package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/brotherlogic/dcrouter/internal/config"
	"github.com/brotherlogic/dcrouter/internal/manager"
	pb "github.com/brotherlogic/devcontainer-manager/proto"
)

func handleList() {
	cfg, err := config.ReadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading config: %v\n", err)
		os.Exit(1)
	}

	if cfg.RouterAddress == "" || cfg.HostAddress == "" {
		fmt.Fprintf(os.Stderr, "Configuration is incomplete. Please run 'dcr config' first.\n")
		os.Exit(1)
	}

	containers, err := manager.FetchContainers(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching containers: %v\n", err)
		os.Exit(1)
	}

	for _, c := range containers {
		formatContainer(os.Stdout, c)
	}
}

func formatContainer(w io.Writer, config *pb.DevcontainerConfig) {
	issueStr := ""
	if config.Request != nil && config.Request.Identifier != nil {
		if issue, ok := config.Request.Identifier.Id.(*pb.Identifier_IssueNumber); ok {
			issueStr = fmt.Sprintf("[%d]", issue.IssueNumber)
		}
	}

	stateStr := ""
	if config.State != pb.State_DCM_READY {
		stateStr = config.State.String()
		stateStr = strings.ToLower(strings.TrimPrefix(stateStr, "State_DCM_"))
		// Handle legacy prefix if needed or simple strings
		stateStr = strings.ToLower(strings.TrimPrefix(stateStr, "dcm_"))
	}

	fmt.Fprintf(w, "%s %s %s\n", config.Id, issueStr, stateStr)
}
