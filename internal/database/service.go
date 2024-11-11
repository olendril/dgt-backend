package database

import (
	"fmt"
	"github.com/olendril/dgt-backend/internal/config"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	db gorm.DB
}

func NewDatabase(conf config.DatabaseConfig) (*Database, error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", conf.User, conf.Password, conf.Host, conf.Port, conf.Name)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to database")
		return nil, err
	}

	return &Database{db: *db}, nil
}
