// Copyright 2025 AgbCloud CLI Contributors
// SPDX-License-Identifier: Apache-2.0

package unit

import (
	"fmt"
	"net"
	"syscall"
	"testing"

	"github.com/agbcloud/agbcloud-cli/internal/client"
)

func TestIsRetryableError(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		retryable bool
	}{
		{
			name:      "nil error",
			err:       nil,
			retryable: false,
		},
		{
			name:      "bad file descriptor",
			err:       fmt.Errorf("dial tcp 47.79.49.191:443: connect: bad file descriptor"),
			retryable: true,
		},
		{
			name:      "connection refused",
			err:       fmt.Errorf("dial tcp 127.0.0.1:8080: connect: connection refused"),
			retryable: true,
		},
		{
			name:      "timeout error",
			err:       fmt.Errorf("context deadline exceeded"),
			retryable: true,
		},
		{
			name:      "i/o timeout",
			err:       fmt.Errorf("read tcp 192.168.1.1:443: i/o timeout"),
			retryable: true,
		},
		{
			name:      "network unreachable",
			err:       fmt.Errorf("dial tcp 10.0.0.1:443: connect: network is unreachable"),
			retryable: true,
		},
		{
			name:      "broken pipe",
			err:       fmt.Errorf("write tcp 192.168.1.1:443: write: broken pipe"),
			retryable: true,
		},
		{
			name:      "no such host",
			err:       fmt.Errorf("dial tcp: lookup nonexistent.example.com: no such host"),
			retryable: true,
		},
		{
			name:      "invalid uri - not retryable",
			err:       fmt.Errorf("invalid uri for request"),
			retryable: false,
		},
		{
			name:      "parse error - not retryable",
			err:       fmt.Errorf("parse error in URL"),
			retryable: false,
		},
		{
			name:      "malformed request - not retryable",
			err:       fmt.Errorf("malformed HTTP request"),
			retryable: false,
		},
		{
			name:      "unsupported protocol - not retryable",
			err:       fmt.Errorf("unsupported protocol scheme"),
			retryable: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := client.IsRetryableError(tt.err)
			if result != tt.retryable {
				t.Errorf("IsRetryableError(%v) = %v, want %v", tt.err, result, tt.retryable)
			}
		})
	}
}

func TestIsRetryableError_NetworkErrors(t *testing.T) {
	// Test with actual network error types
	tests := []struct {
		name      string
		err       error
		retryable bool
	}{
		{
			name: "net.OpError with ECONNREFUSED",
			err: &net.OpError{
				Op:  "dial",
				Net: "tcp",
				Err: syscall.ECONNREFUSED,
			},
			retryable: true,
		},
		{
			name: "net.OpError with ECONNRESET",
			err: &net.OpError{
				Op:  "read",
				Net: "tcp",
				Err: syscall.ECONNRESET,
			},
			retryable: true,
		},
		{
			name: "net.OpError with ETIMEDOUT",
			err: &net.OpError{
				Op:  "dial",
				Net: "tcp",
				Err: syscall.ETIMEDOUT,
			},
			retryable: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := client.IsRetryableError(tt.err)
			if result != tt.retryable {
				t.Errorf("IsRetryableError(%v) = %v, want %v", tt.err, result, tt.retryable)
			}
		})
	}
}

func TestIsRetryableHTTPStatus(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		retryable  bool
	}{
		{
			name:       "200 OK - not retryable",
			statusCode: 200,
			retryable:  false,
		},
		{
			name:       "400 Bad Request - not retryable",
			statusCode: 400,
			retryable:  false,
		},
		{
			name:       "401 Unauthorized - not retryable",
			statusCode: 401,
			retryable:  false,
		},
		{
			name:       "403 Forbidden - not retryable",
			statusCode: 403,
			retryable:  false,
		},
		{
			name:       "404 Not Found - not retryable",
			statusCode: 404,
			retryable:  false,
		},
		{
			name:       "408 Request Timeout - retryable",
			statusCode: 408,
			retryable:  true,
		},
		{
			name:       "429 Too Many Requests - retryable",
			statusCode: 429,
			retryable:  true,
		},
		{
			name:       "500 Internal Server Error - retryable",
			statusCode: 500,
			retryable:  true,
		},
		{
			name:       "502 Bad Gateway - retryable",
			statusCode: 502,
			retryable:  true,
		},
		{
			name:       "503 Service Unavailable - retryable",
			statusCode: 503,
			retryable:  true,
		},
		{
			name:       "504 Gateway Timeout - retryable",
			statusCode: 504,
			retryable:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := client.IsRetryableHTTPStatus(tt.statusCode)
			if result != tt.retryable {
				t.Errorf("IsRetryableHTTPStatus(%d) = %v, want %v", tt.statusCode, result, tt.retryable)
			}
		})
	}
}

func TestRetryConfig(t *testing.T) {
	config := client.DefaultRetryConfig()

	if config.MaxRetries != 3 {
		t.Errorf("Expected MaxRetries to be 3, got %d", config.MaxRetries)
	}

	if config.InitialDelay.Milliseconds() != 500 {
		t.Errorf("Expected InitialDelay to be 500ms, got %v", config.InitialDelay)
	}

	if config.MaxDelay.Seconds() != 5 {
		t.Errorf("Expected MaxDelay to be 5s, got %v", config.MaxDelay)
	}

	if config.BackoffFactor != 2.0 {
		t.Errorf("Expected BackoffFactor to be 2.0, got %f", config.BackoffFactor)
	}
}
