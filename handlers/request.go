package handlers

import (
	"sccsmsserver/i18n"
	"sccsmsserver/pub"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetCurrentUser(c *gin.Context) (userID int32, resStatus i18n.ResKey) {
	resStatus = i18n.StatusOK
	uid, ok := c.Get(pub.CTXUserID)
	if !ok {
		zap.L().Error("GetCurrentUser c.Get(pub.CTXUserID) failed.")
		resStatus = i18n.CodeNeedLogin
		return
	}
	userID, ok = uid.(int32)
	if !ok {
		zap.L().Error("GetCurrentUser uid.(int32) failed.")
		resStatus = i18n.CodeNeedLogin
		return
	}

	return
}
