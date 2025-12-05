package handlers

import (
	"sccsmsserver/db/pg"
	"sccsmsserver/i18n"
	"sccsmsserver/pub"
	"sccsmsserver/setting"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Publish system public information handler
func PubSystemInformationHandler(c *gin.Context) {
	ResponseWithMsg(c, i18n.StatusOK, pg.ServerPubInfo)
}

// Client tests if the server is running
func PubServerPing(c *gin.Context) {
	ResponseWithMsg(c, i18n.StatusOK, gin.H{
		"name":    setting.Conf.Name,
		"apiPath": pub.APIPath,
	})
}

// Generate Front-end DBID
func GenerateFrontendDBID(c *gin.Context) {
	f := new(pg.FrontDBInfo)
	err := c.ShouldBind(f)
	if err != nil {
		zap.L().Error("DeleteOPsHandler invaid parms", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get operator id
	operatorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, f)
		return
	}
	f.Creator.ID = operatorID
	// Generate
	resStatus, _ = f.Generate()
	// Response
	ResponseWithMsg(c, resStatus, f)
}

// Get Front-end DB information handler
func GetFrontendDBInfo(c *gin.Context) {
	f := new(pg.FrontDBInfo)
	err := c.ShouldBind(f)
	if err != nil {
		zap.L().Error("DeleteOPsHandler invaid parms", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get front-end DB info
	resStatus, _ := f.GetInfo()
	//Response
	ResponseWithMsg(c, resStatus, f)
}
