package discord

import (
	"github.com/go-resty/resty/v2"
	"github.com/olendril/dgt-backend/internal/config"
)

type Service struct {
	conf   config.DiscordConfig
	client resty.Client
}

func NewDiscordService(conf config.DiscordConfig) Service {
	return Service{
		conf:   conf,
		client: *resty.New(),
	}
}
