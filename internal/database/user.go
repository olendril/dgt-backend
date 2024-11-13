package database

import (
	"errors"
	"github.com/rs/zerolog/log"
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

func (d *Database) SearchUserByDiscordID(discordID string) (*User, error) {
	var user User

	result := d.db.Where("discord_id = ?", discordID).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Debug().Str("discordID", discordID).Msg("User Not Found")
		return nil, nil
	} else if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to search User")
		return nil, result.Error
	}

	return &user, nil
}

func (d *Database) CreateUser(user User) error {
	result := d.db.Create(&user)

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to create User")
		return errors.New("internal server error when creating user")
	}

	return nil
}

func (d *Database) UpdateUser(user User) error {
	result := d.db.Save(&user)

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to Updating User")
		return errors.New("internal server error when updating user")
	}

	return nil
}
