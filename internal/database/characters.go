package database

import (
	"errors"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"slices"
)

type Character struct {
	gorm.Model
	Name            string         `json:"name"`
	Server          string         `json:"server"`
	GuildID         uint           `json:"guild_id"`
	UserID          uint           `json:"user_id"`
	Class           string         `json:"class"`
	DungeonsSuccess pq.StringArray `gorm:"type:text[]" json:"dungeons_success"`
	QuestSuccess    pq.StringArray `gorm:"type:text[]" json:"quest_success"`
	Level           uint           `json:"level"`
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

func (d *Database) FindCharacterByID(id string) (*Character, error) {

	var character Character

	result := d.db.Where("id = ?", id).First(&character)

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to fetch Character")
		return nil, errors.New("internal server error when fetching character")
	}

	return &character, nil
}

func (d *Database) AddDungeonSuccess(idSuccess string, character Character) error {

	key := slices.IndexFunc(character.DungeonsSuccess, func(s string) bool {
		return s == idSuccess
	})

	if key >= 0 {
		return nil
	}

	_, err := d.data.GetDungeonSuccess(idSuccess)
	if err != nil {
		return err
	}

	character.DungeonsSuccess = append(character.DungeonsSuccess, idSuccess)

	result := d.db.Save(character)

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to add dungeon success")
		return errors.New("internal server error when adding dungeon success")
	}

	return nil
}
