#!/bin/bash

# Copyright 2025 AgbCloud CLI Contributors
# SPDX-License-Identifier: Apache-2.0

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Default values
RUN_UNIT=true
RUN_INTEGRATION=false
VERBOSE=false

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --unit-only)
            RUN_UNIT=true
            RUN_INTEGRATION=false
            shift
            ;;
        --integration-only)
            RUN_UNIT=false
            RUN_INTEGRATION=true
            shift
            ;;
        --all)
            RUN_UNIT=true
            RUN_INTEGRATION=true
            shift
            ;;
        --verbose|-v)
            VERBOSE=true
            shift
            ;;
        --help|-h)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  --unit-only        Run only unit tests (default)"
            echo "  --integration-only Run only integration tests"
            echo "  --all              Run both unit and integration tests"
            echo "  --verbose, -v      Enable verbose output"
            echo "  --help, -h         Show this help message"
            echo ""
            echo "Environment Variables:"
            echo "  SKIP_INTEGRATION_TESTS Set to 'true' to skip integration tests"
            exit 0
            ;;
        *)
            print_error "Unknown option: $1"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

# Set verbose flag for go test
VERBOSE_FLAG=""
if [ "$VERBOSE" = true ]; then
    VERBOSE_FLAG="-v"
fi

# Change to project root directory
cd "$(dirname "$0")/.."

print_status "Starting test execution..."

# Run unit tests
if [ "$RUN_UNIT" = true ]; then
    print_status "Running unit tests..."
    
    if go test $VERBOSE_FLAG ./test/unit/...; then
        print_status "✅ Unit tests passed"
    else
        print_error "❌ Unit tests failed"
        exit 1
    fi
    
    # Also run tests for internal packages (excluding integration tests)
    print_status "Running internal package tests..."
    if go test $VERBOSE_FLAG ./internal/...; then
        print_status "✅ Internal package tests passed"
    else
        print_error "❌ Internal package tests failed"
        exit 1
    fi
fi

# Run integration tests
if [ "$RUN_INTEGRATION" = true ]; then
    print_status "Running integration tests..."
    
    # Check if integration tests should be skipped
    if [ "$SKIP_INTEGRATION_TESTS" = "true" ]; then
        print_warning "⚠️  Integration tests skipped (SKIP_INTEGRATION_TESTS=true)"
    else
        print_status "Running integration tests against real API..."
        print_warning "Note: Integration tests WILL FAIL if network connectivity to https://agb.cloud is unavailable"
        
        if go test $VERBOSE_FLAG -tags=integration ./test/integration/...; then
            print_status "✅ Integration tests passed"
        else
            print_error "❌ Integration tests failed"
            print_warning "Common causes:"
            print_warning "  - Network connectivity issues to https://agb.cloud"
            print_warning "  - API endpoint changes or service unavailability"
            print_warning "  - Firewall or proxy blocking HTTPS requests"
            print_warning "If you want to skip integration tests, set SKIP_INTEGRATION_TESTS=true"
            # Don't exit with error for integration tests as they may fail due to external dependencies
        fi
    fi
fi

print_status "Test execution completed!"

# Run coverage if requested
if [ "$VERBOSE" = true ]; then
    print_status "Generating test coverage report..."
    go test -coverprofile=coverage.out ./test/unit/... ./internal/...
    go tool cover -html=coverage.out -o coverage.html
    print_status "Coverage report generated: coverage.html"
fi 