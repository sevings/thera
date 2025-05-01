package letta

import (
	"context"
	"fmt"
	"net/http"
)

// IdentityType represents the allowed types for an identity
type IdentityType string

const (
	IdentityTypeOrg   IdentityType = "org"
	IdentityTypeUser  IdentityType = "user"
	IdentityTypeOther IdentityType = "other"
)

// IdentityProperty represents a property associated with an identity
type IdentityProperty struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// CreateIdentityRequest represents the request payload for creating an identity
type CreateIdentityRequest struct {
	IdentifierKey string             `json:"identifier_key"`
	Name          string             `json:"name"`
	IdentityType  IdentityType       `json:"identity_type"`
	ProjectID     string             `json:"project_id,omitempty"`
	AgentIDs      []string           `json:"agent_ids,omitempty"`
	BlockIDs      []string           `json:"block_ids,omitempty"`
	Properties    []IdentityProperty `json:"properties,omitempty"`
}

// Identity represents the response from creating an identity
type Identity struct {
	IdentifierKey string             `json:"identifier_key"`
	Name          string             `json:"name"`
	IdentityType  IdentityType       `json:"identity_type"`
	AgentIDs      []string           `json:"agent_ids,omitempty"`
	BlockIDs      []string           `json:"block_ids,omitempty"`
	ID            string             `json:"id,omitempty"`
	ProjectID     string             `json:"project_id,omitempty"`
	Properties    []IdentityProperty `json:"properties,omitempty"`
}

// validate performs basic validation on the create identity request
func (req CreateIdentityRequest) validate() error {
	if req.IdentifierKey == "" {
		return fmt.Errorf("identifier_key is required")
	}

	if req.Name == "" {
		return fmt.Errorf("name is required")
	}

	switch req.IdentityType {
	case IdentityTypeOrg, IdentityTypeUser, IdentityTypeOther:
		// Valid identity type
	default:
		return fmt.Errorf("invalid identity_type: %s", req.IdentityType)
	}

	return nil
}

// CreateIdentity creates a new identity
func (c *APIClient) CreateIdentity(ctx context.Context, request CreateIdentityRequest) (*Identity, error) {
	if err := request.validate(); err != nil {
		return nil, err
	}

	apiReq := APIRequest{
		Method:   http.MethodPost,
		Endpoint: "/v1/identities/",
		Body:     request,
	}

	var response Identity
	err := c.sendRequest(ctx, apiReq, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// UpsertIdentity creates a new identity or updates existing one
func (c *APIClient) UpsertIdentity(ctx context.Context, request CreateIdentityRequest) (*Identity, error) {
	if err := request.validate(); err != nil {
		return nil, err
	}

	apiReq := APIRequest{
		Method:   http.MethodPost,
		Endpoint: "/v1/identities/",
		Body:     request,
	}

	var response Identity
	err := c.sendRequest(ctx, apiReq, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// RetrieveIdentity returns an existing identity
func (c *APIClient) RetrieveIdentity(ctx context.Context, id string) (*Identity, error) {
	apiReq := APIRequest{
		Method:   http.MethodGet,
		Endpoint: fmt.Sprintf("/v1/identities/%s", id),
	}

	var response Identity
	err := c.sendRequest(ctx, apiReq, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// DeleteIdentity removes an existing identity
func (c *APIClient) DeleteIdentity(ctx context.Context, id string) error {
	apiReq := APIRequest{
		Method:   http.MethodDelete,
		Endpoint: fmt.Sprintf("/v1/identities/%s", id),
	}

	return c.sendRequest(ctx, apiReq, nil)
}

// CountIdentities returns the total number of identities
func (c *APIClient) CountIdentities(ctx context.Context) (int, error) {
	apiReq := APIRequest{
		Method:   http.MethodGet,
		Endpoint: "/v1/identities/count",
	}

	var response int
	err := c.sendRequest(ctx, apiReq, &response)
	if err != nil {
		return 0, err
	}

	return response, nil
}
