#!/bin/bash

# SSL Verification Strategy Demo
# This script demonstrates how the CLI automatically handles SSL verification

echo "🔒 SSL Verification Strategy Demo"
echo "================================="
echo

# Build the CLI
echo "📦 Building CLI..."
go build -o agbcloud .
echo

# Test different endpoint types
echo "🌐 Testing different endpoint types:"
echo

echo "1️⃣  Production domain (SSL verification enabled):"
AGB_CLI_ENDPOINT=agb.cloud ./agbcloud config get endpoint
echo "   ✅ Uses SSL verification for security"
echo

echo "2️⃣  IP address (SSL verification disabled):"
AGB_CLI_ENDPOINT=12.34.56.78 ./agbcloud config get endpoint
echo "   ⚠️  Skips SSL verification (certificates don't include IP SANs)"
echo

echo "3️⃣  Localhost (SSL verification disabled):"
AGB_CLI_ENDPOINT=localhost:8080 ./agbcloud config get endpoint
echo "   🛠️  Skips SSL verification for local development"
echo

echo "4️⃣  Development domain (SSL verification disabled):"
AGB_CLI_ENDPOINT=api.local ./agbcloud config get endpoint
echo "   🧪 Skips SSL verification for .local domains"
echo

echo "5️⃣  Manual override - Force SSL verification:"
AGB_CLI_SKIP_SSL_VERIFY=false AGB_CLI_ENDPOINT=12.34.56.78 ./agbcloud config get endpoint
echo "   🔐 Forces SSL verification even for IP addresses"
echo

echo "6️⃣  Manual override - Skip SSL verification:"
AGB_CLI_SKIP_SSL_VERIFY=true AGB_CLI_ENDPOINT=agb.cloud ./agbcloud config get endpoint
echo "   ⚠️  Skips SSL verification even for production domains"
echo

echo "📋 Summary:"
echo "   • Production domains: SSL verification enabled by default"
echo "   • IP addresses, localhost, dev domains: SSL verification disabled by default"
echo "   • Use AGB_CLI_SKIP_SSL_VERIFY to override automatic behavior"
echo "   • Security-first approach: when in doubt, verify SSL"
echo

echo "✅ Demo completed!" 