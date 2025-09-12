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

## Features

- **Version Selection**: Install specific versions or latest
- **Architecture Detection**: Automatic detection with manual override
- **Custom Installation Path**: Install to any directory
- **Environment Variable Support**: Use environment variables for configuration
- **No Emoji Output**: Compatible with all Windows terminals (cmd, PowerShell)
- **Comprehensive Help**: Built-in help system with examples
- **Error Handling**: Robust error handling and fallback mechanisms
- **PATH Management**: Automatic PATH updates with manual fallback instructions

## Compatibility

- Windows 10/11
- PowerShell 5.1+
- Windows PowerShell and PowerShell Core
- Command Prompt (cmd) - no emoji display issues
- Windows Terminal
- Constrained Language Mode support 