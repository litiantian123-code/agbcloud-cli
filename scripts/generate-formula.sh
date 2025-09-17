#!/bin/bash
# Generate Homebrew Formula for test builds
set -e

# Input parameters
VERSION=${1:-"dev-$(date +%Y%m%d-%H%M)"}
GIT_COMMIT=${2:-"$(git rev-parse --short HEAD)"}
TIMESTAMP=${3:-"$(date +%Y%m%d-%H%M)"}

# Sanitize timestamp for class name (remove hyphens)
SANITIZED_TIMESTAMP=$(echo "$TIMESTAMP" | tr '-' '')

# Package directory (should contain .tar.gz files and .sha256 files)
PACKAGE_DIR=${PACKAGE_DIR:-"packages"}

echo "Generating Formula for version: $VERSION"
echo "Timestamp: $TIMESTAMP"
echo "Git commit: $GIT_COMMIT"

# Check if package directory exists
if [[ ! -d "$PACKAGE_DIR" ]]; then
    echo "Error: Package directory '$PACKAGE_DIR' not found"
    exit 1
fi

# Function to get SHA256 from .sha256 file
get_sha256() {
    local platform_arch=$1
    local sha256_file="$PACKAGE_DIR/agbcloud-$VERSION-$platform_arch.tar.gz.sha256"
    
    if [[ -f "$sha256_file" ]]; then
        cut -d' ' -f1 "$sha256_file"
    else
        echo "Error: SHA256 file not found: $sha256_file" >&2
        echo "MISSING_SHA256"
    fi
}

# Get SHA256 values for all platforms
DARWIN_AMD64_SHA256=$(get_sha256 "darwin-amd64")
DARWIN_ARM64_SHA256=$(get_sha256 "darwin-arm64")
LINUX_AMD64_SHA256=$(get_sha256 "linux-amd64")
LINUX_ARM64_SHA256=$(get_sha256 "linux-arm64")

# Check if all SHA256 values are available
if [[ "$DARWIN_AMD64_SHA256" == "MISSING_SHA256" ]] || 
   [[ "$DARWIN_ARM64_SHA256" == "MISSING_SHA256" ]] || 
   [[ "$LINUX_AMD64_SHA256" == "MISSING_SHA256" ]] || 
   [[ "$LINUX_ARM64_SHA256" == "MISSING_SHA256" ]]; then
    echo "Error: Some SHA256 files are missing"
    exit 1
fi

# Generate Formula file
FORMULA_FILE="Formula/agbcloud@dev-$TIMESTAMP.rb"

echo "Generating $FORMULA_FILE..."

# Use sed to replace template variables
sed -e "s/<%= sanitized_timestamp %>/$SANITIZED_TIMESTAMP/g" \
    -e "s/<%= timestamp %>/$TIMESTAMP/g" \
    -e "s/<%= version %>/$VERSION/g" \
    -e "s/<%= git_commit %>/$GIT_COMMIT/g" \
    -e "s/<%= darwin_amd64_sha256 %>/$DARWIN_AMD64_SHA256/g" \
    -e "s/<%= darwin_arm64_sha256 %>/$DARWIN_ARM64_SHA256/g" \
    -e "s/<%= linux_amd64_sha256 %>/$LINUX_AMD64_SHA256/g" \
    -e "s/<%= linux_arm64_sha256 %>/$LINUX_ARM64_SHA256/g" \
    scripts/formula-template.rb.tpl > "$FORMULA_FILE"

echo "[OK] Formula generated: $FORMULA_FILE"

# Validate the generated Formula
if command -v ruby >/dev/null 2>&1; then
    echo "Validating Formula syntax..."
    if ruby -c "$FORMULA_FILE" >/dev/null 2>&1; then
        echo "[OK] Formula syntax is valid"
    else
        echo "[WARN] Formula syntax validation failed"
        ruby -c "$FORMULA_FILE"
    fi
fi

echo "Formula generation completed successfully!" 