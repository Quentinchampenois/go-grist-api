package grist

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewGristClient(t *testing.T) {
	t.Run("With valid input", func(t *testing.T) {
		ctx := context.Background()
		baseURL := "https://getgrist.com"
		apiKey := "valid-key"

		client, err := NewGristClient(ctx, baseURL, apiKey)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if client.Endpoint != baseURL {
			t.Errorf("Expected Endpoint to be %s, got %s", baseURL, client.Endpoint)
		}

		if client.ApiKey != apiKey {
			t.Errorf("Expected ApiKey to be %s, got %s", apiKey, client.ApiKey)
		}

		if client.Context != ctx {
			t.Error("Context was not properly assigned")
		}

		// Check HTTPClient settings
		if client.HTTPClient.Timeout != 10*time.Second {
			t.Errorf("Expected HTTPClient timeout to be 10s, got %v", client.HTTPClient.Timeout)
		}

		transport, ok := client.HTTPClient.Transport.(*http.Transport)
		if !ok {
			t.Fatal("Expected HTTPClient.Transport to be *http.Transport")
		}
		if transport.MaxIdleConns != 10 {
			t.Errorf("Expected MaxIdleConns to be 10, got %d", transport.MaxIdleConns)
		}
	})
	t.Run("With trailing slash in endpoint removes it", func(t *testing.T) {
		ctx := context.Background()
		endpoint := "https://getgrist.com/"
		expected := "https://getgrist.com"
		apiKey := "valid-key"

		client, err := NewGristClient(ctx, endpoint, apiKey)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if client.Endpoint != expected {
			t.Errorf("Expected Endpoint to be %s, got %s", expected, client.Endpoint)
		}
	})
	t.Run("With empty endpoint returns error", func(t *testing.T) {
		ctx := context.Background()
		endpoint := ""
		apiKey := "valid-key"

		_, err := NewGristClient(ctx, endpoint, apiKey)
		if err == nil {
			t.Fatalf("Expected error for empty endpoint")
		}

		assert.Equal(t, "endpoint cannot be empty", err.Error())
	})
	t.Run("With empty api key returns error", func(t *testing.T) {
		ctx := context.Background()
		endpoint := "https://getgrist.com"
		apiKey := ""

		_, err := NewGristClient(ctx, endpoint, apiKey)
		if err == nil {
			t.Fatalf("Expected error for empty apiKey")
		}

		assert.Equal(t, "API key cannot be empty", err.Error())
	})
}

// TestClient_ApiEndpoint checks that the full API endpoint is correctly formed.
func TestClient_ApiEndpoint(t *testing.T) {
	ctx := context.Background()
	baseURL := "https://getgrist.com"
	apiKey := "valid-key"

	client, err := NewGristClient(ctx, baseURL, apiKey)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	expected := baseURL + "/api"
	got := client.ApiEndpoint()

	if got != expected {
		t.Errorf("ApiEndpoint() = %s, expected %s", got, expected)
	}
}
