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

	database := Database{db: *db}

	err = database.Migrate()

	if err != nil {
		return nil, err
	}

	return &database, nil
}

func (d *Database) Migrate() error {
	err := d.db.AutoMigrate(&User{})

	if err != nil {
		log.Error().Err(err).Msg("Failed to migrate User")
		return err
	}

	return nil
}
