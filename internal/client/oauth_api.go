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

// OAuthAPI interface for OAuth related operations
type OAuthAPI interface {
	GetGoogleLoginURL(ctx context.Context, fromUrlPath string) (OAuthGoogleLoginResponse, *http.Response, error)
}

// OAuthAPIService implements OAuthAPI interface
type OAuthAPIService service

// OAuthGoogleLoginResponse represents the response from Google OAuth login API
type OAuthGoogleLoginResponse struct {
	Code           string               `json:"code"`
	RequestID      string               `json:"requestId"`
	Success        bool                 `json:"success"`
	Data           OAuthGoogleLoginData `json:"data"`
	TraceID        string               `json:"traceId"`
	HTTPStatusCode int                  `json:"httpStatusCode"`
}

// OAuthGoogleLoginData represents the data field in Google OAuth login response
type OAuthGoogleLoginData struct {
	InvokeURL string `json:"invokeUrl"`
}

// GetGoogleLoginURL retrieves the Google OAuth login URL
func (o *OAuthAPIService) GetGoogleLoginURL(ctx context.Context, fromUrlPath string) (OAuthGoogleLoginResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		localVarReturnValue OAuthGoogleLoginResponse
	)

	// Build the request path
	localVarPath := "/api/oauth/google/login"

	// Use the configured server URL (defaults to agb.cloud)
	serverURL, err := o.client.cfg.ServerURLWithContext(ctx, "GetGoogleLoginURL")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath = serverURL + localVarPath

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}

	// Set headers
	localVarHeaderParams["Accept"] = "application/json"

	// Note: OAuth Google login endpoint does not require authorization

	// Add fromUrlPath query parameter
	if fromUrlPath != "" {
		localVarQueryParams.Add("fromUrlPath", fromUrlPath)
	}

	// Create a context without authentication for OAuth endpoint
	ctxWithoutAuth := context.Background()
	if deadline, ok := ctx.Deadline(); ok {
		var cancel context.CancelFunc
		ctxWithoutAuth, cancel = context.WithDeadline(ctxWithoutAuth, deadline)
		defer cancel()
	}

	req, err := o.client.prepareRequest(ctxWithoutAuth, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := o.client.callAPI(req)
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

	err = o.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}
