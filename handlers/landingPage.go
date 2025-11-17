package handlers

import (
	"sccsmsserver/db/pg"
	"sccsmsserver/i18n"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Get Landing Page Info handler
func GetLandingPageInfoHandler(c *gin.Context) {
	info, resStatus, _ := pg.GetLandingPageInfo()
	// Response
	ResponseWithMsg(c, resStatus, info)
}

// Modify Landing Page Info handler
func ModifyLandingPageInfoHandler(c *gin.Context) {
	info := new(pg.LandingPageInfo)
	err := c.ShouldBind(info)
	if err != nil {
		zap.L().Error("ModifyLandingPageInfoHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}

	// Get current operator ID
	modifyUserId, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, i18n.CodeInternalError, info)
		return
	}
	info.Modifier.ID = modifyUserId
	// Modify
	resStatus, _ = info.Modify()
	// Resoponse
	ResponseWithMsg(c, resStatus, info)
}
