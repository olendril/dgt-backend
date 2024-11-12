package database

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	DiscordID    string    `json:"discord_id"`
	Expiration   time.Time `json:"expiration"`
}
