package config

import (
	"context"
	"github.com/pkg/errors"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	Port           int            `env:"PORT, default=8080"`
	Database       DatabaseConfig `env:", prefix=DATABASE_"`
	DiscordBaseURL string         `env:"DISCORD_BASE_URL, default=https://discord.com/api/v10"`
}

type DatabaseConfig struct {
	Host     string `env:"HOST, default=db"`
	Port     string `env:"PORT, default=5432"`
	User     string `env:"USER, default=root"`
	Password string `env:"PASSWORD, default=root"`
	Name     string `env:"NAME, default=dgt"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process(context.Background(), &cfg)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to build config from env")
	}
	return &cfg, nil
}
