# AgbCloud CLI

A command-line interface for AgbCloud services.

## Installation

### From Source

```bash
git clone https://github.com/agbcloud/agbcloud-cli.git
cd agbcloud-cli
make build
```

### From Release

Download the latest release for your platform from the [releases page](https://github.com/agbcloud/agbcloud-cli/releases).

## Usage

### Basic Commands

```bash
# Show help
agbcloud --help

# Show version
agbcloud version

# Log in to AgbCloud
agbcloud login

# Configuration management
agbcloud config list
agbcloud config get endpoint
agbcloud config set endpoint your-endpoint-here

# Set configuration via environment variables (recommended)
export AGB_CLI_ENDPOINT=agb.cloud  # Domain only, https:// added automatically
```

### SSL Certificate Verification

The CLI automatically determines SSL verification behavior based on the endpoint:

**SSL Verification Enabled (Secure)**:
- Production domains (e.g., `agb.cloud`, `api.example.com`)
- Standard HTTPS port (443)

**SSL Verification Disabled (Development)**:
- IP addresses (e.g., `12.34.56.78`, `[2001:db8::1]`)
- Localhost (`localhost`, `127.0.0.1`)
- Development domains (`.local`, `.dev`, `.test`, `.internal`)
- Non-standard ports (e.g., `:8080`, `:8443`)

**Manual Override**:
```bash
# Force skip SSL verification
export AGB_CLI_SKIP_SSL_VERIFY=true

# Force SSL verification (even for IP addresses)
export AGB_CLI_SKIP_SSL_VERIFY=false
```

### Authentication

```bash
# Log in using OAuth in browser
agbcloud login
```

## Environment Variables

- `AGB_CLI_ENDPOINT`: AgbCloud API endpoint domain (default: agb.cloud, https:// prefix added automatically)
- `AGB_CLI_SKIP_SSL_VERIFY`: Override SSL verification behavior ("true" to skip, "false" to enforce)

## Configuration

The CLI uses the following configuration:

- **API Endpoint**: 
  - `AGB_CLI_ENDPOINT` environment variable (domain only, https:// added automatically)
  - Default: `agb.cloud`
- **Authentication**: 
  - OAuth token-based authentication via `agbcloud login`
- **SSL Verification**: 
  - Automatic based on endpoint type (IP addresses, localhost, .local/.dev/.test/.internal domains, non-standard ports skip verification)
  - Override with `AGB_CLI_SKIP_SSL_VERIFY` environment variable

## Development

### Prerequisites

- Go 1.23 or later
- Make

### Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Clean build artifacts
make clean
```

### Testing

```bash
make test
```

### Formatting

```bash
make fmt
```

## Cross-Platform Support

This CLI supports the following platforms:

- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64, arm64)

## Project Structure

```
agbcloud-cli/
├── cmd/                 # Command definitions
│   ├── version.go      # Version command
│   └── config.go       # Configuration commands
├── internal/           # Internal packages
│   ├── config/         # Configuration management
│   └── client/         # API client
├── pkg/               # Public packages
│   └── version/       # Version information
├── .github/           # GitHub Actions
│   └── workflows/
├── main.go           # Entry point
├── Makefile         # Build scripts
├── go.mod          # Go module
├── .gitignore      # Git ignore file
├── LICENSE         # Apache-2.0 license
└── README.md       # This file
```

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request 