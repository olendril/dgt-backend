package config

import (
	"context"
	"github.com/pkg/errors"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	Port         int            `env:"PORT, default=8080"`
	FrontendURL  string         `env:"FRONTEND_URL, default=http://localhost:5173"`
	Database     DatabaseConfig `env:", prefix=DATABASE_"`
	Discord      DiscordConfig  `env:", prefix=DISCORD_"`
	CookieDomain string         `env:"COOKIE_DOMAIN, default=localhost"`
}

type DatabaseConfig struct {
	Host     string `env:"HOST, default=db"`
	Port     string `env:"PORT, default=5432"`
	User     string `env:"USER, default=root"`
	Password string `env:"PASSWORD, default=root"`
	Name     string `env:"NAME, default=dgt"`
}

type DiscordConfig struct {
	BaseUrl      string `env:"BASE_URL, default=https://discord.com/api/v10"`
	ClientID     string `env:"CLIENT_ID"`
	ClientSecret string `env:"CLIENT_SECRET"`
	Scopes       string `env:"SCOPES, default=identify"`
	Redirect     string `env:"REDIRECT_URL, default=http://localhost:8080/redirect"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process(context.Background(), &cfg)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to build config from env")
	}
	return &cfg, nil
}
