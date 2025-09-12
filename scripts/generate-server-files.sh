#!/bin/bash

# Generate server files for AgbCloud CLI distribution
# This script creates only the essential files needed for PowerShell installation

set -e

VERSION=${VERSION:-"dev-$(date +%Y%m%d-%H%M)"}
BASE_URL="https://agbcloud-internal.oss-cn-hangzhou.aliyuncs.com"

echo "🚀 Generating PowerShell installer files for version: $VERSION"

# Create output directory
mkdir -p server-files

# 1. Generate latest.json for version API (essential for PowerShell script)
echo "📄 Creating latest.json..."
cat > server-files/latest.json << EOF
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
cp scripts/install-windows-simple.ps1 server-files/install.ps1

echo "✅ PowerShell installer files generated successfully!"
echo ""
echo "📁 Generated files:"
ls -la server-files/
echo ""
echo "🚀 Upload these files to your OSS bucket:"
echo "   - server-files/install.ps1 → $BASE_URL/install.ps1"
echo "   - server-files/latest.json → $BASE_URL/latest.json"
echo "   - packages/*.exe → $BASE_URL/"
echo "   - packages/*.exe.sha256 → $BASE_URL/"
echo ""
echo "📋 Windows 用户安装命令:"
echo "   powershell -Command \"irm $BASE_URL/install.ps1 | iex\"" 