package manager

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"os/exec"
	"os/user"
	"strings"
	"time"

	"github.com/brotherlogic/dcrouter/internal/config"
	pb "github.com/brotherlogic/devcontainer-manager/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// FetchContainers fetches devcontainers from the remote host using an SSH tunnel
func FetchContainers(cfg *config.Config) ([]*pb.DevcontainerConfig, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("could not get current user: %w", err)
	}

	// Find a free local port
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return nil, fmt.Errorf("could not find a free local port: %w", err)
	}
	localPort := listener.Addr().(*net.TCPAddr).Port
	listener.Close()

	// Establish SSH tunnel
	jumpArg := fmt.Sprintf("%s@%s", usr.Username, cfg.RouterAddress)
	hostArg := fmt.Sprintf("%s@%s", usr.Username, cfg.HostAddress)
	routerTunnelArg := fmt.Sprintf("%d:localhost:%d", localPort, localPort)
	hostTunnelArg := fmt.Sprintf("%d:localhost:50051", localPort)

	var stderr bytes.Buffer
	cmd := exec.Command("ssh", "-L", routerTunnelArg, jumpArg, "ssh", "-N", "-L", hostTunnelArg, hostArg)
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("SSH tunnel failed to start: %w", err)
	}
	defer func() {
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
	}()

	errCh := make(chan error, 1)
	go func() {
		errCh <- cmd.Wait()
	}()

	// Wait for the local port to become dialable (max 5 seconds)
	target := fmt.Sprintf("localhost:%d", localPort)
	dialed := false
	for i := 0; i < 50; i++ {
		select {
		case err := <-errCh:
			errMsg := "SSH process exited unexpectedly"
			if err != nil {
				errMsg += fmt.Sprintf(" (%v)", err)
			}
			if stderr.Len() > 0 {
				errMsg += fmt.Sprintf(": %s", strings.TrimSpace(stderr.String()))
			}
			return nil, fmt.Errorf("%s", errMsg)
		default:
		}

		conn, err := net.DialTimeout("tcp", target, 100*time.Millisecond)
		if err == nil {
			conn.Close()
			dialed = true
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	if !dialed {
		errMsg := fmt.Sprintf("SSH tunnel failed: timeout waiting for local port %d to be dialable", localPort)
		if stderr.Len() > 0 {
			errMsg += fmt.Sprintf(", stderr: %s", strings.TrimSpace(stderr.String()))
		}
		return nil, fmt.Errorf("%s", errMsg)
	}

	// Make gRPC call
	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("gRPC dial failed: %w", err)
	}
	defer conn.Close()

	client := pb.NewManagerServiceClient(conn)
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.List(ctx, &pb.ListRequest{})
	if err != nil {
		return nil, fmt.Errorf("gRPC list failed: %w", err)
	}

	return resp.GetConfigs(), nil
}
