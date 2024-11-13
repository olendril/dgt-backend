package discord

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/olendril/dgt-backend/internal/config"
	"github.com/rs/zerolog/log"
)

type Service struct {
	conf   config.DiscordConfig
	client resty.Client
}

func NewDiscordService(conf config.DiscordConfig) Service {
	return Service{
		conf:   conf,
		client: *resty.New(),
	}
}

func (s Service) GetGrantAuthorizationLink() string {
	// https://discord.com/developers/docs/topics/oauth2
	return fmt.Sprintf("https://discord.com/oauth2/authorize?response_type=code&client_id=%s&scope=%s&"+
		"redirect_uri=%s&prompt=none&integration_type=1", s.conf.ClientID, s.conf.Scopes, s.conf.Redirect)
}

type GetAccessTokenResponse struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
}

func (s Service) GetAccessToken(code string) (*GetAccessTokenResponse, error) {
	post, err := s.client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{
			"grant_type":    "authorization_code",
			"code":          code,
			"redirect_uri":  s.conf.Redirect,
			"client_id":     s.conf.ClientID,
			"client_secret": s.conf.ClientSecret,
		}).Post("https://discord.com/api/oauth2/token")

	if err != nil {
		log.Error().Err(err).Msg("Error getting access token")
		return nil, err
	}

	var response GetAccessTokenResponse

	err = json.Unmarshal(post.Body(), &response)

	if err != nil {
		return nil, err
	}

	return &response, nil
}

type UserResponse struct {
	ID         string `json:"id"`
	Username   string `json:"username"`
	GlobalName string `json:"global_name"`
	Avatar     string `json:"avatar"`
}

func (s Service) GetUserInfo(accessToken string) (*UserResponse, error) {
	resp, err := s.client.R().
		SetAuthScheme("Bearer").
		SetAuthToken(accessToken).
		Get(s.conf.BaseUrl + "/users/@me")

	if err != nil {
		log.Error().Err(err).Msg("Error getting user info")
	}

	var response UserResponse

	log.Info().Interface("response", resp.Body()).Msg("Response received")

	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
