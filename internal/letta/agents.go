package letta

import (
	"context"
	"fmt"
	"net/http"
)

// CreateAgentRequest represents the request payload for creating an agent
type CreateAgentRequest struct {
	Name                   string           `json:"name,omitempty"`
	MemoryBlocks           []MemoryBlock    `json:"memory_blocks,omitempty"`
	Tools                  []string         `json:"tools,omitempty"`
	ToolIDs                []string         `json:"tool_ids,omitempty"`
	SourceIDs              []string         `json:"source_ids,omitempty"`
	BlockIDs               []string         `json:"block_ids,omitempty"`
	ToolRules              []any            `json:"tool_rules,omitempty"`
	Tags                   []string         `json:"tags,omitempty"`
	System                 string           `json:"system,omitempty"`
	AgentType              string           `json:"agent_type,omitempty"`
	LLMConfig              *LLMConfig       `json:"llm_config,omitempty"`
	EmbeddingConfig        *EmbeddingConfig `json:"embedding_config,omitempty"`
	InitialMessageSequence []InitialMessage `json:"initial_message_sequence,omitempty"`
	Description            string           `json:"description,omitempty"`
	Metadata               map[string]any   `json:"metadata,omitempty"`
	ProjectID              string           `json:"project_id,omitempty"`
	IdentityIDs            []string         `json:"identity_ids,omitempty"`
}

// MemoryBlock represents a block in the agent's memory
type MemoryBlock struct {
	Value       string         `json:"value"`
	Label       string         `json:"label"`
	Limit       int            `json:"limit,omitempty"`
	Name        string         `json:"name,omitempty"`
	IsTemplate  bool           `json:"is_template,omitempty"`
	Description string         `json:"description,omitempty"`
	Metadata    map[string]any `json:"metadata,omitempty"`
}

// LLMConfig represents the LLM configuration for an agent
type LLMConfig struct {
	Model              string  `json:"model"`
	ModelEndpointType  string  `json:"model_endpoint_type"`
	ContextWindow      int     `json:"context_window"`
	ModelEndpoint      string  `json:"model_endpoint,omitempty"`
	ProviderName       string  `json:"provider_name,omitempty"`
	ModelWrapper       string  `json:"model_wrapper,omitempty"`
	Temperature        float64 `json:"temperature,omitempty"`
	MaxTokens          int     `json:"max_tokens,omitempty"`
	EnableReasoner     bool    `json:"enable_reasoner,omitempty"`
	MaxReasoningTokens int     `json:"max_reasoning_tokens,omitempty"`
}

// EmbeddingConfig represents the embedding configuration for an agent
type EmbeddingConfig struct {
	EmbeddingEndpointType string `json:"embedding_endpoint_type"`
	EmbeddingModel        string `json:"embedding_model"`
	EmbeddingDim          int    `json:"embedding_dim"`
	EmbeddingEndpoint     string `json:"embedding_endpoint,omitempty"`
	EmbeddingChunkSize    int    `json:"embedding_chunk_size,omitempty"`
}

// InitialMessage represents an initial message in the agent's memory
type InitialMessage struct {
	Role    string `json:"role"`
	Content []any  `json:"content"`
	Name    string `json:"name,omitempty"`
}

// Agent represents the response from creating an agent
type Agent struct {
	ID              string           `json:"id"`
	Name            string           `json:"name"`
	System          string           `json:"system"`
	AgentType       string           `json:"agent_type"`
	LLMConfig       *LLMConfig       `json:"llm_config"`
	EmbeddingConfig *EmbeddingConfig `json:"embedding_config"`
	MemoryBlocks    []MemoryBlock    `json:"memory"`
	Tools           []Tool           `json:"tools"`
	Sources         []Source         `json:"sources"`
	Tags            []string         `json:"tags"`
	Description     string           `json:"description,omitempty"`
	Metadata        map[string]any   `json:"metadata,omitempty"`
	ProjectID       string           `json:"project_id,omitempty"`
	IdentityIDs     []string         `json:"identity_ids,omitempty"`
}

// Tool represents a tool used by an agent
type Tool struct {
	ID          string         `json:"id,omitempty"`
	Name        string         `json:"name,omitempty"`
	Description string         `json:"description,omitempty"`
	Metadata    map[string]any `json:"metadata,omitempty"`
}

// Source represents a source used by an agent
type Source struct {
	ID          string         `json:"id,omitempty"`
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Metadata    map[string]any `json:"metadata,omitempty"`
}

// CreateAgent creates a new agent
func (c *APIClient) CreateAgent(ctx context.Context, request CreateAgentRequest) (*Agent, error) {
	apiReq := APIRequest{
		Method:   http.MethodPost,
		Endpoint: "/v1/agents/",
		Body:     request,
	}

	var response Agent
	err := c.sendRequest(ctx, apiReq, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// RetrieveAgent returns an existing agent
func (c *APIClient) RetrieveAgent(ctx context.Context, id string) (*Agent, error) {
	apiReq := APIRequest{
		Method:   http.MethodGet,
		Endpoint: fmt.Sprintf("/v1/agents/%s", id),
	}

	var response Agent
	err := c.sendRequest(ctx, apiReq, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// DeleteAgent removes an existing agent
func (c *APIClient) DeleteAgent(ctx context.Context, id string) error {
	apiReq := APIRequest{
		Method:   http.MethodDelete,
		Endpoint: fmt.Sprintf("/v1/agents/%s", id),
	}

	return c.sendRequest(ctx, apiReq, nil)
}

// CountAgents returns the total number of agents
func (c *APIClient) CountAgents(ctx context.Context) (int, error) {
	apiReq := APIRequest{
		Method:   http.MethodGet,
		Endpoint: "/v1/agents/count",
	}

	var response int
	err := c.sendRequest(ctx, apiReq, &response)
	if err != nil {
		return 0, err
	}

	return response, nil
}
