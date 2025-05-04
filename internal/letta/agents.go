package letta

import (
	"context"
	"fmt"
	"net/http"
)

type AgentType string

const (
	AgentTypeMemgpt         AgentType = "memgpt_agent"
	AgentTypeSleeptime      AgentType = "sleeptime_agent"
	AgentTypeSplitThread    AgentType = "split_thread_agent"
	AgentTypeVoiceConvo     AgentType = "voice_convo_agent"
	AgentTypeVoiceSleeptime AgentType = "voice_sleeptime_agent"
)

type CreateAgentRequest struct {
	AgentType                    AgentType          `json:"agent_type,omitempty"`
	BaseTemplateID               string             `json:"base_template_id,omitempty"`
	BlockIDs                     []string           `json:"block_ids,omitempty"`
	ContextWindowLimit           int                `json:"context_window_limit,omitempty"`
	Description                  string             `json:"description,omitempty"`
	Embedding                    string             `json:"embedding,omitempty"`
	EmbeddingChunkSize           int                `json:"embedding_chunk_size,omitempty"`
	EmbeddingConfig              *EmbeddingConfig   `json:"embedding_config,omitempty"`
	EnableReasoner               *bool              `json:"enable_reasoner,omitempty"`
	EnableSleeptime              *bool              `json:"enable_sleeptime,omitempty"`
	FromTemplate                 string             `json:"from_template,omitempty"`
	IdentityIDs                  []string           `json:"identity_ids,omitempty"`
	IncludeBaseToolRules         *bool              `json:"include_base_tool_rules,omitempty"`
	IncludeBaseTools             *bool              `json:"include_base_tools,omitempty"`
	IncludeMultiAgentTools       *bool              `json:"include_multi_agent_tools,omitempty"`
	InitialMessageSequence       []InitialMessage   `json:"initial_message_sequence,omitempty"`
	LLMConfig                    *LLMConfig         `json:"llm_config,omitempty"`
	MaxReasoningTokens           int                `json:"max_reasoning_tokens,omitempty"`
	MaxTokens                    int                `json:"max_tokens,omitempty"`
	MemoryBlocks                 []MemoryBlock      `json:"memory_blocks,omitempty"`
	MemoryVariables              map[string]*string `json:"memory_variables,omitempty"`
	MessageBufferAutoclear       *bool              `json:"message_buffer_autoclear,omitempty"`
	Metadata                     map[string]any     `json:"metadata,omitempty"`
	Model                        string             `json:"model,omitempty"`
	Name                         string             `json:"name,omitempty"`
	Project                      string             `json:"project,omitempty"`
	ProjectID                    string             `json:"project_id,omitempty"`
	ResponseFormat               *ResponseFormat    `json:"response_format,omitempty"`
	SourceIDs                    []string           `json:"source_ids,omitempty"`
	System                       string             `json:"system,omitempty"`
	Tags                         []string           `json:"tags,omitempty"`
	Template                     *bool              `json:"template,omitempty"`
	TemplateID                   string             `json:"template_id,omitempty"`
	ToolExecEnvironmentVariables map[string]*string `json:"tool_exec_environment_variables,omitempty"`
	ToolIDs                      []string           `json:"tool_ids,omitempty"`
	ToolRules                    []any              `json:"tool_rules,omitempty"`
	Tools                        []string           `json:"tools,omitempty"`
}

type AgentMemory struct {
	Blocks         []MemoryBlock `json:"blocks"`
	PromptTemplate string        `json:"prompt_template"`
}

type MemoryBlock struct {
	Description string         `json:"description,omitempty"`
	ID          string         `json:"id,omitempty"`
	IsTemplate  bool           `json:"is_template,omitempty"`
	Label       string         `json:"label,omitempty"`
	Limit       int            `json:"limit,omitempty"`
	Metadata    map[string]any `json:"metadata,omitempty"`
	Name        string         `json:"name,omitempty"`
	Value       string         `json:"value"`
	CreatedByID string         `json:"created_by_id,omitempty"`
	UpdatedByID string         `json:"last_updated_by_id,omitempty"`
}

