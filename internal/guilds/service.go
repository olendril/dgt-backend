package guilds

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	guild_api "github.com/olendril/dgt-backend/doc/guilds"
	"github.com/olendril/dgt-backend/internal/database"
	"github.com/olendril/dgt-backend/internal/discord"
	"github.com/olendril/dgt-backend/internal/utils"
	"net/http"
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

	err = s.database.CreateGuild(guild)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, nil)
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
		Code:   guild.Code,
		Name:   guild.Name,
		Server: guild.Server,
	}

	c.JSON(200, guildResponse)
}
