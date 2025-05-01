package thera

import (
	"context"
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
