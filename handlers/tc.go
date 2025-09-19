package handlers

import (
	"sccsmsserver/db/pg"
	"sccsmsserver/i18n"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Add Training Crouse handler
func AddTCHandler(c *gin.Context) {
	tc := new(pg.TC)
	err := c.ShouldBind(tc)
	if err != nil {
		zap.L().Error("AddTCHandler invalid param:", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operator ID
	operatorID, resStatus := GetCurrentUser(c)
	if resStatus != i18n.StatusOK {
		zap.L().Error("AddTCHandler getCurrentUser failed:", zap.Error(err))
		ResponseWithMsg(c, resStatus, tc)
		return
	}
	tc.Creator.ID = operatorID
	// Add
	resStatus, _ = tc.Add()
	// Response
	ResponseWithMsg(c, resStatus, tc)
}

// Modify Training course handler
func EditTCHandler(c *gin.Context) {
	tc := new(pg.TC)
	err := c.ShouldBind(tc)
	if err != nil {
		zap.L().Error("EditTCHandler invalid param:", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operator ID
	operatorID, resStatus := GetCurrentUser(c)
	if resStatus != i18n.StatusOK {
		zap.L().Error("EditTCHandler getCurrentUser failed:", zap.Error(err))
		ResponseWithMsg(c, resStatus, tc)
		return
	}
	tc.Modifier.ID = operatorID
	// Modify
	resStatus, _ = tc.Edit()
	// Response
	ResponseWithMsg(c, resStatus, tc)
}

// Delete Training Course Handler
func DeleteTCHandler(c *gin.Context) {
	tc := new(pg.TC)
	err := c.ShouldBind(tc)
	if err != nil {
		zap.L().Error("DeleteTCHandler invalid param:", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operator ID
	operatorID, resStatus := GetCurrentUser(c)
	if resStatus != i18n.StatusOK {
		zap.L().Error("DeleteTCHandler getCurrentUser failed:", zap.Error(err))
		ResponseWithMsg(c, resStatus, tc)
		return
	}
	// Delete
	resStatus, _ = tc.Delete(operatorID)
	// Response
	ResponseWithMsg(c, resStatus, tc)
}

// Batch delete Training Courses handler
func DeleteTCsHandler(c *gin.Context) {
	// Get parameters
	tcs := new([]pg.TC)
	err := c.ShouldBind(tcs)
	if err != nil {
		zap.L().Error("DeleteTCsHandler invaid parms", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
	}
	// Get Operator ID
	operatorID, resStatus := GetCurrentUser(c)
	if resStatus != i18n.StatusOK {
		zap.L().Error("DeleteTCsHandler getCurrentUser failed", zap.Error(err))
		ResponseWithMsg(c, resStatus, tcs)
		return
	}
	// Batch delete
	resStatus, _ = pg.DeleteTCs(tcs, operatorID)
	// Response
	ResponseWithMsg(c, resStatus, tcs)
}

// Get Training Course List Handler
func GetTCListHandler(c *gin.Context) {
	tcs, resStatus, _ := pg.GetTCList()
	ResponseWithMsg(c, resStatus, tcs)
}

// Get Training Course frontend cache handler
func GetTCCacheHandler(c *gin.Context) {
	tcc := new(pg.TCCache)
	err := c.ShouldBind(tcc)
	if err != nil {
		zap.L().Error("GetTCCacheHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Cache
	resStatus, _ := tcc.GetTCCache()
	// Response
	ResponseWithMsg(c, resStatus, tcc)
}

// Check Training Course Name Exist Handler
func CheckTCNameExistHandler(c *gin.Context) {
	tc := new(pg.TC)
	err := c.ShouldBind(tc)
	if err != nil {
		zap.L().Error("CheckTCNameExistHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Check Name Exist
	resStatus, _ := tc.CheckNameExist()
	// Response
	ResponseWithMsg(c, resStatus, tc)
}
