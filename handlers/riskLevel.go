package handlers

import (
	"sccsmsserver/db/pg"
	"sccsmsserver/i18n"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Add Risk Level handler
func AddRLHandler(c *gin.Context) {
	rl := new(pg.RiskLevel)
	err := c.ShouldBind(rl)
	if err != nil {
		zap.L().Error("AddRLHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operator ID
	operatorID, resStatus := GetCurrentUser(c)
	if resStatus != i18n.StatusOK {
		zap.L().Error("AddRLHandler getCurrentUser failed", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInternalError, rl)
		return
	}
	rl.Creator.ID = operatorID
	// Add
	resStatus, _ = rl.Add()
	// Response
	ResponseWithMsg(c, resStatus, rl)
}

// Get Risk Level list handler
func GetRLListHandler(c *gin.Context) {
	// Get Risk Level list
	rls, resStatus, _ := pg.GetRLList()
	// Response
	ResponseWithMsg(c, resStatus, rls)
}

// Modify Risk Level handler
func EditRLHandler(c *gin.Context) {
	rl := new(pg.RiskLevel)
	err := c.ShouldBind(rl)
	if err != nil {
		zap.L().Error("EditRLHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operator ID
	operatorID, resStatus := GetCurrentUser(c)
	if resStatus != i18n.StatusOK {
		zap.L().Error("EditRLHandler getCurrentUser failed", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInternalError, rl)
		return
	}
	rl.Modifier.ID = operatorID
	// Modify
	resStatus, _ = rl.Edit()
	// Response
	ResponseWithMsg(c, resStatus, rl)
}

// Delete Risk Level handler
func DeleteRLHandler(c *gin.Context) {
	rl := new(pg.RiskLevel)
	err := c.ShouldBind(rl)
	if err != nil {
		zap.L().Error("DeleteRLHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operator ID
	operatorID, resStatus := GetCurrentUser(c)
	if resStatus != i18n.StatusOK {
		zap.L().Error("DeleteRLHandler getCurrentUser failed", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInternalError, rl)
		return
	}
	rl.Modifier.ID = operatorID
	// Delete
	resStatus, _ = rl.Delete()
	// Resopnse
	ResponseWithMsg(c, resStatus, rl)
}

// CheckRLNameExistHandler 检查风险等级名称是否存在
func CheckRLNameExistHandler(c *gin.Context) {
	rl := new(pg.RiskLevel)
	err := c.ShouldBind(rl)
	if err != nil {
		zap.L().Error("CheckRLNameExistHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Check
	resStatus, _ := rl.CheckNameExist()
	// Response
	ResponseWithMsg(c, resStatus, rl)
}

// Get the latest Risk Level front-end cache handler
func GetRLsCacheHandler(c *gin.Context) {
	rlc := new(pg.RLCache)
	err := c.ShouldBind(rlc)
	if err != nil {
		zap.L().Error("GetRLsCacheHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Retrieve data
	resStatus, _ := rlc.GetRLsCache()
	// Response
	ResponseWithMsg(c, resStatus, rlc)
}

// Batch delete Risk Level handler
func DeleteRLsHandler(c *gin.Context) {
	// Parse the parameters
	rls := new([]pg.RiskLevel)
	err := c.ShouldBind(rls)
	if err != nil {
		zap.L().Error("DeleteRLsHandler invaid parms", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Opeartor ID
	modifyUserId, resStatus := GetCurrentUser(c)
	if resStatus != i18n.StatusOK {
		zap.L().Error("DeleteRLsHandler getCurrentUser failed", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInternalError, rls)
		return
	}
	//批量删除
	resStatus, _ = pg.DeleteRLs(rls, modifyUserId)

	ResponseWithMsg(c, resStatus, rls)
}
