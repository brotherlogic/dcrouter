# Dev Container Router (dcrouter) Project Plan

This document outlines the tasks required to implement the `dcr` CLI and the necessary changes to the `devcontainer-manager` ecosystem.

## 1. `devcontainer-manager` Requirements
The following changes are needed in the [brotherlogic/devcontainer-manager](https://github.com/brotherlogic/devcontainer-manager) project to support discovery.

- [ ] **Workspace Tracking**: Ensure DevPod workspaces are named consistently with the project name.
- [ ] **Host Configuration**: Ensure DevPod is installed and configured on the host machine to allow `ssh <project>.devpod` connections.

## 2. `dcrouter` CLI Implementation (Go)
The CLI will be a Go binary named `dcr`.

### Phase 1: Research & Setup
- [x] Initialize Go module.
- [x] Implement configuration handling (store `router_address` and `host_address` in `~/.config/dcrouter/config.json`).

### Phase 2: Name Resolution
- [x] Implement Convention-based Resolution: Resolve `<name>` to `<name>.devpod`.

### Phase 3: SSH Wrapper
- [x] Construct the nested SSH command: `ssh -t <user>@<router> "ssh -t <user>@<host> 'ssh -t <workspace>'"`
- [x] Use `syscall.Exec` to replace the Go process with the SSH process. This ensures terminal signals (resize, Ctrl+C) are handled correctly by the system SSH client.

### Phase 4: Self-Update
- [x] Implement an automatic update check against GitHub releases (standard for `brotherlogic` projects).

## 3. Development Workflow
- [ ] **Branching Strategy**: All features developed in branches (e.g., `feat/ssh-wrapper`).
- [ ] **CI Pipeline (GitHub Actions)**:
    - Trigger on push to any branch.
    - Run `go test ./...`, `go vet`, and `staticcheck`.
    - Cross-compile for `linux/amd64` and `darwin/arm64` to verify build integrity.
- [ ] **Automated PRs**: Use GitHub Actions to automatically create PRs for feature branches.
- [ ] **Releases (GoReleaser)**:
    - Trigger on tag push (e.g., `v1.0.0`).
    - Automatically build binaries and create GitHub Releases.

## 4. Infrastructure Prep
- [ ] Ensure the Router machine has the necessary SSH private keys to access the Host machine.
- [ ] Ensure the Host machine's `authorized_keys` includes the Router's public key.

## 5. Verification & Testing
- [x] **Mocking Strategy**: Abstract SSH execution behind an interface to allow CI to verify command construction without live network access.
- [x] **Convention Test**: Verify the CLI correctly resolves the workspace name based on the naming convention.
- [x] **Integration Test**: Verify the generated SSH command string matches the expected nested format.
- [x] **End-to-End**: Test connection from a remote machine through the router into a test container.
