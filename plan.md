# Dev Container Router (dcrouter) Project Plan

This document outlines the tasks required to implement the `dc` CLI and the necessary changes to the `devcontainer-manager` ecosystem.

## 1. `devcontainer-manager` Requirements
The following changes are needed in the [brotherlogic/devcontainer-manager](https://github.com/brotherlogic/devcontainer-manager) project to support discovery.

- [ ] **Port Allocation Strategy**: Define a unique SSH port for each managed devcontainer (e.g., starting from 2222).
- [ ] **Port Advertising**: Update the manager to generate/update a `mappings.json` file in the repository.
    - Format: `{"containers": {"music": {"port": 2222}, "api": {"port": 2223}}}`
- [ ] **Container Configuration**: Ensure each devcontainer is started with the correct port mapping (e.g., `-p <host_port>:22`) and has an SSH server running.

## 2. `dcrouter` CLI Implementation (Go)
The CLI will be a Go binary named `dc`.

### Phase 1: Research & Setup
- [x] Initialize Go module.
- [x] Implement configuration handling (store `router_address` and `host_address` in `~/.config/dcrouter/config.json`).

### Phase 2: Discovery Engine
- [x] Implement GitHub Fetcher: Pull `mappings.json` from the `devcontainer-manager` repo.
- [x] Implement Caching: Store the mappings locally with a TTL (e.g., 5-10 minutes) to avoid rate limits.
- [x] Implement Name Resolution: Match the user input (e.g., `dc music`) against the mapping.

### Phase 3: SSH Wrapper
- [x] Construct the nested SSH command: `ssh -t <user>@<router> "ssh -t <user>@<host> -p <port>"`
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
- [x] **Mock Test**: Verify the CLI correctly parses a mock `mappings.json`.
- [x] **Integration Test**: Verify the generated SSH command string matches the expected nested format.
- [x] **End-to-End**: Test connection from a remote machine through the router into a test container.
