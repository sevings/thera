package thera

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"thera/internal/letta"

	"go.uber.org/zap"
)

type Thera struct {
	db  *DB
	api *letta.APIClient
	log *zap.SugaredLogger
	cfg Config
}

func NewThera(db *DB, cfg Config) (*Thera, error) {
	api, err := letta.NewClient(cfg.Letta)
	if err != nil {
		return nil, err
	}

	th := &Thera{
		db:  db,
		api: api,
		log: zap.L().Named("thera").Sugar(),
		cfg: cfg,
	}
	return th, nil
}

func (th *Thera) Start() error {
	ctx := context.Background()
	healthResp, err := th.api.Health(ctx)
	if err != nil {
		return err
	}

	th.log.Infow("Letta Health Check",
		"version", healthResp.Version,
		"status", healthResp.Status)

	return nil
}

func (th *Thera) CreateChat(ctx context.Context, userID int64, userName, userBio string) (*Chat, error) {
	chat, err := th.db.GetChatByUserID(userID)
	if err == nil {
		return chat, nil
	}

	if err != ErrNotFound {
		return nil, fmt.Errorf("error checking existing chat: %w", err)
	}

	identityReq := letta.CreateIdentityRequest{
		IdentifierKey: fmt.Sprintf("user_%d", userID),
		IdentityType:  letta.IdentityTypeUser,
	}

	if userName == "" {
		identityReq.Name = fmt.Sprintf("User %d", userID)
	} else {
		identityReq.Name = fmt.Sprintf("%s (%d)", userName, userID)
	}

	identity, err := th.api.CreateIdentity(ctx, identityReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create identity: %w", err)
	}

	// Load agent configuration from file
	var agentConfig letta.CreateAgentRequest
	if th.cfg.AgentPath != "" {
		configData, err := os.ReadFile(th.cfg.AgentPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read agent config: %w", err)
		}

		err = json.Unmarshal(configData, &agentConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to parse agent config: %w", err)
		}
	}

	// Override or set specific fields
	agentConfig.Name = fmt.Sprintf("Agent for User %d", userID)
	agentConfig.IdentityIDs = []string{identity.ID}

	// Add human details to memory blocks if provided
	if userName != "" || userBio != "" {
		var details []string
		if userName != "" {
			details = append(details, fmt.Sprintf("Name: %s", userName))
		}
		if userBio != "" {
			details = append(details, fmt.Sprintf("Bio: %s", userBio))
		}

		humanBlockValue := strings.Join(details, "\n")

		humanBlockFound := false
		for i := range agentConfig.MemoryBlocks {
			if agentConfig.MemoryBlocks[i].Label == "human" {
				agentConfig.MemoryBlocks[i].Value += "\n" + humanBlockValue
				humanBlockFound = true
				break
			}
		}

		if !humanBlockFound {
			humanBlock := letta.MemoryBlock{
				Label: "human",
				Limit: 5000,
				Value: humanBlockValue,
			}

			if agentConfig.MemoryBlocks == nil {
				agentConfig.MemoryBlocks = make([]letta.MemoryBlock, 0)
			}

			agentConfig.MemoryBlocks = append(agentConfig.MemoryBlocks, humanBlock)
		}
	}

	agent, err := th.api.CreateAgent(ctx, agentConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create agent: %w", err)
	}

	chat = &Chat{
		UserID:     userID,
		IdentityID: identity.ID,
		AgentID:    agent.ID,
		TalkLevel:  TalkLevelAloud,
	}

	if !th.cfg.Release {
		chat.TalkLevel = TalkLevelVerbose
	}

	chat, err = th.db.CreateChat(chat)
	if err != nil {
		return nil, fmt.Errorf("failed to save chat: %w", err)
	}

	th.log.Infow("Created new agent for user",
		"userID", userID,
		"agentID", agent.ID,
		"identityID", identity.ID,
		"userName", userName)

	return chat, nil
}

// getChatForUser retrieves an existing chat for a user or creates a new one
func (th *Thera) getChatForUser(ctx context.Context, userID int64) (*Chat, error) {
	return th.CreateChat(ctx, userID, "", "")
}

func (th *Thera) SendMessage(ctx context.Context, userID int64, messageText string) ([]letta.Message, error) {
	chat, err := th.getChatForUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat for user: %w", err)
	}

	messageReq := letta.MessageRequest{
		Messages: []letta.MessageContent{
			{
				Role:     letta.RoleUser,
				Content:  messageText,
				SenderID: &chat.IdentityID,
			},
		},
	}

	resp, err := th.api.SendMessage(ctx, chat.AgentID, messageReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	th.log.Infow("Message sent",
		"userID", userID,
		"agentID", chat.AgentID,
		"messageLength", len(messageText))

	if chat.TalkLevel == TalkLevelVerbose {
		return resp.Messages, nil
	}

	messages := make([]letta.Message, 0)
	for _, m := range resp.Messages {
		switch m.GetMessageType() {
		case letta.MessageTypeAssistant:
			messages = append(messages, m)
		case letta.MessageTypeReasoning:
			if chat.TalkLevel == TalkLevelThoughts {
				messages = append(messages, m)
			}
		}
	}

	return messages, nil
}

// UpdateTalkLevel updates the talk level for a chat by user ID
func (th *Thera) UpdateTalkLevel(userID int64, newTalkLevel TalkLevel) error {
	return th.db.UpdateChatTalkLevel(userID, newTalkLevel)
}