type ResponseFormat struct {
	Type string `json:"type,omitempty"`
}

type LLMConfig struct {
	ContextWindow            int     `json:"context_window" koanf:"context_window"`
	EnableReasoner           *bool   `json:"enable_reasoner,omitempty" koanf:"enable_reasoner"`
	Handle                   string  `json:"handle,omitempty" koanf:"handle"`
	MaxReasoningTokens       int     `json:"max_reasoning_tokens,omitempty" koanf:"max_reasoning_tokens"`
	MaxTokens                int     `json:"max_tokens,omitempty" koanf:"max_tokens"`
	Model                    string  `json:"model" koanf:"model"`
	ModelEndpoint            string  `json:"model_endpoint,omitempty" koanf:"model_endpoint"`
	ModelEndpointType        string  `json:"model_endpoint_type" koanf:"model_endpoint_type"`
	ModelWrapper             string  `json:"model_wrapper,omitempty" koanf:"model_wrapper"`
	ProviderName             string  `json:"provider_name,omitempty" koanf:"provider_name"`
	PutInnerThoughtsInKwargs *bool   `json:"put_inner_thoughts_in_kwargs,omitempty"`
	ReasoningEffort          string  `json:"reasoning_effort,omitempty"`
	Temperature              float64 `json:"temperature,omitempty" koanf:"temperature"`
}

type EmbeddingConfig struct {
	AzureDeployment       string `json:"azure_deployment,omitempty"`
	AzureEndpoint         string `json:"azure_endpoint,omitempty"`
	AzureVersion          string `json:"azure_version,omitempty"`
	EmbeddingChunkSize    int    `json:"embedding_chunk_size,omitempty" koanf:"embedding_chunk_size"`
	EmbeddingDim          int    `json:"embedding_dim" koanf:"embedding_dim"`
	EmbeddingEndpoint     string `json:"embedding_endpoint,omitempty" koanf:"embedding_endpoint"`
	EmbeddingEndpointType string `json:"embedding_endpoint_type" koanf:"embedding_endpoint_type"`
	EmbeddingModel        string `json:"embedding_model" koanf:"embedding_model"`
	Handle                string `json:"handle,omitempty"`
}

type InitialMessage struct {
	Content []any  `json:"content"`
	Name    string `json:"name,omitempty"`
	Role    string `json:"role"`
}

type Agent struct {
	AgentType                    AgentType                `json:"agent_type"`
	BaseTemplateID               string                   `json:"base_template_id,omitempty"`
	CreatedAt                    *FlexTime                `json:"created_at,omitempty"`
	CreatedByID                  string                   `json:"created_by_id,omitempty"`
	Description                  string                   `json:"description,omitempty"`
	EmbeddingConfig              *EmbeddingConfig         `json:"embedding_config"`
	EnableSleeptime              bool                     `json:"enable_sleeptime,omitempty"`
	ID                           string                   `json:"id"`
	IdentityIDs                  []string                 `json:"identity_ids,omitempty"`
	LLMConfig                    *LLMConfig               `json:"llm_config"`
	LastUpdatedByID              string                   `json:"last_updated_by_id,omitempty"`
	Memory                       AgentMemory              `json:"memory"`
	MessageBufferAutoclear       bool                     `json:"message_buffer_autoclear,omitempty"`
	MessageIDs                   []string                 `json:"message_ids,omitempty"`
	Metadata                     map[string]any           `json:"metadata,omitempty"`
	MultiAgentGroup              *MultiAgentGroup         `json:"multi_agent_group,omitempty"`
	Name                         string                   `json:"name"`
	ProjectID                    string                   `json:"project_id,omitempty"`
	ResponseFormat               *ResponseFormat          `json:"response_format,omitempty"`
	Sources                      []Source                 `json:"sources"`
	System                       string                   `json:"system"`
	Tags                         []string                 `json:"tags"`
	TemplateID                   string                   `json:"template_id,omitempty"`
	ToolExecEnvironmentVariables []ToolExecEnvironmentVar `json:"tool_exec_environment_variables,omitempty"`
	ToolRules                    []ToolRule               `json:"tool_rules,omitempty"`
	Tools                        []Tool                   `json:"tools"`
	UpdatedAt                    *FlexTime                `json:"updated_at,omitempty"`
}

