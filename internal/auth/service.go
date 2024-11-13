package auth

import (
	"github.com/gin-gonic/gin"
	authapi "github.com/olendril/dgt-backend/doc/auth"
	"github.com/olendril/dgt-backend/internal/database"
	"github.com/olendril/dgt-backend/internal/discord"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
	"time"
)

type Service struct {
	discordClient discord.Service
	database      database.Database
	FrontendURL   string
}

func NewService(discordClient discord.Service, database database.Database, frontendURL string) *Service {
	return &Service{
		discordClient: discordClient,
		database:      database,
		FrontendURL:   frontendURL,
	}
}

func (s Service) GetLogin(c *gin.Context) {
	response := authapi.LoginResponse{
		Link: s.discordClient.GetGrantAuthorizationLink(),
	}

	c.JSON(http.StatusOK, response)
}

func (s Service) GetRedirect(c *gin.Context) {
	code, exist := c.GetQuery("code")

	if !exist {
		log.Error().Msg("Code not present in the redirect from discord")
		c.JSON(500, gin.H{})
		return
	}

	token, err := s.discordClient.GetAccessToken(code)
	if err != nil || token == nil {
		log.Error().Msg("Error getting access token")
		c.JSON(500, gin.H{})
		return
	}

	discordInfos, err := s.discordClient.GetUserInfo(token.AccessToken)
	if err != nil {
		log.Error().Err(err).Msg("Error getting user info")
		c.JSON(500, gin.H{})
		return
	}

	log.Debug().Interface("infos", discordInfos).Msg("Successfully got user info")

	user, err := s.database.SearchUserByDiscordID(discordInfos.ID)

	if err != nil {
		log.Error().Err(err).Msg("Error searching user")
		c.JSON(500, gin.H{})
		return
	}

	expiration, _ := time.ParseDuration(strconv.Itoa(token.ExpiresIn) + "s")
	expirationDate := time.Now().Add(expiration)

	if user == nil {

		newUser := database.User{
			DiscordID:    discordInfos.ID,
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
			Expiration:   expirationDate,
		}

		err := s.database.CreateUser(newUser)
		if err != nil {
			log.Error().Err(err).Msg("Error creating new user")
			c.JSON(500, gin.H{})
			return
		}
	} else {
		user.AccessToken = token.AccessToken
		user.RefreshToken = token.RefreshToken
		user.Expiration = expirationDate

		err := s.database.UpdateUser(*user)
		if err != nil {
			log.Error().Err(err).Msg("Error updating user")
			c.JSON(500, gin.H{})
			return
		}
	}

	// avatarURL := fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", discordInfos.ID, discordInfos.Avatar)

	c.SetCookie("access_token", token.AccessToken, token.ExpiresIn, "/", "localhost", true, true)
	c.Redirect(http.StatusFound, s.FrontendURL)
}
