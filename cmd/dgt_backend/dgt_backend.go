package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/olendril/dgt-backend/api"
	config "github.com/olendril/dgt-backend/internal/config"
	"github.com/olendril/dgt-backend/internal/database"
	"github.com/olendril/dgt-backend/internal/monitoring"
	"github.com/rs/zerolog/log"
	"net/http"
)

func main() {
	// create a type that satisfies the `api.ServerInterface`, which contains an implementation of every operation from the generated code
	server := monitoring.NewServer()

	conf, err := config.NewConfig()

	if err != nil {
		log.Error().Err(err).Msg("Error loading config")
		return
	} else {
		log.Info().Msg("Config loaded")
	}

	_, err = database.NewDatabase(conf.Database)

	if err != nil {
		log.Error().Err(err).Msg("Error connecting to database")
		return
	} else {
		log.Info().Msg("Connected to database")
	}

	r := gin.Default()

	api.RegisterHandlers(r, server)

	s := &http.Server{
		Handler: r,
	}

	s.Addr = fmt.Sprintf("0.0.0.0:%d", conf.Port)

	// And we serve HTTP until the world ends.
	log.Fatal().Err(s.ListenAndServe()).Msg("failed to start http server")
}
