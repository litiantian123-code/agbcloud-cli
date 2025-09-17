// Copyright 2025 AgbCloud CLI Contributors
// SPDX-License-Identifier: Apache-2.0

package unit

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"errors"

	"github.com/agbcloud/agbcloud-cli/internal/client"
	"github.com/stretchr/testify/assert"
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

	t.Logf("[OK] GetUploadCredential test passed!")
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

	t.Logf("[OK] CreateImage test passed!")
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

	t.Logf("[OK] Parameter validation tests passed!")
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

	t.Logf("[OK] GetImageTask test passed!")
	t.Logf("   - Status: %s", response.Data.Status)
	t.Logf("   - TaskMsg: %s", response.Data.TaskMsg)
	t.Logf("   - ImageID: %v", response.Data.ImageID)
}

// TestImageAPIListImages tests the ListImages method with mock server
func TestImageAPIListImages(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and path
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/api/image/list" {
			t.Errorf("Expected path /api/image/list, got %s", r.URL.Path)
		}

		// Verify query parameters
		loginToken := r.URL.Query().Get("loginToken")
		sessionId := r.URL.Query().Get("sessionId")
		imageType := r.URL.Query().Get("imageType")
		page := r.URL.Query().Get("page")
		pageSize := r.URL.Query().Get("pageSize")

		if loginToken != "test-login-token" {
			t.Errorf("Expected loginToken test-login-token, got %s", loginToken)
		}
		if sessionId != "test-session-id" {
			t.Errorf("Expected sessionId test-session-id, got %s", sessionId)
		}
		if imageType != "User" {
			t.Errorf("Expected imageType User, got %s", imageType)
		}
		if page != "1" {
			t.Errorf("Expected page 1, got %s", page)
		}
		if pageSize != "10" {
			t.Errorf("Expected pageSize 10, got %s", pageSize)
		}

		// Mock successful response
		response := client.ImageListResponse{
			Code:      "success",
			RequestID: "test-list-request-id",
			Success:   true,
			Data: client.ImageListData{
				Images: []client.ImageInfo{
					{
						ImageID:    "img-12345",
						ImageName:  "my-custom-image",
						Status:     "IMAGE_AVAILABLE",
						Type:       "CodeSpace",
						OSType:     "Linux",
						UpdateTime: "2025-01-01T10:30:00Z",
					},
					{
						ImageID:    "img-67890",
						ImageName:  "another-image",
						Status:     "IMAGE_BUILDING",
						Type:       "CodeSpace",
						OSType:     "Linux",
						UpdateTime: "2025-01-01T11:15:00Z",
					},
				},
				Total:    2,
				Page:     1,
				PageSize: 10,
			},
			TraceID:        "test-list-trace-id",
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

	// Call ListImages
	response, httpResp, err := apiClient.ImageAPI.ListImages(ctx, "test-login-token", "test-session-id", "User", 1, 10, nil)

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
	if response.Data.Total != 2 {
		t.Errorf("Expected Total 2, got %d", response.Data.Total)
	}
	if response.Data.Page != 1 {
		t.Errorf("Expected Page 1, got %d", response.Data.Page)
	}
	if response.Data.PageSize != 10 {
		t.Errorf("Expected PageSize 10, got %d", response.Data.PageSize)
	}
	if len(response.Data.Images) != 2 {
		t.Errorf("Expected 2 images, got %d", len(response.Data.Images))
	}

	// Verify first image
	firstImage := response.Data.Images[0]
	if firstImage.ImageID != "img-12345" {
		t.Errorf("Expected first image ID 'img-12345', got '%s'", firstImage.ImageID)
	}
	if firstImage.ImageName != "my-custom-image" {
		t.Errorf("Expected first image name 'my-custom-image', got '%s'", firstImage.ImageName)
	}
	if firstImage.Status != "IMAGE_AVAILABLE" {
		t.Errorf("Expected first image status 'IMAGE_AVAILABLE', got '%s'", firstImage.Status)
	}

	t.Logf("[OK] ListImages test passed!")
	t.Logf("   - Total images: %d", response.Data.Total)
	t.Logf("   - Page: %d, PageSize: %d", response.Data.Page, response.Data.PageSize)
	t.Logf("   - First image: %s (%s)", firstImage.ImageName, firstImage.Status)
}

