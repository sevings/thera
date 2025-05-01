package letta

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Config represents the configuration for the Letta client
type Config struct {
	BaseURL string `koanf:"base_url"`
	Token   string
	Timeout time.Duration
}

// APIClient provides common methods for API interactions
type APIClient struct {
	config     Config
	httpClient *http.Client
	logger     *zap.Logger
}

// APIRequest represents a generic API request
type APIRequest struct {
	Method      string
	Endpoint    string
	Body        any
	QueryParams map[string]string
}

// NewClient creates a new API client
func NewClient(config Config) (*APIClient, error) {
	if config.Timeout == 0 {
		config.Timeout = 30
	}

	httpClient := &http.Client{
		Timeout: config.Timeout * time.Second,
	}

	return &APIClient{
		config:     config,
		httpClient: httpClient,
		logger:     zap.L().Named("letta"),
	}, nil
}

// sendRequest sends an HTTP request and handles the response
func (c *APIClient) sendRequest(ctx context.Context, req APIRequest, respHandler any) error {
	logger := c.logger.With(
		zap.String("method", req.Method),
		zap.String("endpoint", req.Endpoint),
	)

	start := time.Now()
	logger.Debug("Initiating API request")

	url := fmt.Sprintf("%s%s", c.config.BaseURL, req.Endpoint)

	var body io.Reader
	if req.Body != nil {
		jsonPayload, err := json.Marshal(req.Body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		body = bytes.NewBuffer(jsonPayload)
	}

	httpReq, err := http.NewRequestWithContext(ctx, req.Method, url, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	c.addRequestHeaders(httpReq, req)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		logger.Error("Request failed",
			zap.Error(err),
			zap.Duration("elapsed", time.Since(start)),
		)
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.handleErrorResponse(resp, logger, start)
	}

	return c.parseResponse(resp, respHandler, logger, start)
}

// addRequestHeaders adds standard headers to the request
func (c *APIClient) addRequestHeaders(req *http.Request, apiReq APIRequest) {
	if apiReq.Body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if c.config.Token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.Token))
	}

	if len(apiReq.QueryParams) > 0 {
		q := req.URL.Query()
		for k, v := range apiReq.QueryParams {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}
}

// handleErrorResponse processes and logs error responses
func (c *APIClient) handleErrorResponse(resp *http.Response, logger *zap.Logger, start time.Time) error {
	errorBody, _ := io.ReadAll(resp.Body)

	logger.Error("Unexpected status code",
		zap.Int("status_code", resp.StatusCode),
		zap.String("error_body", string(errorBody)),
		zap.Duration("elapsed", time.Since(start)),
	)

	return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
}

// parseResponse reads and parses the API response
func (c *APIClient) parseResponse(resp *http.Response, respHandler any, logger *zap.Logger, start time.Time) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Failed to read response body",
			zap.Error(err),
			zap.Duration("elapsed", time.Since(start)),
		)
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if respHandler == nil {
		return nil
	}

	if err := json.Unmarshal(body, respHandler); err != nil {
		logger.Error("Failed to parse response",
			zap.Error(err),
			zap.String("response_body", string(body)),
			zap.Duration("elapsed", time.Since(start)),
		)
		return fmt.Errorf("failed to parse response: %w", err)
	}

	logger.Info("Request completed successfully",
		zap.Duration("elapsed", time.Since(start)),
	)

	return nil
}
