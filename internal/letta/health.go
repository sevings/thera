package letta

import (
	"context"
	"net/http"
)

// HealthResponse represents the response from the health check endpoint
type HealthResponse struct {
	Version string `json:"version"`
	Status  string `json:"status"`
}

// Health performs a health check
func (c *APIClient) Health(ctx context.Context) (*HealthResponse, error) {
	apiReq := APIRequest{
		Method:   http.MethodGet,
		Endpoint: "/v1/health/",
	}

	var response HealthResponse
	err := c.sendRequest(ctx, apiReq, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