// TestImageAPIListImagesParameterValidation tests parameter validation for ListImages
func TestImageAPIListImagesParameterValidation(t *testing.T) {
	// Create client configuration
	cfg := client.NewConfiguration()
	cfg.Servers[0].URL = "http://localhost:8080" // Dummy URL

	// Create API client
	apiClient := client.NewAPIClient(cfg)

	// Test context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	t.Run("ListImages_EmptyLoginToken", func(t *testing.T) {
		_, _, err := apiClient.ImageAPI.ListImages(ctx, "", "test-session-id", "User", 1, 10, nil)
		if err == nil {
			t.Error("Expected error for empty loginToken")
		}
		if !contains(err.Error(), "loginToken parameter is required") {
			t.Errorf("Expected loginToken error message, got: %s", err.Error())
		}
	})

	t.Run("ListImages_EmptySessionId", func(t *testing.T) {
		_, _, err := apiClient.ImageAPI.ListImages(ctx, "test-login-token", "", "User", 1, 10, nil)
		if err == nil {
			t.Error("Expected error for empty sessionId")
		}
		if !contains(err.Error(), "sessionId parameter is required") {
			t.Errorf("Expected sessionId error message, got: %s", err.Error())
		}
	})

	t.Run("ListImages_EmptyImageType", func(t *testing.T) {
		_, _, err := apiClient.ImageAPI.ListImages(ctx, "test-login-token", "test-session-id", "", 1, 10, nil)
		if err == nil {
			t.Error("Expected error for empty imageType")
		}
		if !contains(err.Error(), "imageType parameter is required") {
			t.Errorf("Expected imageType error message, got: %s", err.Error())
		}
	})

	t.Run("ListImages_InvalidPage", func(t *testing.T) {
		_, _, err := apiClient.ImageAPI.ListImages(ctx, "test-login-token", "test-session-id", "User", 0, 10, nil)
		if err == nil {
			t.Error("Expected error for invalid page")
		}
		if !contains(err.Error(), "page must be greater than 0") {
			t.Errorf("Expected page validation error message, got: %s", err.Error())
		}
	})

	t.Run("ListImages_InvalidPageSize", func(t *testing.T) {
		_, _, err := apiClient.ImageAPI.ListImages(ctx, "test-login-token", "test-session-id", "User", 1, 0, nil)
		if err == nil {
			t.Error("Expected error for invalid pageSize")
		}
		if !contains(err.Error(), "pageSize must be greater than 0") {
			t.Errorf("Expected pageSize validation error message, got: %s", err.Error())
		}
	})

	t.Logf("[OK] ListImages parameter validation tests passed!")
}

