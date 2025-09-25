package handlers

import (
	"sccsmsserver/db/pg"
	"sccsmsserver/i18n"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Get online user list handler
func GetOnlineUserHandler(c *gin.Context) {
	ous, resStatus, _ := pg.GetAllOnlineUser()
	ResponseWithMsg(c, resStatus, ous)
}

// Remove online user handler
func RemoveOnlineUserHandler(c *gin.Context) {
	ou := new(pg.OnlineUser)
	err := c.ShouldBind(ou)
	if err != nil {
		zap.L().Error("RemoveOnlineUserHandler invaid parms", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// remove from cache
	resStatus, _ := ou.Del()
	ResponseWithMsg(c, resStatus, ou)
}
