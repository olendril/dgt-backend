package characters

import (
	"github.com/gin-gonic/gin"
	character_api "github.com/olendril/dgt-backend/doc/character"
	"github.com/olendril/dgt-backend/internal/database"
	"github.com/olendril/dgt-backend/internal/discord"
	"github.com/olendril/dgt-backend/internal/utils"
	"github.com/rs/zerolog/log"
	"net/http"
	"slices"
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
			Id:             strconv.Itoa(int(character.ID)),
			DungeonSuccess: character.DungeonsSuccess,
			Class:          character.Class,
			GuildId:        strconv.Itoa(int(character.GuildID)),
			Level:          int(character.Level),
			Name:           character.Name,
			Server:         character.Server,
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if guild == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Guild not found"})
		return
	}

	err = s.db.CreateCharacter(database.Character{
		Name:            requestBody.Name,
		Server:          requestBody.Server,
		GuildID:         guild.ID,
		UserID:          user.ID,
		Class:           requestBody.Class,
		DungeonsSuccess: []string{},
		Level:           uint(requestBody.Level),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (s Service) GetCharactersId(c *gin.Context, id string) {
	_, err := utils.CheckAuth(c, s.db)
	if err != nil {
		return
	}

	character, err := s.db.FindCharacterByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if character == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "character not found"})
		return
	}

	characterResponse := character_api.CharacterResponse{
		Id:             strconv.Itoa(int(character.ID)),
		Class:          character.Class,
		DungeonSuccess: character.DungeonsSuccess,
		GuildId:        strconv.Itoa(int(character.GuildID)),
		Level:          int(character.Level),
		Name:           character.Name,
		Server:         character.Server,
	}

	c.JSON(200, characterResponse)
}

func (s Service) PostCharactersIdSuccessSuccessID(c *gin.Context, id string, successID string) {
	user, err := utils.CheckAuth(c, s.db)
	if err != nil {
		return
	}

	character, err := s.db.FindCharacterByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if character == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "character not found"})
		return
	}

	idParsed, err := strconv.Atoi(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Info().Interface("id", idParsed).Send()

	key := slices.IndexFunc(user.Characters, func(s database.Character) bool {
		return int(s.ID) == idParsed
	})

	if key < 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "character doesn't belong to user"})
	}

	err = s.db.AddDungeonSuccess(successID, *character)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{})
}

func (s Service) PutCharactersIdSuccessDungeons(c *gin.Context, id string) {
	user, err := utils.CheckAuth(c, s.db)
	if err != nil {
		return
	}

	var requestBody []string
	err = c.ShouldBindJSON(&requestBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	character, err := s.db.FindCharacterByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if character == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "character not found"})
		return
	}

	idParsed, err := strconv.Atoi(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Info().Interface("id", idParsed).Send()

	key := slices.IndexFunc(user.Characters, func(s database.Character) bool {
		return int(s.ID) == idParsed
	})

	if key < 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "character doesn't belong to user"})
	}

	character.DungeonsSuccess = requestBody

	err = s.db.UpdateCharacter(*character)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{})
}
