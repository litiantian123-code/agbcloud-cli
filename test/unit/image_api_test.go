// Copyright 2025 AgbCloud CLI Contributors
// SPDX-License-Identifier: Apache-2.0

package unit

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/agbcloud/agbcloud-cli/internal/client"
)

// TestImageAPIGetUploadCredential tests the GetUploadCredential method with mock server
func TestImageAPIGetUploadCredential(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and path
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/api/image/getUploadCredential" {
			t.Errorf("Expected path /api/image/getUploadCredential, got %s", r.URL.Path)
		}

		// Verify query parameters
		loginToken := r.URL.Query().Get("loginToken")
		sessionId := r.URL.Query().Get("sessionId")
		if loginToken != "test-login-token" {
			t.Errorf("Expected loginToken test-login-token, got %s", loginToken)
		}
		if sessionId != "test-session-id" {
			t.Errorf("Expected sessionId test-session-id, got %s", sessionId)
		}

		// Mock successful response
		response := client.ImageUploadCredentialResponse{
			Code:      "success",
			RequestID: "test-request-id",
			Success:   true,
			Data: client.ImageUploadCredentialData{
				OssURL: "https://test-oss-url.com/upload",
				TaskID: "test-task-id-12345",
			},
			TraceID:        "test-trace-id",
			HTTPStatusCode: 200,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client configuration
	cfg := client.NewConfiguration()
	cfg.Servers[0].URL = server.URL

	// Create API client
	apiClient := client.NewAPIClient(cfg)

	// Test context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Call GetUploadCredential
	response, httpResp, err := apiClient.ImageAPI.GetUploadCredential(ctx, "test-login-token", "test-session-id")

	// Verify no error
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify HTTP response
	if httpResp.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", httpResp.StatusCode)
	}

	// Verify response data
	if !response.Success {
		t.Error("Expected success to be true")
	}
	if response.Code != "success" {
		t.Errorf("Expected code 'success', got '%s'", response.Code)
	}
	if response.Data.TaskID != "test-task-id-12345" {
		t.Errorf("Expected TaskID 'test-task-id-12345', got '%s'", response.Data.TaskID)
	}
	if response.Data.OssURL != "https://test-oss-url.com/upload" {
		t.Errorf("Expected OssURL 'https://test-oss-url.com/upload', got '%s'", response.Data.OssURL)
	}

	t.Logf("✅ GetUploadCredential test passed!")
	t.Logf("   - TaskID: %s", response.Data.TaskID)
	t.Logf("   - OssURL: %s", response.Data.OssURL)
}

// TestImageAPICreateImage tests the CreateImage method with mock server
func TestImageAPICreateImage(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and path
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/api/image/create" {
			t.Errorf("Expected path /api/image/create, got %s", r.URL.Path)
		}

		// Verify Content-Type
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %s", contentType)
		}

		// Parse request body
		var requestBody map[string]string
		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		// Verify request body parameters
		expectedParams := map[string]string{
			"loginToken":    "test-login-token",
			"sessionId":     "test-session-id",
			"imageName":     "test-image-name",
			"taskId":        "test-task-id",
			"sourceImageId": "agb-code-space-2",
		}

		for key, expectedValue := range expectedParams {
			if actualValue, exists := requestBody[key]; !exists {
				t.Errorf("Missing parameter %s in request body", key)
			} else if actualValue != expectedValue {
				t.Errorf("Expected %s '%s', got '%s'", key, expectedValue, actualValue)
			}
		}

		// Mock successful response
		response := client.ImageCreateResponse{
			Code:           "success",
			RequestID:      "test-create-request-id",
			Success:        true,
			Data:           "DBT18253495250135421-17573241195022",
			TraceID:        "test-create-trace-id",
			HTTPStatusCode: 200,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client configuration
	cfg := client.NewConfiguration()
	cfg.Servers[0].URL = server.URL

	// Create API client
	apiClient := client.NewAPIClient(cfg)

	// Test context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Call CreateImage
	response, httpResp, err := apiClient.ImageAPI.CreateImage(ctx,
		"test-login-token",
		"test-session-id",
		"test-image-name",
		"test-task-id",
		"agb-code-space-2")

	// Verify no error
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify HTTP response
	if httpResp.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", httpResp.StatusCode)
	}

	// Verify response data
	if !response.Success {
		t.Error("Expected success to be true")
	}
	if response.Code != "success" {
		t.Errorf("Expected code 'success', got '%s'", response.Code)
	}
	if response.Data != "DBT18253495250135421-17573241195022" {
		t.Errorf("Expected Data 'DBT18253495250135421-17573241195022', got '%s'", response.Data)
	}

	t.Logf("✅ CreateImage test passed!")
	t.Logf("   - Data: %s", response.Data)
}

