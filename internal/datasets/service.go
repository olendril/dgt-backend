package datasets

import (
	"embed"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type DungeonSuccess struct {
	Name    string `json:"name"`
	Dungeon string `json:"dungeon"`
	Level   int    `json:"level"`
}

type Service struct {
	dungeons map[string]DungeonSuccess
}

//go:embed dungeons_fr.json
var dungeons embed.FS

func NewService() (*Service, error) {

	var dungeonsSuccess map[string]DungeonSuccess

	data, err := dungeons.ReadFile("dungeons_fr.json")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &dungeonsSuccess)

	if err != nil {
		return nil, err
	}

	return &Service{dungeons: dungeonsSuccess}, nil
}

func (s Service) GetDungeonSuccess(id string) (*DungeonSuccess, error) {
	val, ok := s.dungeons[id]
	if !ok {
		return nil, errors.New("success not found")
	}

	return &val, nil
}

func (s Service) GetSuccessDungeons(c *gin.Context) {
	c.JSON(http.StatusOK, s.dungeons)
}
