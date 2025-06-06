// Package httpclient provides a simplified HTTP client for making JSON API requests.
//
// This package offers both a default client for simple use cases and a configurable
// Client type for more advanced scenarios. All HTTP methods support context cancellation
// and customizable options including headers and status code handling.
//
// Basic usage with the default client:
//
//	var result MyStruct
//	err := httpclient.Get(ctx, "https://api.example.com/data", &result)
//
// Using custom headers:
//
//	err := httpclient.Post(ctx, url, requestBody, &result,
//	    httpclient.WithHeader("Authorization", "Bearer token"))
//
// Handling non-200 status codes:
//
//	var status int
//	err := httpclient.Get(ctx, url, &result, httpclient.WithStatus(&status))
//	if status == 404 {
//	    // Handle not found
//	}
//
// Custom client configuration:
//
//	client := &httpclient.Client{
//	    Client: &http.Client{Timeout: 30 * time.Second},
//	    MarshalFunc: customMarshal,
//	    UnmarshalFunc: customUnmarshal,
//	}
//	err := client.Get(ctx, url, &result)
package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Client wraps http.Client with JSON utilities
type Client struct {
	// Client is the underlying HTTP client used for requests, defaults to http.DefaultClient
	Client *http.Client
	// MarshalFunc is used to marshal Go values into JSON, defaults to json.Marshal
	MarshalFunc func(v any) ([]byte, error)
	// UnmarshalFunc is used to unmarshal JSON data into a Go value, defaults to json.Unmarshal
	UnmarshalFunc func(data []byte, v any) error
}

// Get performs a GET request and unmarshals JSON response
func (c *Client) Get(ctx context.Context, url string, result interface{}, opts ...Option) error {
	options := buildOptions(opts...)

	req, err := c.buildRequest(ctx, http.MethodGet, url, nil, options)
	if err != nil {
		return fmt.Errorf("failed to create GET request: %w", err)
	}

	client := c.getClient(options)
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make GET request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	return c.parseResponse(resp, result, options)
}

// Post performs a POST request with JSON body and unmarshals JSON response
func (c *Client) Post(ctx context.Context, url string, body interface{}, result interface{}, opts ...Option) error {
	options := buildOptions(opts...)

	req, err := c.buildRequest(ctx, http.MethodPost, url, body, options)
	if err != nil {
		return fmt.Errorf("failed to create POST request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := c.getClient(options)
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make POST request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	return c.parseResponse(resp, result, options)
}

// Patch performs a PATCH request with JSON body and unmarshals JSON response
func (c *Client) Patch(ctx context.Context, url string, body interface{}, result interface{}, opts ...Option) error {
	options := buildOptions(opts...)

	req, err := c.buildRequest(ctx, http.MethodPatch, url, body, options)
	if err != nil {
		return fmt.Errorf("failed to create PATCH request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := c.getClient(options)
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make PATCH request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	return c.parseResponse(resp, result, options)
}

// Delete performs a DELETE request and unmarshals JSON response
func (c *Client) Delete(ctx context.Context, url string, result interface{}, opts ...Option) error {
	options := buildOptions(opts...)

	req, err := c.buildRequest(ctx, http.MethodDelete, url, nil, options)
	if err != nil {
		return fmt.Errorf("failed to create DELETE request: %w", err)
	}

	client := c.getClient(options)
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make DELETE request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	return c.parseResponse(resp, result, options)
}

// parseResponse reads and unmarshals JSON response
func (c *Client) parseResponse(resp *http.Response, result interface{}, options *Options) error {
	// Set status code if pointer provided
	if options.Status != nil {
		*options.Status = resp.StatusCode
	}

	// Return error for non-OK status codes unless Status pointer is provided
	if options.Status == nil && resp.StatusCode >= 400 {
		return fmt.Errorf("HTTP error: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if result != nil {
		if err := c.unmarshal(body, result); err != nil {
			// If status is being captured, don't fail on unmarshal errors for non-OK responses
			if options.Status != nil && resp.StatusCode >= 400 {
				// For non-OK responses with status capture, ignore unmarshal errors
				return nil
			}
			return fmt.Errorf("failed to unmarshal JSON response: %w", err)
		}
	}

	return nil
}

func (c *Client) getClient(options *Options) *http.Client {
	if options.Client != nil {
		return options.Client
	}
	if c.Client == nil {
		return http.DefaultClient
	}
	return c.Client
}

// buildRequest creates an HTTP request with the given method, URL, and body
func (c *Client) buildRequest(ctx context.Context, method, url string, body interface{}, options *Options) (*http.Request, error) {
	var bodyBytes []byte
	if body != nil {
		var err error
		bodyBytes, err = c.marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Apply headers from options
	for key, value := range options.Headers {
		req.Header.Set(key, value)
	}

	return req, nil
}

func (c *Client) marshal(v any) ([]byte, error) {
	if c.MarshalFunc != nil {
		return c.MarshalFunc(v)
	}
	return json.Marshal(v)
}

func (c *Client) unmarshal(data []byte, v any) error {
	if c.UnmarshalFunc != nil {
		return c.UnmarshalFunc(data, v)
	}
	return json.Unmarshal(data, v)
}
