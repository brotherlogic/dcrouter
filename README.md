# dcrouter

`dcrouter` is a CLI tool designed to facilitate SSH connections to devcontainers running on a local network host, jumping through a router machine.

## How it Works

The tool automates a nested SSH connection:
`Local Machine` -> `Router Machine` -> `Host Machine` -> `Devcontainer (Port)`

It resolves container names (e.g., `music`) to specific ports by fetching a mapping from the [brotherlogic/devcontainer-manager](https://github.com/brotherlogic/devcontainer-manager) repository.

## Installation

To install `dcrouter`, you can use `go install`:

```bash
go install github.com/brotherlogic/dcrouter/cmd/dc@latest
```

## Configuration

Before using `dcrouter`, you need to configure the addresses for your router and the host machine where the containers are running.

```bash
dc config --router router.yournetwork.com --host host.yournetwork.com
```

You can view your current configuration by running:

```bash
dc config
```

## Usage

To connect to a devcontainer, simply provide its name:

```bash
dc music
```

This will:
1. Fetch (or use cached) mappings to find the port for `music`.
2. Construct and execute a nested SSH command using `syscall.Exec`.
3. Seamlessly drop you into the container's shell.

### Commands

| Command | Description |
| --- | --- |
| `dc <name>` | Connect to a specific devcontainer |
| `dc config` | View or update router/host configuration |
| `dc version` | Show the current version |
| `dc update` | Manually check for tool updates |
| `dc help` | Show usage information |

## Auto-Updates

`dcrouter` automatically checks for updates against GitHub Releases every 24 hours. If a new version is available, it will notify you. You can also manually trigger a check with `dc update`.

## Architecture Details

- **Shell Integration**: Uses `syscall.Exec` to replace the Go process with the system `ssh` client. This ensures that all terminal signals (window resizing, `Ctrl+C`, etc.) are handled natively by SSH.
- **Caching**: Container mappings are cached locally for 5 minutes to ensure fast response times while staying up-to-date with changes in the manager repository.
- **Identity**: Assumes the same username is used across all machines. SSH keys for the host machine must be present on the router machine.
