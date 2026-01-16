# Configuration

ftsctl stores configuration in `~/.config/ftsctl/config.yaml` (Linux/macOS) or `%APPDATA%\ftsctl\config.yaml` (Windows).

## Set the API URL

```bash
ftsctl config set-base-url https://your-ftsnext-server.org
```

## View Current Configuration

```bash
ftsctl config show
```

## Configuration File

You can also edit the config file directly:

```yaml
api:
  base_url: "https://your-ftsnext-server.org"
```

## Authentication

ftsctl supports three authentication methods. Configure only one at a time.

### Basic Authentication

```yaml
api:
  base_url: "https://your-ftsnext-server.org"

auth:
  basic:
    user: "your-username"
    password: "${FTS_PASSWORD}"  # Use environment variable
```

### OAuth2 (Client Credentials)

```yaml
api:
  base_url: "https://your-ftsnext-server.org"

auth:
  oauth2:
    token_url: "https://keycloak.example.org/realms/myrealm/protocol/openid-connect/token"
    client_id: "ftsctl"
    client_secret: "${FTS_CLIENT_SECRET}"
    scope: "openid"  # optional
```

### Certificate Authentication (mTLS)

```yaml
api:
  base_url: "https://your-ftsnext-server.org"  # Must be HTTPS

auth:
  certificate:
    cert_file: "/path/to/client.crt"
    key_file: "/path/to/client.key"
    ca_file: "/path/to/ca.crt"  # optional
```

## Environment Variables

Secrets in the config file support `${VAR}` syntax for environment variables:

```yaml
auth:
  basic:
    user: "admin"
    password: "${FTS_PASSWORD}"
```

Then set the environment variable:

```bash
export FTS_PASSWORD="your-secret-password"
```

## Command Line Override

Override the base URL for a single command:

```bash
ftsctl --base-url https://other-server.org project list
```

## Next Steps

- [Commands](./commands.md) - Available commands
