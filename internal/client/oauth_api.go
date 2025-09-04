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
	GetLoginProviderURL(ctx context.Context, fromUrlPath, loginClient, oauthProvider string) (OAuthLoginProviderResponse, *http.Response, error)
}

// OAuthAPIService implements OAuthAPI interface
type OAuthAPIService service

// OAuthGoogleLoginResponse represents the response from Google OAuth login API (legacy)
type OAuthGoogleLoginResponse struct {
	Code           string               `json:"code"`
	RequestID      string               `json:"requestId"`
	Success        bool                 `json:"success"`
	Data           OAuthGoogleLoginData `json:"data"`
	TraceID        string               `json:"traceId"`
	HTTPStatusCode int                  `json:"httpStatusCode"`
}

// OAuthGoogleLoginData represents the data field in Google OAuth login response (legacy)
type OAuthGoogleLoginData struct {
	InvokeURL string `json:"invokeUrl"`
}

// OAuthLoginProviderResponse represents the response from OAuth login provider API
type OAuthLoginProviderResponse struct {
	Code           string                 `json:"code"`
	RequestID      string                 `json:"requestId"`
	Success        bool                   `json:"success"`
	Data           OAuthLoginProviderData `json:"data"`
	TraceID        string                 `json:"traceId"`
	HTTPStatusCode int                    `json:"httpStatusCode"`
}

// OAuthLoginProviderData represents the data field in OAuth login provider response
type OAuthLoginProviderData struct {
	InvokeURL string `json:"invokeUrl"`
}

// GetGoogleLoginURL retrieves the Google OAuth login URL (legacy method)
func (o *OAuthAPIService) GetGoogleLoginURL(ctx context.Context, fromUrlPath string) (OAuthGoogleLoginResponse, *http.Response, error) {
	// Convert to new API call with default values
	newResponse, httpResp, err := o.GetLoginProviderURL(ctx, fromUrlPath, "CLI", "GOOGLE_LOCALHOST")
	if err != nil {
		return OAuthGoogleLoginResponse{}, httpResp, err
	}

	// Convert new response to legacy response format
	legacyResponse := OAuthGoogleLoginResponse{
		Code:      newResponse.Code,
		RequestID: newResponse.RequestID,
		Success:   newResponse.Success,
		Data: OAuthGoogleLoginData{
			InvokeURL: newResponse.Data.InvokeURL,
		},
		TraceID:        newResponse.TraceID,
		HTTPStatusCode: newResponse.HTTPStatusCode,
	}

	return legacyResponse, httpResp, nil
}

// GetLoginProviderURL retrieves the OAuth login provider URL with new parameters
func (o *OAuthAPIService) GetLoginProviderURL(ctx context.Context, fromUrlPath, loginClient, oauthProvider string) (OAuthLoginProviderResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		localVarReturnValue OAuthLoginProviderResponse
	)

	// Build the request path
	localVarPath := "/api/oauth/login_provider"

	// Use the configured server URL (defaults to agb.cloud)
	serverURL, err := o.client.cfg.ServerURLWithContext(ctx, "GetLoginProviderURL")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath = serverURL + localVarPath

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}

	// Set headers
	localVarHeaderParams["Accept"] = "application/json"

	// Note: OAuth login provider endpoint does not require authorization

	// Add query parameters
	if fromUrlPath != "" {
		localVarQueryParams.Add("fromUrlPath", fromUrlPath)
	}

	// Set default values if not provided
	if loginClient == "" {
		loginClient = "CLI"
	}
	if oauthProvider == "" {
		oauthProvider = "GOOGLE_LOCALHOST"
	}

	localVarQueryParams.Add("loginClient", loginClient)
	localVarQueryParams.Add("oauthProvider", oauthProvider)

	// Use the original context but without authentication for OAuth endpoint
	// Note: We don't need to strip authentication since we're not adding any auth headers
	req, err := o.client.prepareRequest(ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams)
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
