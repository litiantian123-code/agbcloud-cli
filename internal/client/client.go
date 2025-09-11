// Copyright 2025 AgbCloud CLI Contributors
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	JsonCheck = regexp.MustCompile(`(?i:(?:application|text)/(?:[^;]+\+)?json)`)
)

// APIClient manages communication with the AgbCloud API
type APIClient struct {
	cfg    *Configuration
	common service // Reuse a single struct instead of allocating one for each service on the heap.

	// API Services
	OAuthAPI OAuthAPI
	ImageAPI ImageAPI
}

type service struct {
	client *APIClient
}

// NewAPIClient creates a new API client
func NewAPIClient(cfg *Configuration) *APIClient {
	if cfg.HTTPClient == nil {
		cfg.HTTPClient = &http.Client{
			Timeout: 30 * time.Second,
		}
	}

	c := &APIClient{}
	c.cfg = cfg
	c.common.client = c

	// API Services
	c.OAuthAPI = (*OAuthAPIService)(&c.common)
	c.ImageAPI = (*ImageAPIService)(&c.common)

	return c
}

// GetConfig returns the configuration
func (c *APIClient) GetConfig() *Configuration {
	return c.cfg
}

// callAPI do the request.
func (c *APIClient) callAPI(request *http.Request) (*http.Response, error) {
	// Log request information for debugging (only shown with -v flag)
	log.Debugf("\n=== HTTP Request Information ===")
	log.Debugf("URL: %s", request.URL.String())
	log.Debugf("Method: %s", request.Method)
	log.Debugf("Headers: %v", request.Header)

	if request.Body != nil {
		// Read body for logging without consuming it
		bodyBytes, err := io.ReadAll(request.Body)
		if err == nil {
			log.Debugf("Request Body: %s", string(bodyBytes))
			// Restore the body
			request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}
	} else {
		log.Debugf("Request Body: None")
	}
	log.Debugf("=" + strings.Repeat("=", 49))

	if c.cfg.Debug {
		dump, err := httputil.DumpRequestOut(request, true)
		if err != nil {
			return nil, err
		}
		log.Debugf("\n%s\n", string(dump))
	}

	resp, err := c.cfg.HTTPClient.Do(request)
	if err != nil {
		log.Debugf("\n=== HTTP Request Error ===")
		log.Debugf("Error Type: %T", err)
		log.Debugf("Error Message: %s", err.Error())
		log.Debugf("Request URL: %s", request.URL.String())
		log.Debugf("=" + strings.Repeat("=", 49))
		return resp, err
	}

	// Log response information for debugging (only shown with -v flag)
	log.Debugf("\n=== HTTP Response Information ===")
	log.Debugf("Status Code: %d", resp.StatusCode)
	log.Debugf("Status: %s", resp.Status)
	log.Debugf("Headers: %v", resp.Header)

	// Log response body for debugging
	if resp.Body != nil {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err == nil {
			log.Debugf("Response Body: %s", string(bodyBytes))
			// Restore the body
			resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		} else {
			log.Debugf("Response Body: Error reading body - %v", err)
		}
	} else {
		log.Debugf("Response Body: None")
	}

	log.Debugf("=" + strings.Repeat("=", 49))

	if c.cfg.Debug {
		dump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			return resp, err
		}
		log.Debugf("\n%s\n", string(dump))
	}
	return resp, err
}

