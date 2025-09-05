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
	GetLoginProviderURL(ctx context.Context, fromUrlPath, loginClient, oauthProvider string) (OAuthLoginProviderResponse, *http.Response, error)
	LoginTranslate(ctx context.Context, loginClient, oauthProvider, authCode string) (OAuthLoginTranslateResponse, *http.Response, error)
	RefreshToken(ctx context.Context, keepAliveToken, sessionId string) (OAuthRefreshTokenResponse, *http.Response, error)
	Logout(ctx context.Context, loginToken, sessionId string) (OAuthLogoutResponse, *http.Response, error)
}

// OAuthAPIService implements OAuthAPI interface
type OAuthAPIService service

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

// OAuthLoginTranslateResponse represents the response from OAuth login translate API
type OAuthLoginTranslateResponse struct {
	Code           string                  `json:"code"`
	RequestID      string                  `json:"requestId"`
	Success        bool                    `json:"success"`
	Data           OAuthLoginTranslateData `json:"data"`
	TraceID        string                  `json:"traceId"`
	HTTPStatusCode int                     `json:"httpStatusCode"`
}

// OAuthLoginTranslateData represents the data field in OAuth login translate response
// This matches the actual AgbCloud API response format
type OAuthLoginTranslateData struct {
	LoginToken     string `json:"loginToken"`
	SessionId      string `json:"sessionId"`
	KeepAliveToken string `json:"keepAliveToken"`
	ExpiresAt      string `json:"expiresAt"`
}

// OAuthRefreshTokenResponse represents the response from OAuth refresh token API
type OAuthRefreshTokenResponse struct {
	Code           string                `json:"code"`
	RequestID      string                `json:"requestId"`
	Success        bool                  `json:"success"`
	Data           OAuthRefreshTokenData `json:"data"`
	TraceID        string                `json:"traceId"`
	HTTPStatusCode int                   `json:"httpStatusCode"`
}

// OAuthRefreshTokenData represents the data field in OAuth refresh token response
// This matches the actual AgbCloud API response format
type OAuthRefreshTokenData struct {
	LoginToken     string `json:"loginToken"`
	SessionId      string `json:"sessionId"`
	KeepAliveToken string `json:"keepAliveToken"`
	ExpiresAt      string `json:"expiresAt"`
}

// OAuthLogoutResponse represents the response from OAuth logout API
type OAuthLogoutResponse struct {
	Code           string          `json:"code"`
	RequestID      string          `json:"requestId"`
	Success        bool            `json:"success"`
	Data           OAuthLogoutData `json:"data"`
	TraceID        string          `json:"traceId"`
	HTTPStatusCode int             `json:"httpStatusCode"`
}

// OAuthLogoutData represents the data field in OAuth logout response
type OAuthLogoutData struct {
	Message string `json:"message"`
}

// GetLoginProviderURL retrieves the OAuth login provider URL with loginClient parameter
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

	// Only add loginClient parameter, let backend use defaults for others
	if loginClient == "" {
		loginClient = "CLI"
	}
	localVarQueryParams.Add("loginClient", loginClient)

	// Only add oauthProvider if explicitly provided (let backend use default)
	if oauthProvider != "" {
		localVarQueryParams.Add("oauthProvider", oauthProvider)
	}

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

// RefreshToken refreshes the login session using keepAliveToken and sessionId
func (o *OAuthAPIService) RefreshToken(ctx context.Context, keepAliveToken, sessionId string) (OAuthRefreshTokenResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		localVarReturnValue OAuthRefreshTokenResponse
	)

	// Build the request path
	localVarPath := "/api/biz_login/refresh"

	// Use the configured server URL (defaults to agb.cloud)
	serverURL, err := o.client.cfg.ServerURLWithContext(ctx, "RefreshToken")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath = serverURL + localVarPath

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}

	// Set headers
	localVarHeaderParams["Accept"] = "application/json"

	// Add required query parameters
	if keepAliveToken == "" {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: "keepAliveToken parameter is required"}
	}
	localVarQueryParams.Add("keepAliveToken", keepAliveToken)

	if sessionId == "" {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: "sessionId parameter is required"}
	}
	localVarQueryParams.Add("sessionId", sessionId)

	// Prepare request
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

// LoginTranslate translates OAuth authorization code to access token
func (o *OAuthAPIService) LoginTranslate(ctx context.Context, loginClient, oauthProvider, authCode string) (OAuthLoginTranslateResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		localVarReturnValue OAuthLoginTranslateResponse
	)

	// Build the request path
	localVarPath := "/api/oauth/auth_code/login_translate"

	// Use the configured server URL (defaults to agb.cloud)
	serverURL, err := o.client.cfg.ServerURLWithContext(ctx, "LoginTranslate")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath = serverURL + localVarPath

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}

	// Set headers
	localVarHeaderParams["Accept"] = "application/json"

	// Add required query parameters
	if loginClient == "" {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: "loginClient parameter is required"}
	}
	localVarQueryParams.Add("loginClient", loginClient)

	if oauthProvider == "" {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: "oauthProvider parameter is required"}
	}
	localVarQueryParams.Add("oauthProvider", oauthProvider)

	if authCode == "" {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: "authCode parameter is required"}
	}
	localVarQueryParams.Add("authCode", authCode)

	// Prepare request
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

// Logout logs out the user by invalidating the session
func (o *OAuthAPIService) Logout(ctx context.Context, loginToken, sessionId string) (OAuthLogoutResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		localVarReturnValue OAuthLogoutResponse
	)

	// Build the request path
	localVarPath := "/api/biz_login/logout"

	// Use the configured server URL (defaults to agb.cloud)
	serverURL, err := o.client.cfg.ServerURLWithContext(ctx, "Logout")
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
