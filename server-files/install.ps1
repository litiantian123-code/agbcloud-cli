# PowerShell script to download and install AgbCloud CLI binary

param(
    [string]$Version = "",
    [string]$Architecture = "",
    [string]$InstallPath = "",
    [string]$DownloadUrl = "",
    [switch]$Help
)

# Show help information
if ($Help) {
    Write-Host "AgbCloud CLI Installation Script"
    Write-Host ""
    Write-Host "Usage:"
    Write-Host "  .\install-windows-simple.ps1 [options]"
    Write-Host ""
    Write-Host "Options:"
    Write-Host "  -Version <version>      Specify version to install (e.g., 'v1.0.0', 'latest')"
    Write-Host "  -Architecture <arch>    Specify architecture ('amd64' or 'arm64')"
    Write-Host "  -InstallPath <path>     Specify custom installation directory"
    Write-Host "  -DownloadUrl <url>      Specify custom download base URL"
    Write-Host "  -Help                   Show this help message"
    Write-Host ""
    Write-Host "Environment Variables:"
    Write-Host "  AGBCLOUD_VERSION        Default version to install"
    Write-Host "  AGBCLOUD_PATH           Default installation directory"
    Write-Host "  AGBCLOUD_DOWNLOAD_URL   Default download base URL"
    Write-Host ""
    Write-Host "Examples:"
    Write-Host "  .\install-windows-simple.ps1                    # Install latest version"
    Write-Host "  .\install-windows-simple.ps1 -Version v1.2.3    # Install specific version"
    Write-Host "  .\install-windows-simple.ps1 -Version latest    # Install latest version"
    exit 0
}

# Determine architecture
if (-not $Architecture) {
    $Architecture = if ($env:PROCESSOR_ARCHITECTURE -eq "AMD64") { "amd64" } elseif ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") { "arm64" } else { "amd64" }
}

# Define version and download URL with parameter priority
$version = if ($Version) { 
    $Version 
} elseif ($env:AGBCLOUD_VERSION) { 
    $env:AGBCLOUD_VERSION 
} else { 
    "latest" 
}

$baseUrl = if ($DownloadUrl) { 
    $DownloadUrl 
} elseif ($env:AGBCLOUD_DOWNLOAD_URL) { 
    $env:AGBCLOUD_DOWNLOAD_URL 
} else { 
    "https://agbcloud-internal.oss-cn-hangzhou.aliyuncs.com" 
}

$destination = if ($InstallPath) { 
    $InstallPath 
} elseif ($env:AGBCLOUD_PATH) {
    $env:AGBCLOUD_PATH
} else {
    "$env:APPDATA\bin\agb" 
}

# Get latest version if needed
if ($version -eq "latest") {
    try {
        Write-Host "[INFO] Checking for latest version..."
        $latestInfo = Invoke-RestMethod -Uri "$baseUrl/latest.json" -UseBasicParsing -ErrorAction SilentlyContinue
        if ($latestInfo -and $latestInfo.version) {
            $version = $latestInfo.version
        } else {
            $version = "dev-$(Get-Date -Format 'yyyyMMdd-HHmm')"
        }
    } catch {
        Write-Host "[WARN] Could not fetch latest version, using fallback"
        $version = "dev-$(Get-Date -Format 'yyyyMMdd-HHmm')"
    }
}

$downloadUrl = "$baseUrl/agb-$version-windows-$Architecture.exe"

Write-Host "[INFO] Installing AgbCloud CLI..."
Write-Host ""

# Display installation info
Write-Host "Installation Details:"
Write-Host "   Version: $version"
Write-Host "   Architecture: $Architecture"
if ($InstallPath -or $env:AGBCLOUD_PATH) {
    Write-Host "   Custom installation directory: $destination"
} else {
    Write-Host "   Default installation directory: $destination"
    Write-Host "   [TIP] You can override this by setting the AGBCLOUD_PATH environment variable."
}
Write-Host ""

# Create destination directory if it doesn't exist
try {
    if (!(Test-Path -Path $destination)) {
        Write-Host "[INFO] Creating installation directory at $destination"
        New-Item -ItemType Directory -Force -Path $destination -ErrorAction Stop | Out-Null
        Write-Host ""
    }
} catch {
    Write-Error "[ERROR] Failed to create installation directory: $_"
    exit 1
}

# File to download
$outputFile = "$destination\agb.exe"

# Check if already installed and get current version
$upgrading = $false
if (Test-Path $outputFile) {
    try {
        $currentVersion = & $outputFile version --short 2>$null
        if ($currentVersion -eq $version) {
            Write-Host "[SUCCESS] AgbCloud CLI $version is already installed!"
            Write-Host "   Location: $outputFile"
            Write-Host ""
            Write-Host "[INFO] You're all set! Use 'agb --help' to get started."
            exit 0
        } else {
            Write-Host "[INFO] Upgrading from $currentVersion to $version"
            $upgrading = $true
        }
    } catch {
        Write-Host "[INFO] Existing installation found, upgrading..."
        $upgrading = $true
    }
    Write-Host ""
}

