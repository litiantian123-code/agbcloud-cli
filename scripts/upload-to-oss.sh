#!/bin/bash
# Upload packages to Alibaba Cloud OSS
set -e

# OSS Configuration
OSS_ENDPOINT="oss-cn-hangzhou.aliyuncs.com"
OSS_BUCKET="agbcloud-internal"
OSS_PREFIX="agbcloud/releases"

# Package directory
PACKAGE_DIR=${PACKAGE_DIR:-"packages"}
VERSION=${VERSION:-"dev-$(date +%Y%m%d-%H%M)"}

echo "Uploading packages to OSS..."
echo "Endpoint: $OSS_ENDPOINT"
echo "Bucket: $OSS_BUCKET"
echo "Prefix: $OSS_PREFIX"
echo "Version: $VERSION"

# Check required environment variables
if [[ -z "$OSS_ACCESS_KEY_ID" ]]; then
    echo "Error: OSS_ACCESS_KEY_ID environment variable is not set"
    exit 1
fi

if [[ -z "$OSS_ACCESS_KEY_SECRET" ]]; then
    echo "Error: OSS_ACCESS_KEY_SECRET environment variable is not set"
    exit 1
fi

# Check if package directory exists
if [[ ! -d "$PACKAGE_DIR" ]]; then
    echo "Error: Package directory '$PACKAGE_DIR' not found"
    exit 1
fi

# Check if ossutil is installed
if ! command -v ossutil >/dev/null 2>&1; then
    echo "Installing ossutil..."
    
    # Download and install ossutil
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        if [[ "$(uname -m)" == "arm64" ]]; then
            OSSUTIL_URL="https://gosspublic.alicdn.com/ossutil/1.7.19/ossutil-v1.7.19-darwin-arm64.zip"
        else
            OSSUTIL_URL="https://gosspublic.alicdn.com/ossutil/1.7.19/ossutil-v1.7.19-darwin-amd64.zip"
        fi
    else
        # Linux
        OSSUTIL_URL="https://gosspublic.alicdn.com/ossutil/1.7.19/ossutil-v1.7.19-linux-amd64.zip"
    fi
    
    curl -L "$OSSUTIL_URL" -o ossutil.zip
    unzip -q ossutil.zip
    chmod +x ossutil*
    sudo mv ossutil* /usr/local/bin/ossutil
    rm -f ossutil.zip
    
    echo "✓ ossutil installed successfully"
fi

# Configure ossutil
echo "Configuring ossutil..."
ossutil config -e "$OSS_ENDPOINT" -i "$OSS_ACCESS_KEY_ID" -k "$OSS_ACCESS_KEY_SECRET"

# Upload all package files
echo "Uploading packages..."
upload_count=0

for package in "$PACKAGE_DIR"/*.tar.gz; do
    if [[ -f "$package" ]]; then
        filename=$(basename "$package")
        oss_path="oss://$OSS_BUCKET/$OSS_PREFIX/$filename"
        
        echo "Uploading $filename..."
        
        # Upload the package
        if ossutil cp "$package" "$oss_path" --force; then
            echo "✓ Uploaded: $filename"
            
            # Set public read permission
            if ossutil set-acl "$oss_path" public-read; then
                echo "✓ Set public-read ACL for: $filename"
            else
                echo "⚠ Failed to set ACL for: $filename"
            fi
            
            # Generate public URL
            public_url="https://$OSS_BUCKET.$OSS_ENDPOINT/$OSS_PREFIX/$filename"
            echo "  Public URL: $public_url"
            
            upload_count=$((upload_count + 1))
        else
            echo "✗ Failed to upload: $filename"
            exit 1
        fi
    fi
done

# Upload SHA256 files
echo "Uploading SHA256 files..."
for sha256_file in "$PACKAGE_DIR"/*.sha256; do
    if [[ -f "$sha256_file" ]]; then
        filename=$(basename "$sha256_file")
        oss_path="oss://$OSS_BUCKET/$OSS_PREFIX/$filename"
        
        echo "Uploading $filename..."
        
        if ossutil cp "$sha256_file" "$oss_path" --force; then
            echo "✓ Uploaded: $filename"
            
            # Set public read permission
            ossutil set-acl "$oss_path" public-read
        else
            echo "⚠ Failed to upload: $filename"
        fi
    fi
done

echo ""
echo "Upload completed successfully!"
echo "Total packages uploaded: $upload_count"
echo ""
echo "Download URLs:"
for package in "$PACKAGE_DIR"/*.tar.gz; do
    if [[ -f "$package" ]]; then
        filename=$(basename "$package")
        echo "  https://$OSS_BUCKET.$OSS_ENDPOINT/$OSS_PREFIX/$filename"
    fi
done 