// TestImageAPIService_StartImage tests the StartImage method with mock server
func TestImageAPIService_StartImage(t *testing.T) {
	// Test successful start image request
	t.Run("Success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify request method and path
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/api/image/start", r.URL.Path)

			// Verify headers
			assert.Equal(t, "application/json", r.Header.Get("Accept"))
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

			// Verify request body
			body, err := io.ReadAll(r.Body)
			assert.NoError(t, err)

			var requestBody client.ImageStartRequest
			err = json.Unmarshal(body, &requestBody)
			assert.NoError(t, err)

			assert.Equal(t, "test-login-token", requestBody.LoginToken)
			assert.Equal(t, "test-session-id", requestBody.SessionId)
			assert.Equal(t, "test-image-id", requestBody.ImageId)
			assert.Equal(t, 2, requestBody.CPU)
			assert.Equal(t, 4, requestBody.Memory)

			// Send success response
			response := client.ImageStartResponse{
				Code:           "SUCCESS",
				RequestID:      "req-12345",
				Success:        true,
				Data:           true,
				TraceID:        "trace-12345",
				HTTPStatusCode: 200,
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		// Create client with test server
		cfg := client.NewConfiguration()
		cfg.Servers[0].URL = server.URL
		apiClient := client.NewAPIClient(cfg)

		// Call StartImage API
		ctx := context.Background()
		resp, httpResp, err := apiClient.ImageAPI.StartImage(ctx, "test-login-token", "test-session-id", "test-image-id", 2, 4)

		// Verify results
		assert.NoError(t, err)
		assert.NotNil(t, httpResp)
		assert.Equal(t, http.StatusOK, httpResp.StatusCode)
		assert.True(t, resp.Success)
		assert.Equal(t, "SUCCESS", resp.Code)
		assert.Equal(t, "req-12345", resp.RequestID)
		assert.True(t, resp.Data)
	})

	// Test with optional parameters as zero values (should be included in request body as 0)
	t.Run("SuccessWithOptionalParametersZero", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify request body
			body, err := io.ReadAll(r.Body)
			assert.NoError(t, err)

			var requestBody client.ImageStartRequest
			err = json.Unmarshal(body, &requestBody)
			assert.NoError(t, err)

			assert.Equal(t, "test-login-token", requestBody.LoginToken)
			assert.Equal(t, "test-session-id", requestBody.SessionId)
			assert.Equal(t, "test-image-id", requestBody.ImageId)

			// Verify optional parameters are included as zero values in JSON body
			assert.Equal(t, 0, requestBody.CPU)
			assert.Equal(t, 0, requestBody.Memory)

			// Send success response
			response := client.ImageStartResponse{
				Code:           "SUCCESS",
				RequestID:      "req-12345",
				Success:        true,
				Data:           true,
				TraceID:        "trace-12345",
				HTTPStatusCode: 200,
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		// Create client with test server
		cfg := client.NewConfiguration()
		cfg.Servers[0].URL = server.URL
		apiClient := client.NewAPIClient(cfg)

		// Call StartImage API with zero values for optional parameters
		ctx := context.Background()
		resp, httpResp, err := apiClient.ImageAPI.StartImage(ctx, "test-login-token", "test-session-id", "test-image-id", 0, 0)

		// Verify results
		assert.NoError(t, err)
		assert.NotNil(t, httpResp)
		assert.Equal(t, http.StatusOK, httpResp.StatusCode)
		assert.True(t, resp.Success)
	})

	// Test missing required parameters
	t.Run("MissingLoginToken", func(t *testing.T) {
		cfg := client.NewConfiguration()
		apiClient := client.NewAPIClient(cfg)

		ctx := context.Background()
		_, _, err := apiClient.ImageAPI.StartImage(ctx, "", "test-session-id", "test-image-id", 2, 4)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "loginToken parameter is required")
	})

	t.Run("MissingSessionId", func(t *testing.T) {
		cfg := client.NewConfiguration()
		apiClient := client.NewAPIClient(cfg)

		ctx := context.Background()
		_, _, err := apiClient.ImageAPI.StartImage(ctx, "test-login-token", "", "test-image-id", 2, 4)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sessionId parameter is required")
	})

	t.Run("MissingImageId", func(t *testing.T) {
		cfg := client.NewConfiguration()
		apiClient := client.NewAPIClient(cfg)

		ctx := context.Background()
		_, _, err := apiClient.ImageAPI.StartImage(ctx, "test-login-token", "test-session-id", "", 2, 4)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "imageId parameter is required")
	})

	// Test API error response
	t.Run("APIError", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response := client.ImageStartResponse{
				Code:      "INVALID_IMAGE_ID",
				RequestID: "req-error-123",
				Success:   false,
				Data:      false,
				TraceID:   "trace-error-123",
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		cfg := client.NewConfiguration()
		cfg.Servers[0].URL = server.URL
		apiClient := client.NewAPIClient(cfg)

		ctx := context.Background()
		_, httpResp, err := apiClient.ImageAPI.StartImage(ctx, "test-login-token", "test-session-id", "invalid-image-id", 2, 4)

		// Should return error due to HTTP status >= 300
		assert.Error(t, err)
		assert.NotNil(t, httpResp)
		assert.Equal(t, http.StatusBadRequest, httpResp.StatusCode)

		// But we should still be able to parse the response
		var apiError *client.GenericOpenAPIError
		assert.True(t, errors.As(err, &apiError))
		assert.NotNil(t, apiError)
	})

	// Test network error
	t.Run("NetworkError", func(t *testing.T) {
		cfg := client.NewConfiguration()
		cfg.Servers[0].URL = "http://non-existent-server:9999"
		apiClient := client.NewAPIClient(cfg)

		ctx := context.Background()
		_, _, err := apiClient.ImageAPI.StartImage(ctx, "test-login-token", "test-session-id", "test-image-id", 2, 4)

		assert.Error(t, err)
		// Should be a network error, not an API error
		var apiError *client.GenericOpenAPIError
		assert.False(t, errors.As(err, &apiError))
	})
}

