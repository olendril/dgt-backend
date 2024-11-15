package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	auth_api "github.com/olendril/dgt-backend/doc/auth"
	guild_api "github.com/olendril/dgt-backend/doc/guilds"
	monitoring_api "github.com/olendril/dgt-backend/doc/monitoring"
	"github.com/olendril/dgt-backend/internal/auth"
	config "github.com/olendril/dgt-backend/internal/config"
	"github.com/olendril/dgt-backend/internal/database"
	"github.com/olendril/dgt-backend/internal/discord"
	"github.com/olendril/dgt-backend/internal/guilds"
	"github.com/olendril/dgt-backend/internal/monitoring"
	"github.com/rs/zerolog/log"
	"net/http"
)

func main() {
	// create a type that satisfies the `api.ServerInterface`, which contains an implementation of every operation from the generated code
	monitoringServer := monitoring.NewServer()

	conf, err := config.NewConfig()

	if err != nil {
		log.Error().Err(err).Msg("Error loading config")
		return
	} else {
		log.Info().Msg("Config loaded")
	}

	databaseService, err := database.NewDatabase(conf.Database)

	if err != nil {
		log.Error().Err(err).Msg("Error connecting to database")
		return
	} else {
		log.Info().Msg("Connected to database")
	}

	discordService := discord.NewDiscordService(conf.Discord)

	r := gin.Default()

	authServer := auth.NewService(discordService, *databaseService, conf.FrontendURL)

	guildServer := guilds.NewServer(discordService, *databaseService)

	monitoring_api.RegisterHandlers(r, monitoringServer)
	auth_api.RegisterHandlers(r, authServer)
	guild_api.RegisterHandlers(r, guildServer)

	s := &http.Server{
		Handler: r,
	}

	s.Addr = fmt.Sprintf("0.0.0.0:%d", conf.Port)

	// And we serve HTTP until the world ends.
	log.Fatal().Err(s.ListenAndServe()).Msg("failed to start http server")
}
