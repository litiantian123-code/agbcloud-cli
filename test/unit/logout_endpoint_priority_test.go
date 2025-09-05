// Copyright 2025 AgbCloud CLI Contributors
// SPDX-License-Identifier: Apache-2.0

package unit

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/agbcloud/agbcloud-cli/internal/client"
	"github.com/agbcloud/agbcloud-cli/internal/config"
)

func TestLogoutUsesEnvironmentEndpoint(t *testing.T) {
	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request path and method
		if r.URL.Path != "/api/biz_login/logout" {
			t.Errorf("Expected path /api/biz_login/logout, got %s", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		// Verify query parameters
		loginToken := r.URL.Query().Get("loginToken")
		sessionId := r.URL.Query().Get("sessionId")
		if loginToken != "test-login-token" {
			t.Errorf("Expected loginToken 'test-login-token', got %s", loginToken)
		}
		if sessionId != "test-session-id" {
			t.Errorf("Expected sessionId 'test-session-id', got %s", sessionId)
		}

		// Return success response
		response := map[string]interface{}{
			"success":        true,
			"code":           "200",
			"requestId":      "test-request-id",
			"traceId":        "test-trace-id",
			"httpStatusCode": 200,
			"data": map[string]interface{}{
				"message": "Logout successful",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	// Create temporary config directory
	tempDir, err := os.MkdirTemp("", "agbcloud-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set config directory environment variable
	originalConfigDir := os.Getenv("AGB_CLI_CONFIG_DIR")
	os.Setenv("AGB_CLI_CONFIG_DIR", tempDir)
	defer func() {
		if originalConfigDir == "" {
			os.Unsetenv("AGB_CLI_CONFIG_DIR")
		} else {
			os.Setenv("AGB_CLI_CONFIG_DIR", originalConfigDir)
		}
	}()

	// Set endpoint environment variable to mock server URL
	originalEndpoint := os.Getenv("AGB_CLI_ENDPOINT")
	os.Setenv("AGB_CLI_ENDPOINT", mockServer.URL)
	defer func() {
		if originalEndpoint == "" {
			os.Unsetenv("AGB_CLI_ENDPOINT")
		} else {
			os.Setenv("AGB_CLI_ENDPOINT", originalEndpoint)
		}
	}()

	// Create config file with different endpoint (should be overridden by env var)
	configPath := filepath.Join(tempDir, "config.json")
	configData := map[string]interface{}{
		"endpoint": "https://config.example.com", // This should be overridden
		"token": map[string]interface{}{
			"loginToken":     "test-login-token",
			"sessionId":      "test-session-id",
			"keepAliveToken": "test-keep-alive-token",
			"expiresAt":      time.Now().Add(time.Hour).Format(time.RFC3339),
		},
	}

	configBytes, err := json.MarshalIndent(configData, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}

	err = os.WriteFile(configPath, configBytes, 0600)
	if err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Load config
	cfg, err := config.GetConfig()
	if err != nil {
		t.Fatalf("Failed to get config: %v", err)
	}

	// Verify that config uses environment endpoint
	if cfg.Endpoint != mockServer.URL {
		t.Errorf("Expected endpoint %q, got %q", mockServer.URL, cfg.Endpoint)
	}

	// Create API client and test logout
	apiClient := client.NewFromConfig(cfg)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Call logout API
	response, httpResp, err := apiClient.OAuthAPI.Logout(ctx, cfg.Token.LoginToken, cfg.Token.SessionId)
	if err != nil {
		t.Fatalf("Logout API call failed: %v", err)
	}

	if httpResp.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", httpResp.StatusCode)
	}

	if !response.Success {
		t.Errorf("Expected success=true, got %v", response.Success)
	}

	if response.Code != "200" {
		t.Errorf("Expected code='200', got %s", response.Code)
	}

	t.Logf("✅ Logout successfully used environment endpoint: %s", mockServer.URL)
	t.Logf("✅ Config file endpoint was correctly overridden by environment variable")
}

func TestLogoutUsesConfigFileEndpoint(t *testing.T) {
	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return success response
		response := map[string]interface{}{
			"success":        true,
			"code":           "200",
			"requestId":      "test-request-id",
			"traceId":        "test-trace-id",
			"httpStatusCode": 200,
			"data": map[string]interface{}{
				"message": "Logout successful",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	// Create temporary config directory
	tempDir, err := os.MkdirTemp("", "agbcloud-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set config directory environment variable
	originalConfigDir := os.Getenv("AGB_CLI_CONFIG_DIR")
	os.Setenv("AGB_CLI_CONFIG_DIR", tempDir)
	defer func() {
		if originalConfigDir == "" {
			os.Unsetenv("AGB_CLI_CONFIG_DIR")
		} else {
			os.Setenv("AGB_CLI_CONFIG_DIR", originalConfigDir)
		}
	}()

	// Ensure no endpoint environment variable is set
	originalEndpoint := os.Getenv("AGB_CLI_ENDPOINT")
	os.Unsetenv("AGB_CLI_ENDPOINT")
	defer func() {
		if originalEndpoint == "" {
			os.Unsetenv("AGB_CLI_ENDPOINT")
		} else {
			os.Setenv("AGB_CLI_ENDPOINT", originalEndpoint)
		}
	}()

	// Create config file with mock server endpoint
	configPath := filepath.Join(tempDir, "config.json")
	configData := map[string]interface{}{
		"endpoint": mockServer.URL,
		"token": map[string]interface{}{
			"loginToken":     "test-login-token",
			"sessionId":      "test-session-id",
			"keepAliveToken": "test-keep-alive-token",
			"expiresAt":      time.Now().Add(time.Hour).Format(time.RFC3339),
		},
	}

	configBytes, err := json.MarshalIndent(configData, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}

	err = os.WriteFile(configPath, configBytes, 0600)
	if err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Load config
	cfg, err := config.GetConfig()
	if err != nil {
		t.Fatalf("Failed to get config: %v", err)
	}

	// Verify that config uses config file endpoint
	if cfg.Endpoint != mockServer.URL {
		t.Errorf("Expected endpoint %q, got %q", mockServer.URL, cfg.Endpoint)
	}

	// Create API client and test logout
	apiClient := client.NewFromConfig(cfg)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Call logout API
	response, httpResp, err := apiClient.OAuthAPI.Logout(ctx, cfg.Token.LoginToken, cfg.Token.SessionId)
	if err != nil {
		t.Fatalf("Logout API call failed: %v", err)
	}

	if httpResp.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", httpResp.StatusCode)
	}

	if !response.Success {
		t.Errorf("Expected success=true, got %v", response.Success)
	}

	t.Logf("✅ Logout successfully used config file endpoint: %s", mockServer.URL)
}

func TestLogoutUsesDefaultEndpoint(t *testing.T) {
	// Create temporary config directory
	tempDir, err := os.MkdirTemp("", "agbcloud-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set config directory environment variable
	originalConfigDir := os.Getenv("AGB_CLI_CONFIG_DIR")
	os.Setenv("AGB_CLI_CONFIG_DIR", tempDir)
	defer func() {
		if originalConfigDir == "" {
			os.Unsetenv("AGB_CLI_CONFIG_DIR")
		} else {
			os.Setenv("AGB_CLI_CONFIG_DIR", originalConfigDir)
		}
	}()

	// Ensure no endpoint environment variable is set
	originalEndpoint := os.Getenv("AGB_CLI_ENDPOINT")
	os.Unsetenv("AGB_CLI_ENDPOINT")
	defer func() {
		if originalEndpoint == "" {
			os.Unsetenv("AGB_CLI_ENDPOINT")
		} else {
			os.Setenv("AGB_CLI_ENDPOINT", originalEndpoint)
		}
	}()

	// Create config file without endpoint
	configPath := filepath.Join(tempDir, "config.json")
	configData := map[string]interface{}{
		"token": map[string]interface{}{
			"loginToken":     "test-login-token",
			"sessionId":      "test-session-id",
			"keepAliveToken": "test-keep-alive-token",
			"expiresAt":      time.Now().Add(time.Hour).Format(time.RFC3339),
		},
	}

	configBytes, err := json.MarshalIndent(configData, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}

	err = os.WriteFile(configPath, configBytes, 0600)
	if err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Load config
	cfg, err := config.GetConfig()
	if err != nil {
		t.Fatalf("Failed to get config: %v", err)
	}

	// Verify that config uses default endpoint
	expectedEndpoint := "https://agb.cloud"
	if cfg.Endpoint != expectedEndpoint {
		t.Errorf("Expected endpoint %q, got %q", expectedEndpoint, cfg.Endpoint)
	}

	// Create API client
	apiClient := client.NewFromConfig(cfg)

	// Verify that the client configuration uses the default endpoint
	clientConfig := apiClient.GetConfig()
	if len(clientConfig.Servers) == 0 {
		t.Fatal("No servers configured in client")
	}

	serverURL := clientConfig.Servers[0].URL
	if !strings.Contains(serverURL, "agb.cloud") {
		t.Errorf("Expected server URL to contain 'agb.cloud', got %q", serverURL)
	}

	t.Logf("✅ Config correctly uses default endpoint: %s", cfg.Endpoint)
	t.Logf("✅ Client correctly configured with server URL: %s", serverURL)
}
