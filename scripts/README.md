# AgbCloud CLI Scripts

This directory contains scripts for the AgbCloud CLI project.

## Test Scripts

### test.sh

The main test runner script used by the Makefile for running unit and integration tests.

#### Usage

```bash
# Run unit tests only
./scripts/test.sh --unit-only

# Run integration tests only
./scripts/test.sh --integration-only

# Run all tests
./scripts/test.sh --all

# Run tests with verbose output
./scripts/test.sh --unit-only --verbose
```

## Build and Release Scripts

### generate-server-files.sh

Generates server files and installation commands for distribution. This script is used in the CI/CD pipeline to create PowerShell installation scripts and version metadata.

#### Usage

```bash
# Generate files for a specific version
VERSION="v1.2.3" ./scripts/generate-server-files.sh

# Generate files with default dev version
./scripts/generate-server-files.sh
```

#### Generated Files

- `server-files/install.ps1` - PowerShell installation script
- `server-files/latest.json` - Version metadata for API

### upload-to-oss.sh

Uploads build artifacts to Alibaba Cloud OSS. Can be used independently or as part of the CI/CD pipeline.

#### Usage

```bash
# Set required environment variables
export OSS_ACCESS_KEY_ID="your-access-key"
export OSS_ACCESS_KEY_SECRET="your-secret-key"
export VERSION="v1.2.3"

# Upload packages
./scripts/upload-to-oss.sh
```

## Features

- **Test Automation**: Comprehensive test runner with multiple modes
- **Build Integration**: Automatic generation of installation files for CI/CD
- **Cloud Storage**: Upload artifacts to OSS with proper permissions
- **Version Management**: Support for version-specific builds and releases

## Compatibility

- Linux/macOS build environments
- Alibaba Cloud OSS integration
- Go 1.23+ test framework 