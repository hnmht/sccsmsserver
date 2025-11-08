package handlers

import (
	"sccsmsserver/db/pg"
	"sccsmsserver/i18n"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Get Construction Site list handler
func GetCSsHandler(c *gin.Context) {
	css, resStatus, _ := pg.GetCSs()
	ResponseWithMsg(c, resStatus, css)
}

// Add Constructor Site handler
func AddCSHandler(c *gin.Context) {
	cs := new(pg.ConstructionSite)
	err := c.ShouldBind(cs)
	if err != nil {
		zap.L().Error("AddCSHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get operator ID
	operatorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, cs)
		return
	}
	cs.Creator.ID = operatorID
	// Add
	resStatus, _ = cs.Add()
	ResponseWithMsg(c, resStatus, cs)
}

// Modify Constructor Site handler
func EditCSHandler(c *gin.Context) {
	cs := new(pg.ConstructionSite)
	err := c.ShouldBind(cs)
	if err != nil {
		zap.L().Error("EditCSHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get operator ID
	operationID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, cs)
		return
	}
	cs.Modifier.ID = operationID
	// Modify
	resStatus, _ = cs.Edit()
	// Response
	ResponseWithMsg(c, resStatus, cs)
}

// Get the latest Constructor Site front-end cache handler
func GetCSCacheHandler(c *gin.Context) {
	csc := new(pg.ConstructionSiteCache)
	err := c.ShouldBind(csc)
	if err != nil {
		zap.L().Error("GetCSCacheHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get the cache
	resStatus, _ := csc.GetCSCache()
	// Response
	ResponseWithMsg(c, resStatus, csc)
}

// Check if the Construction Site Code exist handler
func CheckCSCodeExistHandler(c *gin.Context) {
	cs := new(pg.ConstructionSite)
	err := c.ShouldBind(cs)
	if err != nil {
		zap.L().Error("CheckCSCodeExistHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Check
	resStatus, _ := cs.CheckCodeExist()
	// Response
	ResponseWithMsg(c, resStatus, cs)
}

// Delete Construction Site handler
func DeleteCSHandler(c *gin.Context) {
	cs := new(pg.ConstructionSite)
	err := c.ShouldBind(cs)
	if err != nil {
		zap.L().Error("DeleteCSHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get operation ID
	operationID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, cs)
		return
	}
	cs.Modifier.ID = operationID
	// Delete
	resStatus, _ = cs.Delete()
	// Response
	ResponseWithMsg(c, resStatus, cs)
}

// Batch delete Construction Site handler
func DeleteCSsHandler(c *gin.Context) {
	css := new([]pg.ConstructionSite)
	err := c.ShouldBind(css)
	if err != nil {
		zap.L().Error("DeleteCSsHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get operator ID
	operationID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, css)
		return
	}
	// Batch Delete
	resStatus, _ = pg.DeleteCSs(css, operationID)
	// Response
	ResponseWithMsg(c, resStatus, css)
}
