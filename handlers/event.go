package handlers

import (
	"sccsmsserver/db/pg"
	"sccsmsserver/i18n"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Get User Events Handler
func GetEventsHandler(c *gin.Context) {
	ue := new(pg.UserEvents)
	err := c.ShouldBind(ue)
	if err != nil {
		zap.L().Error("GetEventsHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Retrieve Events
	resStatus, _ := ue.GetEvents()
	// Response
	ResponseWithMsg(c, resStatus, ue)
}
