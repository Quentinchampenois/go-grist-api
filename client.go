package grist

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"
)

const apiPath = "/api"

// Client is the Grist API client
type Client struct {
	Endpoint   string
	ApiKey     string
	HTTPClient *http.Client
	Context    context.Context
}

// ApiEndpoint returns the API base endpoint for the Grist instance
func (c *Client) ApiEndpoint() string {
	return c.Endpoint + apiPath
}

func NewGristClient(ctx context.Context, baseUrl, apiKey string) (*Client, error) {
	if baseUrl == "" {
		return nil, errors.New("endpoint cannot be empty")
	}
	if apiKey == "" {
		return nil, errors.New("API key cannot be empty")
	}

	// Strip last '/' if exists in baseUrl because ApiPath already contains leading '/'
	baseUrl = strings.TrimRight(baseUrl, "/")

	return &Client{
		Endpoint: baseUrl,
		ApiKey:   apiKey,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns: 10,
			},
		},
		Context: ctx,
	}, nil
}
