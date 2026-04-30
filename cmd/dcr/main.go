package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/brotherlogic/dcrouter/internal/config"
	"github.com/brotherlogic/dcrouter/internal/engine"
	"github.com/brotherlogic/dcrouter/internal/ssh"
	"github.com/brotherlogic/dcrouter/internal/update"
)

var (
	version = "dev"
)

func main() {
	update.CheckForUpdate(version, false)

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "config":
		handleConfig()
	case "version":
		fmt.Printf("dcr version %s\n", version)
	case "update":
		update.CheckForUpdate(version, true)
	case "help":
		printUsage()
	default:
		containerName := os.Args[1]
		cfg, err := config.ReadConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading config: %v\n", err)
			os.Exit(1)
		}

		if cfg.RouterAddress == "" || cfg.HostAddress == "" {
			fmt.Fprintf(os.Stderr, "Configuration is incomplete. Please run 'dcr config' first.\n")
			os.Exit(1)
		}

		port, err := engine.ResolvePort(containerName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error resolving container %q: %v\n", containerName, err)
			os.Exit(1)
		}

		err = ssh.Execute(&ssh.SystemExecutor{}, cfg, port)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error executing SSH: %v\n", err)
			os.Exit(1)
		}
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  dcr config [--router <address>] [--host <address>]  - View or update configuration")
	fmt.Println("  dcr version                                         - Show version information")
	fmt.Println("  dcr update                                          - Check for updates")
	fmt.Println("  dcr <container_name>                                - Connect to a devcontainer")
}

func handleConfig() {
	configCmd := flag.NewFlagSet("config", flag.ExitOnError)
	routerAddr := configCmd.String("router", "", "Set the router address")
	hostAddr := configCmd.String("host", "", "Set the host address")

	err := configCmd.Parse(os.Args[2:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	cfg, err := config.ReadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading config: %v\n", err)
		os.Exit(1)
	}

	updated := false
	if *routerAddr != "" {
		cfg.RouterAddress = *routerAddr
		updated = true
	}
	if *hostAddr != "" {
		cfg.HostAddress = *hostAddr
		updated = true
	}

	if updated {
		if err := config.WriteConfig(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing config: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Configuration updated.")
	} else {
		path, _ := config.GetConfigPath()
		fmt.Printf("Config path: %s\n", path)
		fmt.Printf("Router Address: %s\n", cfg.RouterAddress)
		fmt.Printf("Host Address:   %s\n", cfg.HostAddress)
	}
}
