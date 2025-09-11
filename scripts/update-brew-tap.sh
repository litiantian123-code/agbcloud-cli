#!/bin/bash
# Update Homebrew Tap repository with new Formula
set -e

# Configuration
VERSION=${VERSION:-"dev-$(date +%Y%m%d-%H%M)"}
GIT_COMMIT=${GIT_COMMIT:-"$(git rev-parse --short HEAD)"}
TIMESTAMP=${TIMESTAMP:-"$(date +%Y%m%d-%H%M)"}
BREW_TAP_REPO=${BREW_TAP_REPO:-"git@gitlab.alibaba-inc.com:InnoArchClub/homebrew-agbcloud.git"}
BREW_TAP_DIR="homebrew-agbcloud"

echo "Updating Homebrew Tap..."
echo "Version: $VERSION"
echo "Timestamp: $TIMESTAMP"
echo "Git commit: $GIT_COMMIT"
echo "Tap repository: $BREW_TAP_REPO"

# Clone or update brew tap repository
if [[ ! -d "$BREW_TAP_DIR" ]]; then
    echo "Cloning brew tap repository..."
    git clone "$BREW_TAP_REPO" "$BREW_TAP_DIR"
else
    echo "Updating existing brew tap repository..."
    cd "$BREW_TAP_DIR"
    git fetch origin
    git reset --hard origin/main
    cd ..
fi

# Copy current scripts to tap repository
echo "Copying scripts to tap repository..."
cp -r scripts "$BREW_TAP_DIR/"

# Generate new Formula
echo "Generating Formula..."
cd "$BREW_TAP_DIR"

# Set package directory relative to tap repo
export PACKAGE_DIR="../packages"

# Generate the Formula
if ./scripts/generate-formula.sh "$VERSION" "$GIT_COMMIT" "$TIMESTAMP"; then
    echo "✓ Formula generated successfully"
else
    echo "✗ Failed to generate Formula"
    exit 1
fi

# Check if Formula was created
FORMULA_FILE="Formula/agbcloud@dev-$TIMESTAMP.rb"
if [[ ! -f "$FORMULA_FILE" ]]; then
    echo "Error: Formula file not found: $FORMULA_FILE"
    exit 1
fi

# Configure git if needed
if ! git config user.name >/dev/null 2>&1; then
    git config user.name "CICD Bot"
fi

if ! git config user.email >/dev/null 2>&1; then
    git config user.email "cicd@your-company.com"
fi

# Add and commit the new Formula
echo "Committing new Formula..."
git add "$FORMULA_FILE"
git add scripts/

# Create commit message
COMMIT_MSG="Add test build agbcloud@dev-$TIMESTAMP

Version: $VERSION
Git commit: $GIT_COMMIT
Build timestamp: $TIMESTAMP

Generated Formula: $FORMULA_FILE"

if git commit -m "$COMMIT_MSG"; then
    echo "✓ Formula committed successfully"
else
    echo "⚠ No changes to commit (Formula may already exist)"
fi

# Push to remote repository
echo "Pushing to remote repository..."
if git push origin main; then
    echo "✓ Formula pushed to remote repository"
else
    echo "✗ Failed to push to remote repository"
    exit 1
fi

cd ..

echo ""
echo "Homebrew Tap update completed successfully!"
echo ""
echo "Users can now install with:"
echo "  brew install agbcloud@dev-$TIMESTAMP"
echo ""
echo "Formula location:"
echo "  $FORMULA_FILE" 