package datasets

import (
	"embed"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type DungeonSuccessArray []DungeonSuccess

type DungeonSuccess struct {
	Name    string            `json:"name"`
	Level   int               `json:"level"`
	Success map[string]string `json:"success"`
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
	tmp := strings.Split(id, "-")

	// The id must be in format dungeon-success
	if len(tmp) != 2 {
		return nil, errors.New("invalid id")
	}

	dungeon, ok := s.dungeons[tmp[0]]
	if !ok {
		return nil, errors.New("dungeon not found")
	}

	_, ok = dungeon.Success[tmp[1]]
	if !ok {
		return nil, errors.New("success not found")
	}

	return &dungeon, nil
}

func (s Service) GetSuccessFromDungeons(id string) ([]string, error) {

	dungeon, ok := s.dungeons[id]
	if !ok {
		return nil, errors.New("dungeon not found")
	}

	var ids []string

	for key, _ := range dungeon.Success {
		ids = append(ids, key)
	}

	return ids, nil
}

func (s Service) GetSuccessDungeons(c *gin.Context) {
	c.JSON(http.StatusOK, s.dungeons)
}