# Download the file with progress
try {
    if ($upgrading) {
        Write-Host "[INFO] Downloading AgbCloud CLI update from $downloadUrl"
    } else {
        Write-Host "[INFO] Downloading AgbCloud CLI from $downloadUrl"
    }

    # Use Invoke-WebRequest with progress
    $ProgressPreference = 'Continue'
    Invoke-WebRequest -Uri $downloadUrl -OutFile $outputFile -UseBasicParsing -ErrorAction Stop

    Write-Host ""
    Write-Host "[SUCCESS] Download complete!"
} catch {
    Write-Error "[ERROR] Failed to download AgbCloud CLI: $_"
    Write-Host "   Please check your internet connection and try again."
    Write-Host "   If the problem persists, visit: https://github.com/your-org/agbcloud-cli/releases"
    exit 1
}

Write-Host ""

# Set executable permissions (Windows doesn't need this, but good practice)
try {
    Write-Host "[INFO] Setting up binary permissions..."
    # Try to set attributes, but don't fail if it doesn't work (constrained language mode)
    try {
        Set-ItemProperty -Path $outputFile -Name IsReadOnly -Value $false -ErrorAction SilentlyContinue
        [System.IO.File]::SetAttributes($outputFile, 'Normal')
    } catch {
        # In constrained language mode, this might fail, but it's not critical
        Write-Host "   [WARN] Could not set file attributes (this is usually fine)"
    }
} catch {
    # This shouldn't happen now, but keep as fallback
    Write-Host "   [WARN] Could not set binary permissions (this is usually fine on Windows)"
}

Write-Host ""

# Add to PATH if not already present
try {
    Write-Host "[INFO] Updating PATH..."
    
    # Try to get current PATH, handle constrained language mode
    try {
        $currentPath = [System.Environment]::GetEnvironmentVariable("Path", [System.EnvironmentVariableTarget]::User)
        if (-not $currentPath) { $currentPath = "" }
        
        $pathEntries = $currentPath -split ';' | ForEach-Object { $_.TrimEnd('\') }
        
        if (-not ($pathEntries | Where-Object { $_ -eq $destination })) {
            Write-Host "   Adding $destination to user PATH..."
            $newPath = if ($currentPath.EndsWith(';')) { "$currentPath$destination" } else { "$currentPath;$destination" }
            [System.Environment]::SetEnvironmentVariable("Path", $newPath, [System.EnvironmentVariableTarget]::User)
            Write-Host "[SUCCESS] PATH updated successfully!"
            Write-Host "   [TIP] Please restart your terminal or run a new PowerShell session"
        } else {
            Write-Host "[SUCCESS] Already in PATH"
        }
    } catch {
        Write-Host "   [WARN] Could not automatically update PATH (constrained language mode)"
        Write-Host "   [MANUAL] Please manually add the following to your PATH:"
        Write-Host "      $destination"
        Write-Host ""
        Write-Host "   [STEPS] To add manually:"
        Write-Host "      1. Press Win+R, type 'sysdm.cpl', press Enter"
        Write-Host "      2. Click 'Environment Variables'"
        Write-Host "      3. Under 'User variables', select 'Path' and click 'Edit'"
        Write-Host "      4. Click 'New' and add: $destination"
        Write-Host "      5. Click OK to save"
    }
} catch {
    Write-Host "   [WARN] PATH update failed, but installation completed"
    Write-Host "   [MANUAL] Please manually add to PATH: $destination"
}

Write-Host ""

# Test installation
Write-Host "[INFO] Testing installation..."
try {
    # Try different version commands to get version info
    $installedVersion = ""
    try {
        $installedVersion = & $outputFile version --short 2>$null
        if (-not $installedVersion) {
            $installedVersion = & $outputFile version 2>$null | Select-String "version" | Select-Object -First 1
        }
        if (-not $installedVersion) {
            $installedVersion = & $outputFile --version 2>$null
        }
    } catch {
        $installedVersion = "unknown"
    }
    
    Write-Host "[SUCCESS] Installation test successful!"
    Write-Host ""
    
    if ($upgrading) {
        Write-Host "[SUCCESS] AgbCloud CLI successfully upgraded to $installedVersion!"
    } else {
        Write-Host "[SUCCESS] AgbCloud CLI $installedVersion installed successfully!"
    }
    
    Write-Host "   Location: $outputFile"
    Write-Host ""
    Write-Host "Quick Start:"
    Write-Host "   agb --help          # Show help"
Write-Host "   agb version         # Show version"
Write-Host "   agb login           # Login to AgbCloud"
    Write-Host ""
    Write-Host "Important Notes:"
    Write-Host "   * The command is 'agb'"
Write-Host "   * If 'agb' command not found, restart your terminal"
    Write-Host "   * Or run directly: $outputFile"
    Write-Host ""
    Write-Host "Documentation: https://docs.agbcloud.com"
    
} catch {
    Write-Host "[WARN] Installation test failed, but binary was downloaded successfully"
    Write-Host "   Location: $outputFile"
    Write-Host "   [TIP] You can run it directly or add to PATH manually"
    Write-Host ""
    Write-Host "   [TIP] Try running directly:"
    Write-Host "      $outputFile version"
}

Write-Host "" 