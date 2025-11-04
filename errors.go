package grist

import (
	"fmt"
)

// APIError represents a non-2xx HTTP response.
type APIError struct {
	StatusCode int
	Status     string
	Body       string
}

func (e *APIError) Error() string {
	if e.Body == "" {
		return fmt.Sprintf("api error: %s", e.Status)
	}
	return fmt.Sprintf("api error: %s - %s", e.Status, e.Body)
}
