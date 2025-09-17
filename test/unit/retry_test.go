package unit

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/agbcloud/agbcloud-cli/internal/client"
)

func TestRetryMechanism(t *testing.T) {
	tests := []struct {
		name             string
		serverBehavior   func(attemptCount *int) http.HandlerFunc
		expectedAttempts int
		shouldSucceed    bool
		description      string
	}{
		{
			name: "Success on first attempt",
			serverBehavior: func(attemptCount *int) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					*attemptCount++
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("success"))
				}
			},
			expectedAttempts: 1,
			shouldSucceed:    true,
			description:      "Request should succeed immediately",
		},
		{
			name: "Success on second attempt",
			serverBehavior: func(attemptCount *int) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					*attemptCount++
					if *attemptCount == 1 {
						// Simulate server error on first attempt (retryable)
						w.WriteHeader(http.StatusInternalServerError)
						w.Write([]byte("server error"))
						return
					}
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("success"))
				}
			},
			expectedAttempts: 2,
			shouldSucceed:    true,
			description:      "Request should succeed on retry",
		},
		{
			name: "Retry on 500 error",
			serverBehavior: func(attemptCount *int) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					*attemptCount++
					if *attemptCount <= 2 {
						w.WriteHeader(http.StatusInternalServerError)
						w.Write([]byte("server error"))
						return
					}
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("success"))
				}
			},
			expectedAttempts: 3,
			shouldSucceed:    true,
			description:      "Should retry on 500 error and eventually succeed",
		},
		{
			name: "No retry on 404 error",
			serverBehavior: func(attemptCount *int) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					*attemptCount++
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte("not found"))
				}
			},
			expectedAttempts: 1,
			shouldSucceed:    false,
			description:      "Should not retry on 404 error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attemptCount := 0

			// Create test server with recovery middleware
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				defer func() {
					if r := recover(); r != nil {
						// Simulate connection error by closing connection
						hj, ok := w.(http.Hijacker)
						if ok {
							conn, _, err := hj.Hijack()
							if err == nil {
								conn.Close()
							}
						}
					}
				}()
				tt.serverBehavior(&attemptCount)(w, r)
			}))
			defer server.Close()

			// Create retry client with fast retry for testing
			retryConfig := &client.RetryConfig{
				MaxRetries:    3,
				InitialDelay:  10 * time.Millisecond,
				MaxDelay:      50 * time.Millisecond,
				BackoffFactor: 1.5,
			}

			baseClient := &http.Client{Timeout: 5 * time.Second}
			retryClient := client.NewRetryableHTTPClient(baseClient, retryConfig)

			// Make request
			req, err := http.NewRequestWithContext(context.Background(), "GET", server.URL, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			resp, err := retryClient.Do(req)

			// Verify results
			if tt.shouldSucceed {
				if err != nil {
					t.Errorf("%s: Expected success but got error: %v", tt.description, err)
				}
				if resp != nil && resp.StatusCode != http.StatusOK {
					t.Errorf("%s: Expected status 200 but got %d", tt.description, resp.StatusCode)
				}
			} else {
				if err == nil && resp != nil && resp.StatusCode == http.StatusOK {
					t.Errorf("%s: Expected failure but got success", tt.description)
				}
			}

			if attemptCount != tt.expectedAttempts {
				t.Errorf("%s: Expected %d attempts but got %d", tt.description, tt.expectedAttempts, attemptCount)
			}

			if resp != nil {
				resp.Body.Close()
			}

			t.Logf("[OK] %s: Made %d attempts as expected", tt.description, attemptCount)
		})
	}
}

func TestRetryableErrorDetection(t *testing.T) {
	// Test the core retry mechanism with known good cases
	tests := []struct {
		name        string
		errorMsg    string
		shouldRetry bool
	}{
		{
			name:        "Bad file descriptor error",
			errorMsg:    "dial tcp 47.236.181.224:443: connect: bad file descriptor",
			shouldRetry: true,
		},
		{
			name:        "Connection refused error",
			errorMsg:    "dial tcp 127.0.0.1:8080: connect: connection refused",
			shouldRetry: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fmt.Errorf("%s", tt.errorMsg)

			retryConfig := &client.RetryConfig{
				MaxRetries:    1,
				InitialDelay:  1 * time.Millisecond,
				MaxDelay:      1 * time.Millisecond,
				BackoffFactor: 1.0,
			}

			baseClient := &http.Client{
				Transport: &errorTransport{err: err},
				Timeout:   1 * time.Second,
			}

			retryClient := client.NewRetryableHTTPClient(baseClient, retryConfig)

			req, _ := http.NewRequest("GET", "http://example.com", nil)
			_, retryErr := retryClient.Do(req)

			// For retryable errors, we should see "request failed after X attempts"
			isRetryAttempted := strings.Contains(retryErr.Error(), "request failed after")

			if tt.shouldRetry && !isRetryAttempted {
				t.Errorf("Expected error '%s' to be retryable, but no retry was attempted", tt.errorMsg)
			}

			t.Logf("[OK] Error '%s' retry behavior: retryable=%v, attempted=%v",
				tt.errorMsg, tt.shouldRetry, isRetryAttempted)
		})
	}
}

// errorTransport is a test helper that always returns a specific error
type errorTransport struct {
	err error
}

func (et *errorTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return nil, et.err
}
