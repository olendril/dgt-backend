package guilds

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	guild_api "github.com/olendril/dgt-backend/doc/guilds"
	"github.com/olendril/dgt-backend/internal/database"
	"github.com/olendril/dgt-backend/internal/discord"
	"github.com/olendril/dgt-backend/internal/utils"
	"github.com/rs/zerolog/log"
	"net/http"
	"slices"
	"strconv"
)

type Service struct {
	discordClient discord.Service
	database      database.Database
}

func NewServer(discordClient discord.Service, database database.Database) Service {
	return Service{
		discordClient: discordClient,
		database:      database,
	}
}

func (s Service) PostGuilds(c *gin.Context) {
	var requestBody guild_api.GuildInfo

	user, err := utils.CheckAuth(c, s.database)
	if err != nil {
		return
	}

	err = c.ShouldBindJSON(&requestBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	code, err := uuid.NewUUID()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	guild := database.Guild{
		Name:   requestBody.Name,
		Server: requestBody.Server,
		UserID: user.ID,
		Code:   code.String(),
	}

	guildResponse, err := s.database.CreateGuild(guild)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, guild_api.GuildResponse{
		Code:   guildResponse.Code,
		Id:     strconv.Itoa(int(guildResponse.ID)),
		Name:   guildResponse.Name,
		Server: guildResponse.Server,
	})
}

func (s Service) GetGuilds(c *gin.Context) {
	user, err := utils.CheckAuth(c, s.database)
	if err != nil {
		return
	}

	var responseGuilds []guild_api.GuildResponse

	guilds, err := s.database.GetOwnedGuilds(*user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if guilds == nil {
		c.JSON(http.StatusOK, responseGuilds)
		return
	}

	for _, value := range *guilds {
		responseGuilds = append(responseGuilds, guild_api.GuildResponse{
			Id:     strconv.Itoa(int(value.ID)),
			Code:   value.Code,
			Name:   value.Name,
			Server: value.Server,
		})
	}

	c.JSON(200, responseGuilds)

}

func (s Service) GetGuildsId(c *gin.Context, id string) {
	_, err := utils.CheckAuth(c, s.database)
	if err != nil {
		return
	}

	guild, err := s.database.FindGuildByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if guild == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "guild not found"})
		return
	}

	guildResponse := guild_api.GuildResponse{
		Id:     strconv.Itoa(int(guild.ID)),
		Code:   guild.Code,
		Name:   guild.Name,
		Server: guild.Server,
	}

	c.JSON(200, guildResponse)
}

func (s Service) GetGuildsIdCharacters(c *gin.Context, id string) {
	c.JSON(501, gin.H{})
}

func (s Service) DeleteGuildsId(c *gin.Context, id string) {
	user, err := utils.CheckAuth(c, s.database)
	if err != nil {
		return
	}

	guild, err := s.database.FindGuildByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if guild == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "guild not found"})
		return
	}

	idParsed, err := strconv.Atoi(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Info().Interface("id", idParsed).Send()

	key := slices.IndexFunc(user.Guilds, func(s database.Guild) bool {
		return int(s.ID) == idParsed
	})

	if key < 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "guild doesn't belong to user"})
	}

	err = s.database.DeleteGuild(*guild)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{})
}