// TestImageAPIParameterValidation tests parameter validation
func TestImageAPIParameterValidation(t *testing.T) {
	// Create client configuration
	cfg := client.NewConfiguration()
	cfg.Servers[0].URL = "http://localhost:8080" // Dummy URL

	// Create API client
	apiClient := client.NewAPIClient(cfg)

	// Test context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	t.Run("GetUploadCredential_EmptyLoginToken", func(t *testing.T) {
		_, _, err := apiClient.ImageAPI.GetUploadCredential(ctx, "", "test-session-id")
		if err == nil {
			t.Error("Expected error for empty loginToken")
		}
		if !contains(err.Error(), "loginToken parameter is required") {
			t.Errorf("Expected loginToken error message, got: %s", err.Error())
		}
	})

	t.Run("GetUploadCredential_EmptySessionId", func(t *testing.T) {
		_, _, err := apiClient.ImageAPI.GetUploadCredential(ctx, "test-login-token", "")
		if err == nil {
			t.Error("Expected error for empty sessionId")
		}
		if !contains(err.Error(), "sessionId parameter is required") {
			t.Errorf("Expected sessionId error message, got: %s", err.Error())
		}
	})

	t.Run("CreateImage_EmptyParameters", func(t *testing.T) {
		testCases := []struct {
			name          string
			loginToken    string
			sessionId     string
			imageName     string
			taskId        string
			sourceImageId string
			expectedError string
		}{
			{"EmptyLoginToken", "", "sid", "img", "task", "src", "loginToken parameter is required"},
			{"EmptySessionId", "token", "", "img", "task", "src", "sessionId parameter is required"},
			{"EmptyImageName", "token", "sid", "", "task", "src", "imageName parameter is required"},
			{"EmptyTaskId", "token", "sid", "img", "", "src", "taskId parameter is required"},
			{"EmptySourceImageId", "token", "sid", "img", "task", "", "sourceImageId parameter is required"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				_, _, err := apiClient.ImageAPI.CreateImage(ctx, tc.loginToken, tc.sessionId, tc.imageName, tc.taskId, tc.sourceImageId)
				if err == nil {
					t.Errorf("Expected error for %s", tc.name)
				}
				if !contains(err.Error(), tc.expectedError) {
					t.Errorf("Expected error message '%s', got: %s", tc.expectedError, err.Error())
				}
			})
		}
	})

	t.Logf("✅ Parameter validation tests passed!")
}

// TestImageAPIGetImageTask tests the GetImageTask method with mock server
func TestImageAPIGetImageTask(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and path
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/api/image/task" {
			t.Errorf("Expected path /api/image/task, got %s", r.URL.Path)
		}

		// Verify query parameters
		loginToken := r.URL.Query().Get("loginToken")
		sessionId := r.URL.Query().Get("sessionId")
		taskId := r.URL.Query().Get("taskId")
		if loginToken != "test-login-token" {
			t.Errorf("Expected loginToken test-login-token, got %s", loginToken)
		}
		if sessionId != "test-session-id" {
			t.Errorf("Expected sessionId test-session-id, got %s", sessionId)
		}
		if taskId != "test-task-id" {
			t.Errorf("Expected taskId test-task-id, got %s", taskId)
		}

		// Mock successful response based on actual API response
		response := client.ImageTaskResponse{
			Code:      "success",
			RequestID: "test-task-request-id",
			Success:   true,
			Data: client.ImageTaskData{
				Status:  "Inline",
				TaskMsg: "Preparing Init Dockerfile",
				ImageID: nil,
			},
			TraceID:        "test-task-trace-id",
			HTTPStatusCode: 200,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client configuration
	cfg := client.NewConfiguration()
	cfg.Servers[0].URL = server.URL

	// Create API client
	apiClient := client.NewAPIClient(cfg)

	// Test context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Call GetImageTask
	response, httpResp, err := apiClient.ImageAPI.GetImageTask(ctx, "test-login-token", "test-session-id", "test-task-id")

	// Verify no error
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify HTTP response
	if httpResp.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", httpResp.StatusCode)
	}

	// Verify response data
	if !response.Success {
		t.Error("Expected success to be true")
	}
	if response.Code != "success" {
		t.Errorf("Expected code 'success', got '%s'", response.Code)
	}
	if response.Data.Status != "Inline" {
		t.Errorf("Expected Status 'Inline', got '%s'", response.Data.Status)
	}
	if response.Data.TaskMsg != "Preparing Init Dockerfile" {
		t.Errorf("Expected TaskMsg 'Preparing Init Dockerfile', got '%s'", response.Data.TaskMsg)
	}
	if response.Data.ImageID != nil {
		t.Errorf("Expected ImageID to be nil, got '%s'", *response.Data.ImageID)
	}

	t.Logf("✅ GetImageTask test passed!")
	t.Logf("   - Status: %s", response.Data.Status)
	t.Logf("   - TaskMsg: %s", response.Data.TaskMsg)
	t.Logf("   - ImageID: %v", response.Data.ImageID)
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || (len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsAt(s, substr))))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
