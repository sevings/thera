package letta

import (
	"context"
	"fmt"
	"net/http"
)

// MessageRequest represents the structure for sending messages to an agent
type MessageRequest struct {
	Messages                   []MessageContent `json:"messages"`
	UseAssistantMessage        *bool            `json:"use_assistant_message,omitempty"`
	AssistantMessageToolName   *string          `json:"assistant_message_tool_name,omitempty"`
	AssistantMessageToolKwargs *string          `json:"assistant_message_tool_kwargs,omitempty"`
	GroupID                    *string          `json:"group_id,omitempty"`
	SenderID                   *string          `json:"sender_id,omitempty"`
}

type Role string

const (
	RoleUser      Role = "user"
	RoleSystem    Role = "system"
	RoleAssistant Role = "assistant"
)

// MessageContent represents the content of a message
type MessageContent struct {
	Role     Role    `json:"role"`
	Content  any     `json:"content"`
	Name     *string `json:"name,omitempty"`
	OtID     *string `json:"otid,omitempty"`
	SenderID *string `json:"sender_id,omitempty"`
}

// MessageResponse represents the response from sending a message
type MessageResponse struct {
	Messages []map[string]any `json:"messages"`
}

// TextContent represents text message content
type TextContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// ToolCallContent represents tool call message content
type ToolCallContent struct {
	Type  string         `json:"type"`
	ID    string         `json:"id"`
	Name  string         `json:"name"`
	Input map[string]any `json:"input"`
}

// ToolReturnContent represents tool return message content
type ToolReturnContent struct {
	Type       string `json:"type"`
	ToolCallID string `json:"id"`
	Name       string `json:"name"`
	Content    string `json:"content"`
	IsError    bool   `json:"is_error"`
}

// ReasoningContent represents reasoning message content
type ReasoningContent struct {
	Type      string  `json:"type"`
	IsNative  bool    `json:"is_native"`
	Reasoning string  `json:"reasoning"`
	Signature *string `json:"signature,omitempty"`
}

// SendMessage sends a message to a specific agent and returns the agent's response
func (c *APIClient) SendMessage(ctx context.Context, agentID string, req MessageRequest) (*MessageResponse, error) {
	if agentID == "" {
		return nil, fmt.Errorf("agent ID is required")
	}

	if len(req.Messages) == 0 {
		return nil, fmt.Errorf("at least one message is required")
	}

	apiReq := APIRequest{
		Method:   http.MethodPost,
		Endpoint: fmt.Sprintf("/v1/agents/%s/messages", agentID),
		Body:     req,
	}

	var response MessageResponse
	err := c.sendRequest(ctx, apiReq, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
