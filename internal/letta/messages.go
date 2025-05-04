package letta

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Role string

const (
	RoleUser      Role = "user"
	RoleSystem    Role = "system"
	RoleAssistant Role = "assistant"

	RoleTool     Role = "tool"
	RoleFunction Role = "function"
)

type MessageContentType string

const (
	MessageContentTypeText              MessageContentType = "text"
	MessageContentTypeToolReturb        MessageContentType = "tool_return"
	MessageContentTypeReasoning         MessageContentType = "reasoning"
	MessageContentTypeOmiitedReasoning  MessageContentType = "omitted_reasoning"
	MessageContentTypeRedactedReasoning MessageContentType = "redacted_reasoning"
)

// Message Type Enum
type MessageType string

const (
	MessageTypeSystem          MessageType = "system_message"
	MessageTypeUser            MessageType = "user_message"
	MessageTypeReasoning       MessageType = "reasoning_message"
	MessageTypeHiddenReasoning MessageType = "hidden_reasoning_message"
	MessageTypeToolCall        MessageType = "tool_call_message"
	MessageTypeToolReturn      MessageType = "tool_return_message"
	MessageTypeAssistant       MessageType = "assistant_message"
)

// Reasoning Source Enum
type ReasoningSource string

const (
	ReasonerModel    ReasoningSource = "reasoner_model"
	NonReasonerModel ReasoningSource = "non_reasoner_model"
)

// Tool Return Status Enum
type ToolReturnStatus string

const (
	ToolReturnSuccess ToolReturnStatus = "success"
	ToolReturnError   ToolReturnStatus = "error"
)

// Hidden Reasoning State Enum
type HiddenReasoningState string

const (
	HiddenReasoningStateRedacted HiddenReasoningState = "redacted"
	HiddenReasoningStateOmitted  HiddenReasoningState = "omitted"
)

// MessageRequest represents the structure for sending messages to an agent
type MessageRequest struct {
	Messages                   []MessageContent `json:"messages"`
	UseAssistantMessage        *bool            `json:"use_assistant_message,omitempty"`
	AssistantMessageToolName   *string          `json:"assistant_message_tool_name,omitempty"`
	AssistantMessageToolKwargs *string          `json:"assistant_message_tool_kwargs,omitempty"`
}

