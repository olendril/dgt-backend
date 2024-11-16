package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/olendril/dgt-backend/internal/database"
	"strings"
	"time"
)

func CheckAuth(c *gin.Context, database database.Database) (*database.User, error) {
	token := extractToken(c)

	if token == "" {
		c.JSON(401, gin.H{
			"error": "Token is missing",
		})
		return nil, errors.New("token is empty")
	}

	user, err := database.SearchUserByAccessToken(token)
	if err != nil || user == nil {
		c.JSON(401, gin.H{
			"error": "User not found",
		})
		return nil, err
	}

	if time.Now().After(user.Expiration) {
		c.JSON(401, gin.H{
			"error": "Token expired",
		})
		return nil, errors.New("token expired")
	}

	return user, nil
}

func extractToken(c *gin.Context) string {
	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}
