package handlers

import (
	"sccsmsserver/db/pg"
	"sccsmsserver/i18n"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Add position handler
func AddPositionHandler(c *gin.Context) {
	p := new(pg.Position)
	err := c.ShouldBind(p)
	if err != nil {
		zap.L().Error("AddOPHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}

	// Get operator id
	operatorID, resStatus := GetCurrentUser(c)
	if resStatus != i18n.StatusOK {
		zap.L().Error("AddOPHandler getCurrentUser failed", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInternalError, p)
		return
	}
	p.Creator.ID = operatorID
	// Add
	resStatus, _ = p.Add()
	// Response
	ResponseWithMsg(c, resStatus, p)
}

// Get position list handler
func GetPositionListHandler(c *gin.Context) {
	ops, resStatus, _ := pg.GetPositionList()
	ResponseWithMsg(c, resStatus, ops)
}

// Edit position handler
func EditPositionHandler(c *gin.Context) {
	p := new(pg.Position)
	err := c.ShouldBind(p)
	if err != nil {
		zap.L().Error("EditOPHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get operator id
	operatorID, resStatus := GetCurrentUser(c)
	if resStatus != i18n.StatusOK {
		zap.L().Error("EditOPHandler getCurrentUser failed", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInternalError, p)
		return
	}
	p.Modifier.ID = operatorID
	// Edit
	resStatus, _ = p.Edit()
	// Response
	ResponseWithMsg(c, resStatus, p)
}

// Delete positon handler
func DeletePositionHandler(c *gin.Context) {
	p := new(pg.Position)
	err := c.ShouldBind(p)
	if err != nil {
		zap.L().Error("DeleteOPHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get operator id
	operatorID, resStatus := GetCurrentUser(c)
	if resStatus != i18n.StatusOK {
		zap.L().Error("DeleteOPHandler getCurrentUser failed", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInternalError, p)
		return
	}
	p.Modifier.ID = operatorID
	// Delete
	resStatus, _ = p.Delete()
	// Response
	ResponseWithMsg(c, resStatus, p)
}

// Check position name exists handler
func CheckPositionNameExistHandler(c *gin.Context) {
	p := new(pg.Position)
	err := c.ShouldBind(p)
	if err != nil {
		zap.L().Error("CheckOPNameExistHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Check
	resStatus, _ := p.CheckNameExist()

	ResponseWithMsg(c, resStatus, p)
}

// Get latest position master data for front-end cache handler
func GetPositionCacheHandler(c *gin.Context) {
	pc := new(pg.PositionCache)
	err := c.ShouldBind(pc)
	if err != nil {
		zap.L().Error("GetOPCacheHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get position data
	resStatus, _ := pc.GetOPsCache()
	// Response
	ResponseWithMsg(c, resStatus, pc)
}

// Batch delete position handler
func DeletePositionsHandler(c *gin.Context) {
	ps := new([]pg.Position)
	err := c.ShouldBind(ps)
	if err != nil {
		zap.L().Error("DeleteOPsHandler invaid parms", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
	}
	// Get operator id
	operatorID, resStatus := GetCurrentUser(c)
	if resStatus != i18n.StatusOK {
		zap.L().Error("DelUDCsHandler getCurrentUser failed", zap.Error(err))
		ResponseWithMsg(c, resStatus, ps)
		return
	}
	// Batch Delete
	resStatus, _ = pg.DeleteOPs(ps, operatorID)
	// Resopnse
	ResponseWithMsg(c, resStatus, ps)
}
