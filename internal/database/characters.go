package database

import (
	"errors"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Character struct {
	gorm.Model
	Name         string         `json:"name"`
	Server       string         `json:"server"`
	GuildID      uint           `json:"guild_id"`
	UserID       uint           `json:"user_id"`
	Class        string         `json:"class"`
	Achievements pq.StringArray `gorm:"type:text[]" json:"achievements"`
	Level        uint           `json:"level"`
}

func (d *Database) CreateCharacter(character Character) error {
	result := d.db.Create(&character)

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to create Character")
		return errors.New("internal server error when creating Character")
	}

	return nil
}

func (d *Database) GetOwnedCharacters(owner User) (*[]Character, error) {

	result := d.db.Preload("Characters").First(&owner)

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to get owned characters")
		return nil, errors.New("internal server error when getting owned characters")
	}

	return &owner.Characters, nil
}
