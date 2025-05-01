package thera

import (
	"errors"
	"os"
	"thera/internal/letta"

	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	TgToken string `koanf:"tg_token"`
	DBPath  string `koanf:"db_path"`
	Release bool
	Letta   letta.Config
}

func LoadConfig() (Config, error) {
	var kConf = koanf.New("/")

	var cfg Config

	path := os.Getenv("THERA_CONFIG")
	if path == "" {
		path = "thera.toml"
	}

	err := kConf.Load(file.Provider(path), toml.Parser())
	if err != nil {
		return cfg, err
	}

	err = kConf.Unmarshal("", &cfg)
	if err != nil {
		return cfg, err
	}

	if cfg.TgToken == "" {
		return cfg, errors.New("telegram token is required")
	}

	return cfg, nil
}