type Tool struct {
	ArgsJSONSchema  map[string]any `json:"args_json_schema,omitempty"`
	CreatedByID     string         `json:"created_by_id,omitempty"`
	Description     string         `json:"description,omitempty"`
	ID              string         `json:"id,omitempty"`
	JSONSchema      map[string]any `json:"json_schema,omitempty"`
	LastUpdatedByID string         `json:"last_updated_by_id,omitempty"`
	Metadata        map[string]any `json:"metadata_,omitempty"`
	Name            string         `json:"name,omitempty"`
	ReturnCharLimit int            `json:"return_char_limit,omitempty"`
	SourceCode      string         `json:"source_code,omitempty"`
	SourceType      string         `json:"source_type,omitempty"`
	Tags            []string       `json:"tags,omitempty"`
	ToolType        string         `json:"tool_type,omitempty"`
}

type Source struct {
	CreatedAt       *FlexTime        `json:"created_at,omitempty"`
	CreatedByID     string           `json:"created_by_id,omitempty"`
	Description     string           `json:"description,omitempty"`
	EmbeddingConfig *EmbeddingConfig `json:"embedding_config,omitempty"`
	ID              string           `json:"id,omitempty"`
	LastUpdatedByID string           `json:"last_updated_by_id,omitempty"`
	Metadata        map[string]any   `json:"metadata,omitempty"`
	Name            string           `json:"name"`
	UpdatedAt       *FlexTime        `json:"updated_at,omitempty"`
}

type ToolRule struct {
	ChildOutputMapping   map[string]any `json:"child_output_mapping,omitempty"`
	DefaultChild         string         `json:"default_child,omitempty"`
	RequireOutputMapping bool           `json:"require_output_mapping,omitempty"`
	ToolName             string         `json:"tool_name,omitempty"`
	Type                 string         `json:"type,omitempty"`
}

type ToolExecEnvironmentVar struct {
	AgentID         string    `json:"agent_id,omitempty"`
	CreatedAt       *FlexTime `json:"created_at,omitempty"`
	CreatedByID     string    `json:"created_by_id,omitempty"`
	Description     string    `json:"description,omitempty"`
	ID              string    `json:"id,omitempty"`
	Key             string    `json:"key"`
	LastUpdatedByID string    `json:"last_updated_by_id,omitempty"`
	UpdatedAt       *FlexTime `json:"updated_at,omitempty"`
	Value           string    `json:"value"`
}

type MultiAgentGroup struct {
	AgentIDs                []string `json:"agent_ids,omitempty"`
	Description             string   `json:"description,omitempty"`
	ID                      string   `json:"id,omitempty"`
	LastProcessedMessageID  string   `json:"last_processed_message_id,omitempty"`
	ManagerAgentID          string   `json:"manager_agent_id,omitempty"`
	ManagerType             string   `json:"manager_type,omitempty"`
	MaxMessageBufferLength  int      `json:"max_message_buffer_length,omitempty"`
	MaxTurns                int      `json:"max_turns,omitempty"`
	MinMessageBufferLength  int      `json:"min_message_buffer_length,omitempty"`
	SharedBlockIDs          []string `json:"shared_block_ids,omitempty"`
	SleeptimeAgentFrequency int      `json:"sleeptime_agent_frequency,omitempty"`
	TerminationToken        string   `json:"termination_token,omitempty"`
	TurnsCounter            int      `json:"turns_counter,omitempty"`
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

// RetrieveAgentMemory retrieves the memory state of a specific agent
func (c *APIClient) RetrieveAgentMemory(ctx context.Context, agentID string) (*AgentMemory, error) {
	apiReq := APIRequest{
		Method:   http.MethodGet,
		Endpoint: fmt.Sprintf("/v1/agents/%s/core-memory", agentID),
	}

	var response AgentMemory
	err := c.sendRequest(ctx, apiReq, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
