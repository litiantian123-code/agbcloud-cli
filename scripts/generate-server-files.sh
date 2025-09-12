#!/bin/bash

# Generate server files for AgbCloud CLI distribution
# This script creates only the essential files needed for PowerShell installation

set -e

VERSION=${VERSION:-"dev-$(date +%Y%m%d-%H%M)"}
BASE_URL="https://agbcloud-internal.oss-cn-hangzhou.aliyuncs.com"

echo "🚀 Generating PowerShell installer files for version: $VERSION"

# Create output directory
# Handle both running from root directory and from scripts directory
if [[ -d "scripts" ]]; then
    # Running from root directory
    mkdir -p server-files
    OUTPUT_DIR="server-files"
else
    # Running from scripts directory
    mkdir -p ../server-files
    OUTPUT_DIR="../server-files"
fi

# 1. Generate latest.json for version API (essential for PowerShell script)
echo "📄 Creating latest.json..."
cat > $OUTPUT_DIR/latest.json << EOF
{
  "version": "$VERSION",
  "releaseDate": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
  "windows": {
    "amd64": {
      "url": "$BASE_URL/agbcloud-$VERSION-windows-amd64.exe",
      "sha256": "$(cat packages/agbcloud-$VERSION-windows-amd64.exe.sha256 2>/dev/null | cut -d' ' -f1 || echo 'PLACEHOLDER_SHA256')"
    },
    "arm64": {
      "url": "$BASE_URL/agbcloud-$VERSION-windows-arm64.exe", 
      "sha256": "$(cat packages/agbcloud-$VERSION-windows-arm64.exe.sha256 2>/dev/null | cut -d' ' -f1 || echo 'PLACEHOLDER_SHA256')"
    }
  }
}
EOF

# 2. Generate install.ps1 (the main PowerShell installer)
echo "📄 Creating install.ps1..."
# Copy the PowerShell installer script (it dynamically fetches latest version)
# Handle both running from root directory and from scripts directory
if [[ -f "scripts/install-windows-simple.ps1" ]]; then
    cp scripts/install-windows-simple.ps1 $OUTPUT_DIR/install.ps1
elif [[ -f "install-windows-simple.ps1" ]]; then
    cp install-windows-simple.ps1 $OUTPUT_DIR/install.ps1
else
    echo "Error: Could not find install-windows-simple.ps1"
    exit 1
fi

echo "✅ PowerShell installer files generated successfully!"
echo ""
echo "📁 Generated files:"
ls -la $OUTPUT_DIR/
echo ""
echo "🚀 Upload these files to your OSS bucket:"
echo "   - server-files/install.ps1 → $BASE_URL/install.ps1"
echo "   - server-files/latest.json → $BASE_URL/latest.json"
echo "   - packages/*.exe → $BASE_URL/"
echo "   - packages/*.exe.sha256 → $BASE_URL/"
echo ""
echo "📋 Windows Installation Commands:"
echo ""
echo "🔄 Install Latest Version (recommended for production):"
echo "   powershell -Command \"irm $BASE_URL/install.ps1 | iex\""
echo ""
echo "🎯 Install Specific Version $VERSION (for testing):"
echo "   powershell -Command \"irm $BASE_URL/install.ps1 | iex\" -Version $VERSION"
echo ""
echo "📖 Additional Options:"
echo "   # Install to custom directory"
echo "   powershell -Command \"irm $BASE_URL/install.ps1 | iex\" -InstallPath \"C:\\Tools\\agbcloud\""
echo ""
echo "   # Install specific architecture"
echo "   powershell -Command \"irm $BASE_URL/install.ps1 | iex\" -Architecture arm64"
echo ""
echo "   # Show help"
echo "   powershell -Command \"irm $BASE_URL/install.ps1 | iex\" -Help"
echo ""
echo "💡 Testing Team Usage:"
echo "   Use the specific version command above to test version $VERSION"
echo "   Use the latest version command for general testing" 