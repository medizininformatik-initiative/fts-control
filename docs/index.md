# ftsctl

Command-line tool for controlling the FTSnext system.

## Quick Links

- [Installation](./getting-started/installation.md)
- [Configuration](./getting-started/configuration.md)
- [CLI Commands](./getting-started/commands.md)

## What is ftsctl?

ftsctl lets you manage FTSnext from the command line:

- **List projects** - See all available projects
- **View project config** - Inspect project settings
- **Monitor transfers** - Check process status
- **Start transfers** - Initiate data transfers

## Getting Started

```bash
# Install
go build -o ftsctl

# Configure
ftsctl config set-base-url https://your-ftsnext-server.org

# List projects
ftsctl project list
```
