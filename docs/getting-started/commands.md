# CLI Commands

## Global Options

```bash
ftsctl [options] <command>
```

| Option | Short | Description |
|--------|-------|-------------|
| `--verbose` | `-v` | Enable verbose output |
| `--base-url` | `-u` | Override API URL |

## config

Manage CLI configuration.

### config show

Display current configuration.

```bash
ftsctl config show
```

### config set-base-url

Set the API URL.

```bash
ftsctl config set-base-url https://your-ftsnext-server.org
```

## project

Inspect projects.

### project list

List all available projects.

```bash
ftsctl project list
```

### project config

Show configuration for a specific project.

```bash
ftsctl project config -n PROJECT_NAME
```

| Option | Short | Required | Description |
|--------|-------|----------|-------------|
| `--projectName` | `-n` | Yes | Project name |

## process

Monitor transfer processes.

### process list

List all process statuses.

```bash
ftsctl process list
```

### process status

Show status of a specific process.

```bash
ftsctl process status -i PROCESS_ID
```

| Option | Short | Required | Description |
|--------|-------|----------|-------------|
| `--processId` | `-i` | Yes | Process ID |

## transfer

Manage data transfers.

### transfer start

Start a transfer for a project.

```bash
# Transfer all consented patients
ftsctl transfer start -n PROJECT_NAME

# Transfer specific patients
ftsctl transfer start -n PROJECT_NAME --ids patient1,patient2,patient3
```

| Option | Short | Required | Description |
|--------|-------|----------|-------------|
| `--projectName` | `-n` | Yes | Project name |
| `--ids` | `-i` | No | Comma-separated patient IDs |
