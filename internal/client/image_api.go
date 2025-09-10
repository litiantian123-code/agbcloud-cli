// Copyright 2025 AgbCloud CLI Contributors
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
)

// ImageAPI interface for image related operations
type ImageAPI interface {
	GetUploadCredential(ctx context.Context, loginToken, sessionId string) (ImageUploadCredentialResponse, *http.Response, error)
	CreateImage(ctx context.Context, loginToken, sessionId, imageName, taskId, sourceImageId string) (ImageCreateResponse, *http.Response, error)
	GetImageTask(ctx context.Context, loginToken, sessionId, taskId string) (ImageTaskResponse, *http.Response, error)
}

// ImageAPIService implements ImageAPI interface
type ImageAPIService service

// ImageUploadCredentialResponse represents the response from /api/image/getUploadCredential API
type ImageUploadCredentialResponse struct {
	Code           string                    `json:"code"`
	RequestID      string                    `json:"requestId"`
	Success        bool                      `json:"success"`
	Data           ImageUploadCredentialData `json:"data"`
	TraceID        string                    `json:"traceId"`
	HTTPStatusCode int                       `json:"httpStatusCode"`
}

// ImageUploadCredentialData represents the data field in image upload credential response
// This structure matches the actual API response structure
type ImageUploadCredentialData struct {
	OssURL string `json:"ossUrl"`
	TaskID string `json:"taskId"`
}

// ImageCreateResponse represents the response from /api/image/create API
type ImageCreateResponse struct {
	Code           string `json:"code"`
	RequestID      string `json:"requestId"`
	Success        bool   `json:"success"`
	Data           string `json:"data"`
	TraceID        string `json:"traceId"`
	HTTPStatusCode int    `json:"httpStatusCode"`
}

// ImageTaskResponse represents the response from /api/image/task API
type ImageTaskResponse struct {
	Code           string        `json:"code"`
	RequestID      string        `json:"requestId"`
	Success        bool          `json:"success"`
	Data           ImageTaskData `json:"data"`
	TraceID        string        `json:"traceId"`
	HTTPStatusCode int           `json:"httpStatusCode"`
}

// ImageTaskData represents the data field in image task response
// This structure matches the actual API response structure
type ImageTaskData struct {
	Status  string  `json:"status"`
	TaskMsg string  `json:"taskMsg"`
	ImageID *string `json:"imageId"`
}

// ImageCreateData represents the data field in image create response
type ImageCreateData struct {
	ImageID   string `json:"imageId,omitempty"`
	ImageName string `json:"imageName,omitempty"`
	Status    string `json:"status,omitempty"`
	Message   string `json:"message,omitempty"`
}

// GetUploadCredential retrieves upload credentials for image upload
func (i *ImageAPIService) GetUploadCredential(ctx context.Context, loginToken, sessionId string) (ImageUploadCredentialResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		localVarReturnValue ImageUploadCredentialResponse
	)

	// Build the request path
	localVarPath := "/api/image/getUploadCredential"

	// Use the configured server URL (defaults to agb.cloud)
	serverURL, err := i.client.cfg.ServerURLWithContext(ctx, "GetUploadCredential")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath = serverURL + localVarPath

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}

	// Set headers
	localVarHeaderParams["Accept"] = "application/json"

	// Add required query parameters
	if loginToken == "" {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: "loginToken parameter is required"}
	}
	localVarQueryParams.Add("loginToken", loginToken)

	if sessionId == "" {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: "sessionId parameter is required"}
	}
	localVarQueryParams.Add("sessionId", sessionId)

	// Prepare request
	req, err := i.client.prepareRequest(ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := i.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = i.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

// CreateImage creates a custom image using the uploaded dockerfile
func (i *ImageAPIService) CreateImage(ctx context.Context, loginToken, sessionId, imageName, taskId, sourceImageId string) (ImageCreateResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarReturnValue ImageCreateResponse
	)

	// Build the request path
	localVarPath := "/api/image/create"

	// Use the configured server URL (defaults to agb.cloud)
	serverURL, err := i.client.cfg.ServerURLWithContext(ctx, "CreateImage")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath = serverURL + localVarPath

	localVarHeaderParams := make(map[string]string)

	// Set headers
	localVarHeaderParams["Accept"] = "application/json"
	localVarHeaderParams["Content-Type"] = "application/json"

	// Validate required parameters
	if loginToken == "" {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: "loginToken parameter is required"}
	}
	if sessionId == "" {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: "sessionId parameter is required"}
	}
	if imageName == "" {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: "imageName parameter is required"}
	}
	if taskId == "" {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: "taskId parameter is required"}
	}
	if sourceImageId == "" {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: "sourceImageId parameter is required"}
	}

	// Create request body
	requestBody := map[string]string{
		"loginToken":    loginToken,
		"sessionId":     sessionId,
		"imageName":     imageName,
		"taskId":        taskId,
		"sourceImageId": sourceImageId,
	}

	// Prepare request
	req, err := i.client.prepareRequest(ctx, localVarPath, localVarHTTPMethod, requestBody, localVarHeaderParams, url.Values{})
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := i.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = i.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

// GetImageTask retrieves the status of an image creation task
func (i *ImageAPIService) GetImageTask(ctx context.Context, loginToken, sessionId, taskId string) (ImageTaskResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		localVarReturnValue ImageTaskResponse
	)

	// Build the request path
	localVarPath := "/api/image/task"

	// Use the configured server URL (defaults to agb.cloud)
	serverURL, err := i.client.cfg.ServerURLWithContext(ctx, "GetImageTask")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath = serverURL + localVarPath

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}

	// Set headers
	localVarHeaderParams["Accept"] = "application/json"

	// Add required query parameters
	if loginToken == "" {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: "loginToken parameter is required"}
	}
	localVarQueryParams.Add("loginToken", loginToken)

	if sessionId == "" {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: "sessionId parameter is required"}
	}
	localVarQueryParams.Add("sessionId", sessionId)

	if taskId == "" {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: "taskId parameter is required"}
	}
	localVarQueryParams.Add("taskId", taskId)

	// Prepare request
	req, err := i.client.prepareRequest(ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := i.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = i.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}
