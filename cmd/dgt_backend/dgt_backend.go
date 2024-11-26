package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	auth_api "github.com/olendril/dgt-backend/doc/auth"
	character_api "github.com/olendril/dgt-backend/doc/character"
	guild_api "github.com/olendril/dgt-backend/doc/guilds"
	monitoring_api "github.com/olendril/dgt-backend/doc/monitoring"
	success_api "github.com/olendril/dgt-backend/doc/success"
	"github.com/olendril/dgt-backend/internal/auth"
	"github.com/olendril/dgt-backend/internal/characters"
	config "github.com/olendril/dgt-backend/internal/config"
	"github.com/olendril/dgt-backend/internal/database"
	"github.com/olendril/dgt-backend/internal/datasets"
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

	datasetService, err := datasets.NewService()

	if err != nil {
		log.Error().Err(err).Msg("Error loading datasets")
		return
	} else {
		log.Info().Msg("Datasets loaded")
	}

	databaseService, err := database.NewDatabase(conf.Database, *datasetService)

	if err != nil {
		log.Error().Err(err).Msg("Error connecting to database")
		return
	} else {
		log.Info().Msg("Connected to database")
	}

	discordService := discord.NewDiscordService(conf.Discord)

	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Authorization")
	corsConfig.AllowAllOrigins = true
	r.Use(cors.New(corsConfig))

	authServer := auth.NewService(discordService, *databaseService, conf.FrontendURL)

	guildServer := guilds.NewServer(discordService, *databaseService)

	characterServer := characters.NewServer(discordService, *databaseService, *datasetService)

	monitoring_api.RegisterHandlers(r, monitoringServer)
	auth_api.RegisterHandlers(r, authServer)
	guild_api.RegisterHandlers(r, guildServer)
	character_api.RegisterHandlers(r, characterServer)
	success_api.RegisterHandlers(r, datasetService)

	s := &http.Server{
		Handler: r,
	}

	s.Addr = fmt.Sprintf("0.0.0.0:%d", conf.Port)

	// And we serve HTTP until the world ends.
	log.Fatal().Err(s.ListenAndServe()).Msg("failed to start http server")
}
