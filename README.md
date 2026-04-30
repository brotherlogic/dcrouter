# dcrouter

`dcrouter` is a CLI tool designed to facilitate SSH connections to devcontainers running on a local network host, jumping through a router machine.

## How it Works

The tool automates a nested SSH connection:
`Local Machine` -> `Router Machine` -> `Host Machine` -> `Devcontainer (Port)`

It resolves container names (e.g., `music`) to specific ports by fetching a mapping from the [brotherlogic/devcontainer-manager](https://github.com/brotherlogic/devcontainer-manager) repository.

## Installation

To install `dcrouter`, you can use `go install`:

```bash
go install github.com/brotherlogic/dcrouter/cmd/dcr@latest
```

## Configuration

Before using `dcrouter`, you need to configure the addresses for your router and the host machine where the containers are running.

```bash
dcr config --router router.yournetwork.com --host host.yournetwork.com
```

You can view your current configuration by running:

```bash
dcr config
```

## Usage

To connect to a devcontainer, simply provide its name:

```bash
dcr music
```

This will:
1. Fetch (or use cached) mappings to find the port for `music`.
2. Construct and execute a nested SSH command using `syscall.Exec`.
3. Seamlessly drop you into the container's shell.

### Commands

| Command | Description |
| --- | --- |
| `dcr <name>` | Connect to a specific devcontainer |
| `dcr config` | View or update router/host configuration |
| `dcr version` | Show the current version |
| `dcr update` | Manually check for tool updates |
| `dcr help` | Show usage information |

## Auto-Updates

`dcrouter` automatically checks for updates against GitHub Releases every 24 hours. If a new version is available, it will notify you. You can also manually trigger a check with `dcr update`.

## Architecture Details

- **Shell Integration**: Uses `syscall.Exec` to replace the Go process with the system `ssh` client. This ensures that all terminal signals (window resizing, `Ctrl+C`, etc.) are handled natively by SSH.
- **Caching**: Container mappings are cached locally for 5 minutes to ensure fast response times while staying up-to-date with changes in the manager repository.
- **Identity**: Assumes the same username is used across all machines. SSH keys for the host machine must be present on the router machine.
