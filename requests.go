package grist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type requestOption func(*http.Request)

// withJSONBody sets the request body to the JSON-encoded value of body
func withJSONBody(body any) (requestOption, error) {
	if body == nil {
		return func(*http.Request) {}, nil
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		return nil, err
	}
	return func(r *http.Request) {
		r.Body = io.NopCloser(&buf)
		r.ContentLength = int64(buf.Len())
		r.Header.Set("Content-Type", "application/json")
	}, nil
}

// withAuth sets the Authorization header to the given API key
func withAuth(apiKey string) requestOption {
	return func(r *http.Request) {
		r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	}
}

// doRequest performs an HTTP request with the given options
func (c *Client) DoRequest(method, endpoint string, opts ...requestOption) (*http.Response, error) {
	req, err := http.NewRequestWithContext(c.Context, method, endpoint, nil)
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		opt(req)
	}

	return c.HTTPClient.Do(req)
}

// GetRequest performs a GET request
func (c *Client) GetRequest(endpoint string, opts ...requestOption) (*http.Response, error) {
	return c.DoRequest(http.MethodGet, endpoint, opts...)
}

// PostRequest performs a POST request
func (c *Client) PostRequest(endpoint string, opts ...requestOption) (*http.Response, error) {
	return c.DoRequest(http.MethodPost, endpoint, opts...)
}

// PatchRequest performs a PATCH request
func (c *Client) PatchRequest(endpoint string, opts ...requestOption) (*http.Response, error) {
	return c.DoRequest(http.MethodPatch, endpoint, opts...)
}

// DeleteRequest performs a DELETE request
func (c *Client) DeleteRequest(endpoint string, opts ...requestOption) (*http.Response, error) {
	return c.DoRequest(http.MethodDelete, endpoint, opts...)
}

func buildURL(base, p string) string {
	return base + p
}

func handleStatus(resp *http.Response, okStatuses ...int) error {
	defer resp.Body.Close()
	ok := false
	for _, s := range okStatuses {
		if resp.StatusCode == s {
			ok = true
			break
		}
	}
	b, _ := io.ReadAll(resp.Body)
	if !ok {
		return &APIError{StatusCode: resp.StatusCode, Status: resp.Status, Body: string(bytes.TrimSpace(b))}
	}
	return nil
}

func handleJSONResponse[T any](resp *http.Response, out *T, okStatuses ...int) error {
	defer resp.Body.Close()

	// status check
	ok := false
	for _, s := range okStatuses {
		if resp.StatusCode == s {
			ok = true
			break
		}
	}
	if !ok {
		b, _ := io.ReadAll(resp.Body)
		return &APIError{StatusCode: resp.StatusCode, Status: resp.Status, Body: string(bytes.TrimSpace(b))}
	}

	dec := json.NewDecoder(resp.Body)
	err := dec.Decode(out)
	if err != nil {
		return fmt.Errorf("handleJSONResponse: decode json: %w", err)
	}
	return nil
}

func handleRawResponse(resp *http.Response, okStatuses ...int) ([]byte, error) {
	defer resp.Body.Close()

	ok := false
	for _, s := range okStatuses {
		if resp.StatusCode == s {
			ok = true
			break
		}
	}
	b, _ := io.ReadAll(resp.Body)
	if !ok {
		return nil, &APIError{StatusCode: resp.StatusCode, Status: resp.Status, Body: string(bytes.TrimSpace(b))}
	}

	b = bytes.TrimSpace(b)
	b = bytes.TrimSuffix(b, []byte("\""))
	b = bytes.TrimPrefix(b, []byte("\""))
	return b, nil
}
