// Copyright 2025 AgbCloud CLI Contributors
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// ImageAPI interface for image related operations
type ImageAPI interface {
	GetUploadCredential(ctx context.Context, loginToken, sessionId string) (ImageUploadCredentialResponse, *http.Response, error)
	CreateImage(ctx context.Context, loginToken, sessionId, imageName, taskId, sourceImageId string) (ImageCreateResponse, *http.Response, error)
	GetImageTask(ctx context.Context, loginToken, sessionId, taskId string) (ImageTaskResponse, *http.Response, error)
	ListImages(ctx context.Context, loginToken, sessionId, imageType string, page, pageSize int, imageIds []string) (ImageListResponse, *http.Response, error)
	StartImage(ctx context.Context, loginToken, sessionId, imageId string, cpu, memory int) (ImageStartResponse, *http.Response, error)
	StopImage(ctx context.Context, loginToken, sessionId, imageId string) (ImageStopResponse, *http.Response, error)
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

// ImageListResponse represents the response from /api/image/list API
type ImageListResponse struct {
	Code           string        `json:"code"`
	RequestID      string        `json:"requestId"`
	Success        bool          `json:"success"`
	Data           ImageListData `json:"data"`
	TraceID        string        `json:"traceId"`
	HTTPStatusCode int           `json:"httpStatusCode"`
}

// ImageListData represents the data field in image list response
type ImageListData struct {
	Images   []ImageInfo `json:"imageList"` // API returns "imageList" not "images"
	Total    int         `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}

// ImageInfo represents individual image information
type ImageInfo struct {
	ImageID      string  `json:"imageId"`
	ImageName    string  `json:"imageName"`
	Status       string  `json:"status"`
	Type         string  `json:"type"`
	OSType       string  `json:"osType"`
	UpdateTime   string  `json:"updateTime"`   // API uses "updateTime" not "updatedAt"
	GmtCreate    *string `json:"gmtCreate"`    // Can be null
	GmtUpdate    *string `json:"gmtUpdate"`    // Can be null
	LastUsedTime *string `json:"lastUsedTime"` // Can be null
	CPU          *int    `json:"cpu"`          // Can be null
	Memory       *int    `json:"memory"`       // Can be null, in GB
}

// ImageStartResponse represents the response from /api/image/start API
type ImageStartResponse struct {
	Code           string `json:"code"`
	RequestID      string `json:"requestId"`
	Success        bool   `json:"success"`
	Data           bool   `json:"data"`
	TraceID        string `json:"traceId"`
	HTTPStatusCode int    `json:"httpStatusCode"`
}

// ImageStartData represents the data field in image start response
// Note: This is kept for backward compatibility but not used in actual API response
type ImageStartData struct {
	InstanceID string `json:"instanceId"`
	Status     string `json:"status"`
}

// ImageStopResponse represents the response from /api/image/stop API
type ImageStopResponse struct {
	Code           string `json:"code"`
	RequestID      string `json:"requestId"`
	Success        bool   `json:"success"`
	Data           bool   `json:"data"`
	TraceID        string `json:"traceId"`
	HTTPStatusCode int    `json:"httpStatusCode"`
}

// ImageStopData represents the data field in image stop response
// Note: This is kept for backward compatibility but not used in actual API response
type ImageStopData struct {
	InstanceID string `json:"instanceId"`
	Status     string `json:"status"`
}

// ImageStartRequest represents the request body for /api/image/start API
type ImageStartRequest struct {
	LoginToken string `json:"loginToken"`
	SessionId  string `json:"sessionId"`
	ImageId    string `json:"imageId"`
	CPU        int    `json:"cpu,omitempty"`
	Memory     int    `json:"memory,omitempty"`
}

// ImageStopRequest represents the request body for /api/image/stop API
type ImageStopRequest struct {
	LoginToken string `json:"loginToken"`
	SessionId  string `json:"sessionId"`
	ImageId    string `json:"imageId"`
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

// ListImages retrieves a list of images with pagination, optionally filtered by image IDs
func (i *ImageAPIService) ListImages(ctx context.Context, loginToken, sessionId, imageType string, page, pageSize int, imageIds []string) (ImageListResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		localVarReturnValue ImageListResponse
	)

	// Build the request path
	localVarPath := "/api/image/list"

	// Use the configured server URL (defaults to agb.cloud)
	serverURL, err := i.client.cfg.ServerURLWithContext(ctx, "ListImages")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath = serverURL + localVarPath

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}

	// Set headers
	localVarHeaderParams["Accept"] = "application/json"

	// Validate required parameters
	if loginToken == "" {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: "loginToken parameter is required"}
	}
	localVarQueryParams.Add("loginToken", loginToken)

	if sessionId == "" {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: "sessionId parameter is required"}
	}
	localVarQueryParams.Add("sessionId", sessionId)

	if imageType == "" {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: "imageType parameter is required"}
	}
	localVarQueryParams.Add("imageType", imageType)

	// Validate pagination parameters
	if page <= 0 {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: "page must be greater than 0"}
	}
	localVarQueryParams.Add("page", fmt.Sprintf("%d", page))

	if pageSize <= 0 {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: "pageSize must be greater than 0"}
	}
	localVarQueryParams.Add("pageSize", fmt.Sprintf("%d", pageSize))

	// Add imageIds parameter if provided
	if len(imageIds) > 0 {
		for _, imageId := range imageIds {
			localVarQueryParams.Add("imageIds", imageId)
		}
	}

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

// StartImage starts an image with specified resources
func (i *ImageAPIService) StartImage(ctx context.Context, loginToken, sessionId, imageId string, cpu, memory int) (ImageStartResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarReturnValue ImageStartResponse
	)

	// Build the request path
	localVarPath := "/api/image/start"

	// Use the configured server URL (defaults to agb.cloud)
	serverURL, err := i.client.cfg.ServerURLWithContext(ctx, "StartImage")
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
	if imageId == "" {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: "imageId parameter is required"}
	}

	// Create request body
	requestBody := ImageStartRequest{
		LoginToken: loginToken,
		SessionId:  sessionId,
		ImageId:    imageId,
		CPU:        cpu,
		Memory:     memory,
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

// StopImage stops a running image instance
func (i *ImageAPIService) StopImage(ctx context.Context, loginToken, sessionId, imageId string) (ImageStopResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarReturnValue ImageStopResponse
	)

	// Build the request path
	localVarPath := "/api/image/stop"

	// Use the configured server URL (defaults to agb.cloud)
	serverURL, err := i.client.cfg.ServerURLWithContext(ctx, "StopImage")
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
	if imageId == "" {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: "imageId parameter is required"}
	}

	// Create request body
	requestBody := ImageStopRequest{
		LoginToken: loginToken,
		SessionId:  sessionId,
		ImageId:    imageId,
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