// MessageContent represents the content of a message
type MessageContent struct {
	Role     Role    `json:"role"`
	Content  any     `json:"content"`
	Name     *string `json:"name,omitempty"`
	OtID     *string `json:"otid,omitempty"`
	SenderID *string `json:"sender_id,omitempty"`
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

// ReasoningContent represents reasoning message content
type ReasoningContent struct {
	Type      string  `json:"type"`
	IsNative  bool    `json:"is_native"`
	Reasoning string  `json:"reasoning"`
	Signature *string `json:"signature,omitempty"`
}

type RedactedReasoningContent struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type OmittedReasoningContent struct {
	Type string `json:"type"`
}

// Usage represents the usage statistics
type Usage struct {
	MessageType      string          `json:"message_type"`
	CompletionTokens int             `json:"completion_tokens"`
	PromptTokens     int             `json:"prompt_tokens"`
	TotalTokens      int             `json:"total_tokens"`
	StepCount        int             `json:"step_count"`
	StepsMessages    [][]StepMessage `json:"steps_messages"`
	RunIDs           []string        `json:"run_ids"`
}

// StepMessage represents a message within steps
type StepMessage struct {
	Role            Role      `json:"role,omitempty"`
	CreatedByID     string    `json:"created_by_id,omitempty"`
	LastUpdatedByID string    `json:"last_updated_by_id,omitempty"`
	CreatedAt       *FlexTime `json:"created_at,omitempty"`
	UpdatedAt       *FlexTime `json:"updated_at,omitempty"`
	ID              string    `json:"id,omitempty"`
	AgentID         string    `json:"agent_id,omitempty"`
	Model           string    `json:"model,omitempty"`
	Content         []any     `json:"content,omitempty"`
	Name            string    `json:"name,omitempty"`
	ToolCalls       []struct {
		ID       string `json:"id"`
		Function struct {
			Arguments string `json:"arguments"`
			Name      string `json:"name"`
		} `json:"function"`
		Type string `json:"type"`
	} `json:"tool_calls,omitempty"`
	ToolCallID  string `json:"tool_call_id,omitempty"`
	StepID      string `json:"step_id,omitempty"`
	OtID        string `json:"otid,omitempty"`
	ToolReturns []struct {
		Status string   `json:"status,omitempty"`
		Stdout []string `json:"stdout,omitempty"`
		Stderr []string `json:"stderr,omitempty"`
	} `json:"tool_returns,omitempty"`
	GroupID  string `json:"group_id,omitempty"`
	SenderID string `json:"sender_id,omitempty"`
}

type MessageResponse struct {
	Messages Messages `json:"messages"`
	Usage    Usage    `json:"usage"`
}

// Custom JSON Unmarshaler for Messages
type Messages []Message

// Message interface for polymorphic handling
type Message interface {
	GetMessageType() MessageType
	GetContent() string
}

// Base Message Struct
type BaseMessage struct {
	ID          string      `json:"id"`
	Date        FlexTime    `json:"date"`
	MessageType MessageType `json:"message_type"`
	Name        *string     `json:"name,omitempty"`
	OtID        *string     `json:"otid,omitempty"`
	SenderID    *string     `json:"sender_id,omitempty"`
	StepID      *string     `json:"step_id,omitempty"`
}

// Message Type Specific Structs (with previous definitions)
type SystemMessage struct {
	BaseMessage
	Content string `json:"content"`
}

func (m *SystemMessage) GetContent() string {
	return m.Content
}

type UserMessageTextContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type UserMessage struct {
	BaseMessage
	Content any `json:"content"`
}

func (m *UserMessage) GetContent() string {
	switch m.Content.(type) {
	case string:
		return m.Content.(string)
	case []map[string]any:
		contentParts := m.Content.([]map[string]any)
		textParts := make([]string, 0, len(contentParts))
		for _, part := range contentParts {
			textParts = append(textParts, part["text"].(string))
		}
		return strings.Join(textParts, "\n\n")
	default:
		return ""
	}
}

type ReasoningMessage struct {
	BaseMessage
	Reasoning string           `json:"reasoning"`
	Source    *ReasoningSource `json:"source,omitempty"`
	Signature *string          `json:"signature,omitempty"`
}

func (m *ReasoningMessage) GetContent() string {
	return m.Reasoning
}

type HiddenReasoningMessage struct {
	BaseMessage
	State           HiddenReasoningState `json:"state"`
	HiddenReasoning *string              `json:"hidden_reasoning,omitempty"`
}

func (m *HiddenReasoningMessage) GetContent() string {
	return *m.HiddenReasoning
}

type ToolCallMessage struct {
	BaseMessage
	ToolCall struct {
		Name       string `json:"name"`
		Arguments  string `json:"arguments"`
		ToolCallID string `json:"tool_call_id"`
	} `json:"tool_call"`
}

func (m *ToolCallMessage) GetContent() string {
	return fmt.Sprintf("%s:\n%s", m.ToolCall.Name, m.ToolCall.Arguments)
}

type ToolReturnMessage struct {
	BaseMessage
	ToolReturn string           `json:"tool_return"`
	Status     ToolReturnStatus `json:"status"`
	ToolCallID string           `json:"tool_call_id"`
	Stdout     []string         `json:"stdout,omitempty"`
	Stderr     []string         `json:"stderr,omitempty"`
}

func (m *ToolReturnMessage) GetContent() string {
	if m.ToolReturn != "" && m.ToolReturn != "None" {
		return m.ToolReturn
	}
	if len(m.Stderr) > 0 {
		return strings.Join(m.Stderr, "\n")
	}
	if len(m.Stdout) > 0 {
		return strings.Join(m.Stdout, "\n")
	}

	return ""
}

type AssistantMessage struct {
	BaseMessage
	Content any `json:"content"`
}

func (m *AssistantMessage) GetContent() string {
	switch m.Content.(type) {
	case string:
		return m.Content.(string)
	case []map[string]any:
		contentParts := m.Content.([]map[string]any)
		textParts := make([]string, 0, len(contentParts))
		for _, part := range contentParts {
			textParts = append(textParts, part["text"].(string))
		}
		return strings.Join(textParts, "\n\n")
	default:
		return ""
	}
}

func (b *BaseMessage) GetMessageType() MessageType {
	return b.MessageType
}

// Implement a custom UnmarshalJSON for the interface
func UnmarshalMessage(data []byte) (Message, error) {
	// First, extract the message type
	var baseMsg struct {
		MessageType MessageType `json:"message_type"`
	}

	if err := json.Unmarshal(data, &baseMsg); err != nil {
		return nil, err
	}

	// Based on the message type, unmarshal into the specific message struct
	switch baseMsg.MessageType {
	case MessageTypeSystem:
		var msg SystemMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			return nil, err
		}
		return &msg, nil

	case MessageTypeUser:
		var msg UserMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			return nil, err
		}
		return &msg, nil

	case MessageTypeReasoning:
		var msg ReasoningMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			return nil, err
		}
		return &msg, nil

	case MessageTypeHiddenReasoning:
		var msg HiddenReasoningMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			return nil, err
		}
		return &msg, nil

	case MessageTypeToolCall:
		var msg ToolCallMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			return nil, err
		}
		return &msg, nil

	case MessageTypeToolReturn:
		var msg ToolReturnMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			return nil, err
		}
		return &msg, nil

	case MessageTypeAssistant:
		var msg AssistantMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			return nil, err
		}
		return &msg, nil

	default:
		return nil, fmt.Errorf("unknown message type: %s", baseMsg.MessageType)
	}
}

// Custom JSON Unmarshaler for slice of Messages
func (m *Messages) UnmarshalJSON(data []byte) error {
	var rawMessages []json.RawMessage
	if err := json.Unmarshal(data, &rawMessages); err != nil {
		return err
	}

	messages := make([]Message, 0, len(rawMessages))
	for _, rawMsg := range rawMessages {
		msg, err := UnmarshalMessage(rawMsg)
		if err != nil {
			return err
		}
		messages = append(messages, msg)
	}

	*m = messages
	return nil
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
