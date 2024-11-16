package database

import (
	"errors"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Guild struct {
	gorm.Model
	Name       string      `json:"name"`
	Server     string      `json:"server"`
	UserID     uint        `json:"user_id"`
	Code       string      `json:"code"`
	Characters []Character `json:"characters"`
}

func (d *Database) CreateGuild(guild Guild) error {
	result := d.db.Create(&guild)

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to create Guild")
		return errors.New("internal server error when creating guild")
	}

	return nil
}

func (d *Database) GetOwnedGuilds(owner User) (*[]Guild, error) {

	result := d.db.Preload("Guilds").First(&owner)

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to create Guild")
		return nil, errors.New("internal server error when creating guild")
	}

	return &owner.Guilds, nil
}

func (d *Database) FindGuildByCode(code string) (*Guild, error) {

	var guild Guild

	result := d.db.Where("code = ?", code).First(&guild)

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to fetch Guild")
		return nil, errors.New("internal server error when fetching guild")
	}

	return &guild, nil
}

func (d *Database) FindGuildByID(id string) (*Guild, error) {

	var guild Guild

	result := d.db.Where("id = ?", id).First(&guild)

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to fetch Guild")
		return nil, errors.New("internal server error when fetching guild")
	}

	return &guild, nil
}
