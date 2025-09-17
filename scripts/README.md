# AgbCloud CLI Installation Scripts

This directory contains installation scripts for the AgbCloud CLI.

## Windows Installation Scripts

### install-windows-simple.ps1

The main Windows installation script with enhanced functionality.

#### Usage

```powershell
# Install latest version
.\install-windows-simple.ps1

# Install specific version
.\install-windows-simple.ps1 -Version v1.2.3

# Show help
.\install-windows-simple.ps1 -Help

# Custom installation directory
.\install-windows-simple.ps1 -InstallPath "C:\Tools\agbcloud"

# Custom architecture (auto-detected by default)
.\install-windows-simple.ps1 -Architecture arm64

# Custom download URL
.\install-windows-simple.ps1 -DownloadUrl "https://custom-server.com"
```

#### Parameters

- `-Version`: Specify version to install (e.g., 'v1.0.0', 'latest')
- `-Architecture`: Specify architecture ('amd64' or 'arm64')
- `-InstallPath`: Specify custom installation directory
- `-DownloadUrl`: Specify custom download base URL
- `-Help`: Show help message

#### Environment Variables

You can also use environment variables instead of parameters:

- `AGBCLOUD_VERSION`: Default version to install
- `AGBCLOUD_PATH`: Default installation directory
- `AGBCLOUD_DOWNLOAD_URL`: Default download base URL

#### Examples

```powershell
# Install latest version
.\install-windows-simple.ps1

# Install specific version
.\install-windows-simple.ps1 -Version v1.2.3

# Install to custom directory
.\install-windows-simple.ps1 -InstallPath "D:\MyTools\agbcloud"

# Use environment variable for version
$env:AGBCLOUD_VERSION = "v1.1.0"
.\install-windows-simple.ps1

# Combine parameters
.\install-windows-simple.ps1 -Version v2.0.0 -Architecture amd64 -InstallPath "C:\CLI\agbcloud"
```

### server-files/install.ps1

The server deployment version of the installation script with the same functionality as the main script.

## Build and Release Scripts

### generate-server-files.sh

Generates server files and installation commands for distribution. This script is typically run after a successful build to provide installation instructions.

#### Usage

```bash
# Generate files for a specific version
VERSION="v1.2.3" ./scripts/generate-server-files.sh

# Generate files with default dev version
./scripts/generate-server-files.sh
```

#### Generated Installation Commands

The script generates two types of installation commands for Windows users:

**1. Latest Version Installation (Production)**
```powershell
powershell -Command "irm https://agbcloud-internal.oss-cn-hangzhou.aliyuncs.com/install.ps1 | iex"
```

**2. Specific Version Installation (Testing)**
```powershell
powershell -Command "irm https://agbcloud-internal.oss-cn-hangzhou.aliyuncs.com/install.ps1 | iex" -Version v1.2.3
```

#### Additional Installation Options

The generated commands also include examples for:
- Custom installation directory
- Specific architecture selection
- Help information

#### Testing Team Workflow

1. **For Production Testing**: Use the latest version command
2. **For Version-Specific Testing**: Use the specific version command with the exact version number
3. **For Custom Scenarios**: Use additional parameters as needed

Example build notification output:
```
[DOC] Windows Installation Commands:

[REFRESH] Install Latest Version (recommended for production):
   powershell -Command "irm https://agbcloud-internal.oss-cn-hangzhou.aliyuncs.com/install.ps1 | iex"

[TARGET] Install Specific Version v1.2.3 (for testing):
   powershell -Command "irm https://agbcloud-internal.oss-cn-hangzhou.aliyuncs.com/install.ps1 | iex" -Version v1.2.3

[TIP] Testing Team Usage:
   Use the specific version command above to test version v1.2.3
   Use the latest version command for general testing
```

## Features

- **Version Selection**: Install specific versions or latest
- **Architecture Detection**: Automatic detection with manual override
- **Custom Installation Path**: Install to any directory
- **Environment Variable Support**: Use environment variables for configuration
- **No Emoji Output**: Compatible with all Windows terminals (cmd, PowerShell)
- **Comprehensive Help**: Built-in help system with examples
- **Error Handling**: Robust error handling and fallback mechanisms
- **PATH Management**: Automatic PATH updates with manual fallback instructions
- **Build Integration**: Automatic generation of installation commands for CI/CD

## Compatibility

- Windows 10/11
- PowerShell 5.1+
- Windows PowerShell and PowerShell Core
- Command Prompt (cmd) - no emoji display issues
- Windows Terminal
- Constrained Language Mode support 