// prepareRequest build the request
func (c *APIClient) prepareRequest(
	ctx context.Context,
	path string, method string,
	postBody interface{},
	headerParams map[string]string,
	queryParams url.Values) (localVarRequest *http.Request, err error) {

	var body *bytes.Buffer

	// Detect postBody type and post.
	if postBody != nil {
		contentType := headerParams["Content-Type"]
		if contentType == "" {
			contentType = detectContentType(postBody)
			headerParams["Content-Type"] = contentType
		}

		body, err = setBody(postBody, contentType)
		if err != nil {
			return nil, err
		}
	}

	// Setup path and query parameters
	url, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	// Override request host, if applicable
	if c.cfg.Host != "" {
		url.Host = c.cfg.Host
	}

	// Override request scheme, if applicable
	if c.cfg.Scheme != "" {
		url.Scheme = c.cfg.Scheme
	}

	// Adding Query Param
	query := url.Query()
	for k, v := range queryParams {
		for _, iv := range v {
			query.Add(k, iv)
		}
	}
	url.RawQuery = query.Encode()

	// Generate a new request
	if body != nil {
		localVarRequest, err = http.NewRequest(method, url.String(), body)
	} else {
		localVarRequest, err = http.NewRequest(method, url.String(), nil)
	}
	if err != nil {
		return nil, err
	}

	// add header parameters, if any
	if len(headerParams) > 0 {
		headers := http.Header{}
		for h, v := range headerParams {
			headers[h] = []string{v}
		}
		localVarRequest.Header = headers
	}

	// Add the user agent to the request.
	localVarRequest.Header.Add("User-Agent", c.cfg.UserAgent)

	if ctx != nil {
		// add context to the request
		localVarRequest = localVarRequest.WithContext(ctx)

		// Walk through any authentication.

		// LoginToken Authentication
		if auth, ok := ctx.Value(ContextLoginToken).(string); ok {
			localVarRequest.Header.Add("Authorization", "Bearer "+auth)
		}
	}

	for header, value := range c.cfg.DefaultHeader {
		localVarRequest.Header.Add(header, value)
	}
	return localVarRequest, nil
}

func (c *APIClient) decode(v interface{}, b []byte, contentType string) (err error) {
	if len(b) == 0 {
		return nil
	}
	if s, ok := v.(*string); ok {
		*s = string(b)
		return nil
	}
	if JsonCheck.MatchString(contentType) {
		if err = json.Unmarshal(b, v); err != nil {
			return err
		}
		return nil
	}
	return errors.New("undefined response type")
}

// Set request body from an interface{}
func setBody(body interface{}, contentType string) (bodyBuf *bytes.Buffer, err error) {
	if bodyBuf == nil {
		bodyBuf = &bytes.Buffer{}
	}

	if reader, ok := body.(io.Reader); ok {
		_, err = bodyBuf.ReadFrom(reader)
	} else if b, ok := body.([]byte); ok {
		_, err = bodyBuf.Write(b)
	} else if s, ok := body.(string); ok {
		_, err = bodyBuf.WriteString(s)
	} else if s, ok := body.(*string); ok {
		_, err = bodyBuf.WriteString(*s)
	} else if JsonCheck.MatchString(contentType) {
		err = json.NewEncoder(bodyBuf).Encode(body)
	}

	if err != nil {
		return nil, err
	}

	if bodyBuf.Len() == 0 {
		err = fmt.Errorf("invalid body type %s\n", contentType)
		return nil, err
	}
	return bodyBuf, nil
}

// detectContentType method is used to figure out `Request.Body` content type for request header
func detectContentType(body interface{}) string {
	contentType := "text/plain; charset=utf-8"
	switch body.(type) {
	case map[string]interface{}, []interface{}, interface{}:
		contentType = "application/json; charset=utf-8"
	case string:
		contentType = "text/plain; charset=utf-8"
	default:
		if b, ok := body.([]byte); ok {
			contentType = http.DetectContentType(b)
		} else {
			contentType = "application/json; charset=utf-8"
		}
	}

	return contentType
}

// GenericOpenAPIError Provides access to the body, error and model on returned errors.
type GenericOpenAPIError struct {
	body  []byte
	error string
	model interface{}
}

// Error returns non-empty string if there was an error.
func (e GenericOpenAPIError) Error() string {
	return e.error
}

// Body returns the raw bytes of the response
func (e GenericOpenAPIError) Body() []byte {
	return e.body
}

// Model returns the unpacked model of the error
func (e GenericOpenAPIError) Model() interface{} {
	return e.model
}
