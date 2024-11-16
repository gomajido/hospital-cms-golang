package http_wrapper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"
)

type IHTTPClientWrapper interface {
	Post(url string, request interface{}, headers http.Header) (*http.Response, error)
	Get(url string, headers http.Header) (*http.Response, error)
	PostWithContext(ctx context.Context, url string, request interface{}, headers http.Header) (*http.Response, error)
	PostMultipart(url string, body *bytes.Buffer, headers http.Header) (*http.Response, error)
	PostMultipartWithContext(ctx context.Context, url string, body *bytes.Buffer, headers http.Header) (*http.Response, error)
}

// HTTPClientWrapper is a struct that wraps an http.Client and adds custom functionality.
type HTTPClientWrapper struct {
	client *http.Client
}

// NewHTTPClientWrapper creates a new HTTPClientWrapper with a specified timeout.
func NewHTTPClientWrapper(timeout time.Duration) *HTTPClientWrapper {
	return &HTTPClientWrapper{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// Get sends a GET request using the wrapped client, adds custom headers, and returns the *http.Response and any error.
func (c *HTTPClientWrapper) Get(url string, headers http.Header) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header = headers

	resp, err := c.client.Do(req)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return nil, fmt.Errorf("%s", HTTP_ERR_CONTEXT_DEADLINE_EXCEEDED)
		}
		return nil, err
	}

	return resp, nil
}

// Post sends a POST request using the wrapped client, adds custom headers, and returns the *http.Response and any error.
func (c *HTTPClientWrapper) Post(url string, request interface{}, headers http.Header) (*http.Response, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header = headers

	resp, err := c.client.Do(req)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			// Handle network timeout error
			return nil, fmt.Errorf("%s", HTTP_ERR_CONTEXT_DEADLINE_EXCEEDED)
		}
		return nil, err
	}

	return resp, nil
}

// Post sends a POST request using the wrapped client, adds custom headers, and returns the *http.Response and any error.
func (c *HTTPClientWrapper) PostWithContext(ctx context.Context, url string, request interface{}, headers http.Header) (*http.Response, error) {
	body, _ := json.Marshal(request)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header = headers

	resp, err := c.client.Do(req)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			// Handle network timeout error
			return nil, fmt.Errorf("%s", HTTP_ERR_CONTEXT_DEADLINE_EXCEEDED)
		}
		return nil, err
	}

	return resp, nil
}

func (c *HTTPClientWrapper) PostMultipart(url string, body *bytes.Buffer, headers http.Header) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	req.Header = headers
	req.Header.Set("Content-Type", "multipart/form-data")

	resp, err := c.client.Do(req)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			// Handle network timeout error
			return nil, fmt.Errorf("%s", HTTP_ERR_CONTEXT_DEADLINE_EXCEEDED)
		}
		return nil, err
	}

	return resp, nil
}

func (c *HTTPClientWrapper) PostMultipartWithContext(ctx context.Context, url string, body *bytes.Buffer, headers http.Header) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return nil, err
	}

	req.Header = headers

	resp, err := c.client.Do(req)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			// Handle network timeout error
			return nil, fmt.Errorf("%s", HTTP_ERR_CONTEXT_DEADLINE_EXCEEDED)
		}
		return nil, err
	}

	return resp, nil
}
