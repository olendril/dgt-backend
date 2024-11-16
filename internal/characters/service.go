package characters

import (
	"github.com/gin-gonic/gin"
	character_api "github.com/olendril/dgt-backend/doc/character"
	"github.com/olendril/dgt-backend/internal/database"
	"github.com/olendril/dgt-backend/internal/discord"
	"github.com/olendril/dgt-backend/internal/utils"
	"net/http"
	"strconv"
)

type Service struct {
	db      database.Database
	discord discord.Service
}

func NewServer(discordClient discord.Service, database database.Database) Service {
	return Service{
		discord: discordClient,
		db:      database,
	}
}

func (s Service) GetCharacters(c *gin.Context) {
	user, err := utils.CheckAuth(c, s.db)

	if err != nil {
		return
	}

	characters, err := s.db.GetOwnedCharacters(*user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if characters == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "character not found"})
		return
	}

	var responseCharacters []character_api.CharacterResponse
	for _, character := range *characters {
		characterResponse := character_api.CharacterResponse{
			Achievements: character.Achievements,
			Class:        character.Class,
			GuildId:      strconv.Itoa(int(character.GuildID)),
			Level:        int(character.Level),
			Name:         character.Name,
			Server:       character.Server,
		}
		responseCharacters = append(responseCharacters, characterResponse)
	}

	c.JSON(http.StatusOK, responseCharacters)
}

// (POST /characters)
func (s Service) PostCharacters(c *gin.Context) {
	user, err := utils.CheckAuth(c, s.db)

	if err != nil {
		return
	}

	var requestBody character_api.CharacterInfo
	err = c.ShouldBindJSON(&requestBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// find guild by access code
	guild, err := s.db.FindGuildByCode(requestBody.GuildCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Guild not found"})
		return
	}

	err = s.db.CreateCharacter(database.Character{
		Name:         requestBody.Name,
		Server:       requestBody.Server,
		GuildID:      guild.ID,
		UserID:       user.ID,
		Class:        requestBody.Class,
		Achievements: []string{},
		Level:        uint(requestBody.Level),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