// TestImageAPIService_StopImage tests the StopImage method with mock server
func TestImageAPIService_StopImage(t *testing.T) {
	// Test successful stop image request
	t.Run("Success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify request method and path
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/api/image/stop", r.URL.Path)

			// Verify headers
			assert.Equal(t, "application/json", r.Header.Get("Accept"))
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

			// Verify request body
			body, err := io.ReadAll(r.Body)
			assert.NoError(t, err)

			var requestBody client.ImageStopRequest
			err = json.Unmarshal(body, &requestBody)
			assert.NoError(t, err)

			assert.Equal(t, "test-login-token", requestBody.LoginToken)
			assert.Equal(t, "test-session-id", requestBody.SessionId)
			assert.Equal(t, "test-image-id", requestBody.ImageId)

			// Send success response
			response := client.ImageStopResponse{
				Code:           "SUCCESS",
				RequestID:      "req-stop-12345",
				Success:        true,
				Data:           true,
				TraceID:        "trace-stop-12345",
				HTTPStatusCode: 200,
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		// Create client with test server
		cfg := client.NewConfiguration()
		cfg.Servers[0].URL = server.URL
		apiClient := client.NewAPIClient(cfg)

		// Call StopImage API
		ctx := context.Background()
		resp, httpResp, err := apiClient.ImageAPI.StopImage(ctx, "test-login-token", "test-session-id", "test-image-id")

		// Verify results
		assert.NoError(t, err)
		assert.NotNil(t, httpResp)
		assert.Equal(t, http.StatusOK, httpResp.StatusCode)
		assert.True(t, resp.Success)
		assert.Equal(t, "SUCCESS", resp.Code)
		assert.Equal(t, "req-stop-12345", resp.RequestID)
		assert.True(t, resp.Data)

	})

	// Test missing required parameters
	t.Run("MissingLoginToken", func(t *testing.T) {
		cfg := client.NewConfiguration()
		apiClient := client.NewAPIClient(cfg)

		ctx := context.Background()
		_, _, err := apiClient.ImageAPI.StopImage(ctx, "", "test-session-id", "test-image-id")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "loginToken parameter is required")
	})

	t.Run("MissingSessionId", func(t *testing.T) {
		cfg := client.NewConfiguration()
		apiClient := client.NewAPIClient(cfg)

		ctx := context.Background()
		_, _, err := apiClient.ImageAPI.StopImage(ctx, "test-login-token", "", "test-image-id")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sessionId parameter is required")
	})

	t.Run("MissingImageId", func(t *testing.T) {
		cfg := client.NewConfiguration()
		apiClient := client.NewAPIClient(cfg)

		ctx := context.Background()
		_, _, err := apiClient.ImageAPI.StopImage(ctx, "test-login-token", "test-session-id", "")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "imageId parameter is required")
	})

	// Test API error response
	t.Run("APIError", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response := client.ImageStopResponse{
				Code:      "INVALID_IMAGE_ID",
				RequestID: "req-error-stop-123",
				Success:   false,
				Data:      false,
				TraceID:   "trace-error-stop-123",
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		cfg := client.NewConfiguration()
		cfg.Servers[0].URL = server.URL
		apiClient := client.NewAPIClient(cfg)

		ctx := context.Background()
		_, httpResp, err := apiClient.ImageAPI.StopImage(ctx, "test-login-token", "test-session-id", "invalid-image-id")

		// Should return error due to HTTP status >= 300
		assert.Error(t, err)
		assert.NotNil(t, httpResp)
		assert.Equal(t, http.StatusBadRequest, httpResp.StatusCode)

		// But we should still be able to parse the response
		var apiError *client.GenericOpenAPIError
		assert.True(t, errors.As(err, &apiError))
		assert.NotNil(t, apiError)
	})

	// Test network error
	t.Run("NetworkError", func(t *testing.T) {
		cfg := client.NewConfiguration()
		cfg.Servers[0].URL = "http://non-existent-server:9999"
		apiClient := client.NewAPIClient(cfg)

		ctx := context.Background()
		_, _, err := apiClient.ImageAPI.StopImage(ctx, "test-login-token", "test-session-id", "test-image-id")

		assert.Error(t, err)
		// Should be a network error, not an API error
		var apiError *client.GenericOpenAPIError
		assert.False(t, errors.As(err, &apiError))
	})
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
