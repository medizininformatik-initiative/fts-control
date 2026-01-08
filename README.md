# ftsctl

A command-line tool for controlling the FTSnext system.

## Installation

```bash
go build -o ftsctl
```

## Configuration

Create a `config.yaml` in the working directory:

```yaml
api:
  base_url: "http://localhost:8080"
```

The base URL can also be overridden via CLI flag: `--base-url` or `-u`

## Usage

```bash
ftsctl [command] [subcommand] [flags]
```

### Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--verbose` | `-v` | Enable verbose output |
| `--base-url` | `-u` | Override base API URL |

### Commands

#### Project

```bash
ftsctl project list              # List all projects
ftsctl project config -n NAME    # Show project configuration
```

#### Process

```bash
ftsctl process list              # List all process statuses
ftsctl process status -i ID      # Show specific process status
```

#### Transfer

```bash
ftsctl transfer start -n NAME              # Start transfer for all consented patients
ftsctl transfer start -n NAME --ids a,b,c  # Start transfer for specific patient IDs
```

## Development

```bash
# Run all tests
go test -v ./...

# Run tests for cmd package
go test -v ./cmd
```
