package monitoring

import (
	"github.com/gin-gonic/gin"
)

// optional code omitted

type Monitoring struct {
}

func NewServer() Monitoring {
	return Monitoring{}
}

// (GET /ping)
func (Monitoring) GetPing(ctx *gin.Context) {
	ctx.JSON(200, nil)
}
