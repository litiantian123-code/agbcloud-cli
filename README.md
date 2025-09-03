# AgbCloud CLI

A command-line interface for AgbCloud services.

## Installation

### From Source

```bash
git clone https://github.com/liyuebing/agbcloud-cli.git
cd agbcloud-cli
make build
```

### From Release

Download the latest release for your platform from the [releases page](https://github.com/liyuebing/agbcloud-cli/releases).

## Usage

```bash
# Show help
agbcloud --help

# Show version
agbcloud version

# Configuration management
agbcloud config list
agbcloud config get api_key
agbcloud config set api_key your-key-here
```

## Environment Variables

- `AGB_API_KEY`: Your AgbCloud API key

## Configuration

The CLI uses the following default configuration:

- **API Endpoint**: `https://agb.cloud`
- **API Key**: Retrieved from `AGB_API_KEY` environment variable

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