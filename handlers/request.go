package handlers

import (
	"sccsmsserver/i18n"
	"sccsmsserver/pub"

	"github.com/gin-gonic/gin"
)

func GetCurrentUser(c *gin.Context) (userID int32, resStatus i18n.ResKey) {
	resStatus = i18n.CodeSuccess
	uid, ok := c.Get(pub.CTXUserID)
	if !ok {
		resStatus = i18n.CodeNeedLogin
		return
	}
	userID, ok = uid.(int32)
	if !ok {
		resStatus = i18n.CodeNeedLogin
		return
	}

	return
}
