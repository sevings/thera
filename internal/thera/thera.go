package thera

import (
	"context"
	"fmt"
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

// getChatForUser retrieves an existing chat for a user or creates a new one
func (th *Thera) getChatForUser(ctx context.Context, userID int64) (*Chat, error) {
	chat, err := th.db.GetChatByUserID(userID)
	if err == nil {
		return chat, nil
	}

	if err != ErrNotFound {
		return nil, fmt.Errorf("error checking existing chat: %w", err)
	}

	identityReq := letta.CreateIdentityRequest{
		IdentifierKey: fmt.Sprintf("user_%d", userID),
		Name:          fmt.Sprintf("User %d", userID),
		IdentityType:  letta.IdentityTypeUser,
	}

	identity, err := th.api.CreateIdentity(ctx, identityReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create identity: %w", err)
	}

	agentReq := letta.CreateAgentRequest{
		Name:            fmt.Sprintf("Agent for User %d", userID),
		System:          th.cfg.Agent.System,
		AgentType:       letta.AgentTypeMemgpt,
		IdentityIDs:     []string{identity.ID},
		LLMConfig:       &th.cfg.Model,
		EmbeddingConfig: &th.cfg.Embedding,
	}

	agent, err := th.api.CreateAgent(ctx, agentReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create agent: %w", err)
	}

	chat, err = th.db.CreateChat(userID, identity.ID, agent.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to save chat: %w", err)
	}

	th.log.Infow("Created new agent for user",
		"userID", userID,
		"agentID", agent.ID,
		"identityID", identity.ID)

	return chat, nil
}

func (th *Thera) SendMessage(ctx context.Context, userID int64, messageText string) ([]map[string]any, error) {
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

	return resp.Messages, nil
}
