# Installation

## Build from Source

```bash
git clone https://github.com/medizininformatik-initiative/fts-control.git
cd fts-control
go build -o ftsctl
```

Move to a directory in your PATH:

```bash
# With sudo
sudo mv ftsctl /usr/local/bin/

# Without sudo
mkdir -p ~/.local/bin
mv ftsctl ~/.local/bin/
export PATH="$HOME/.local/bin:$PATH"  # Add to ~/.bashrc or ~/.zshrc
```

## Verify

```bash
ftsctl --help
```

## Next Steps

- [Configuration](./configuration.md) - Set up your connection
- [Commands](./commands.md) - Available commands
