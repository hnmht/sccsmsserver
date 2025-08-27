package handlers

import (
	"sccsmsserver/db/pg"
	"sccsmsserver/i18n"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Get Constraction Site Options
func GetCSOsHandler(c *gin.Context) {
	sios, resStatus, _ := pg.GetCSOs()
	ResponseWithMsg(c, resStatus, sios)
}

// Modify Construction Site Option Handler
func EditCSOHandler(c *gin.Context) {
	cso := new(pg.ConstructionSiteOption)
	err := c.ShouldBind(cso)
	if err != nil {
		zap.L().Error("EditCSOHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operator ID
	operatorID, resStatus := GetCurrentUser(c)
	if resStatus != i18n.StatusOK {
		zap.L().Error("EditCSOHandler getCurrentUser failed", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInternalError, cso)
		return
	}
	cso.Modifier.ID = operatorID

	resStatus, _ = cso.Edit()
	ResponseWithMsg(c, resStatus, cso)
}

// Get the latest Construction Site Options front-end cache handler
func GetCSOCacheHandler(c *gin.Context) {
	csoc := new(pg.ConstructionSiteOptionCache)
	err := c.ShouldBind(csoc)
	if err != nil {
		zap.L().Error("GetCSOCacheHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get the cache
	resStatus, _ := csoc.GetCSOCache()
	// Response
	ResponseWithMsg(c, resStatus, csoc)
}
