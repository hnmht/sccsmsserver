package handlers

import (
	"sccsmsserver/db/pg"
	"sccsmsserver/i18n"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Add Personal Protective Equipment handler
func AddPPEHandler(c *gin.Context) {
	ppe := new(pg.PPE)
	err := c.ShouldBind(ppe)
	if err != nil {
		zap.L().Error("AddPPEHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operator ID
	operatorID, resStatus := GetCurrentUser(c)
	if resStatus != i18n.StatusOK {
		zap.L().Error("AddPPEHandler getCurrentUser failed", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInternalError, ppe)
		return
	}
	ppe.Creator.ID = operatorID
	// Add
	resStatus, _ = ppe.Add()
	// Response
	ResponseWithMsg(c, resStatus, ppe)
}

// Get PPE list handler
func GetPPEListHandler(c *gin.Context) {
	ppes, resStatus, _ := pg.GetPPEList()
	ResponseWithMsg(c, resStatus, ppes)
}

// Modify PPE master data handler
func EditPPEHandler(c *gin.Context) {
	ppe := new(pg.PPE)
	err := c.ShouldBind(ppe)
	if err != nil {
		zap.L().Error("EditPPEHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operator ID
	operatorID, resStatus := GetCurrentUser(c)
	if resStatus != i18n.StatusOK {
		zap.L().Error("EditPPEHandler getCurrentUser failed", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInternalError, ppe)
		return
	}
	ppe.Modifier.ID = operatorID
	// Modify
	resStatus, _ = ppe.Edit()
	// Response
	ResponseWithMsg(c, resStatus, ppe)
}

// Delete PPE master data handler
func DeletePPEHandler(c *gin.Context) {
	ppe := new(pg.PPE)
	err := c.ShouldBind(ppe)
	if err != nil {
		zap.L().Error("DeletePPEHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operator ID
	operatorID, resStatus := GetCurrentUser(c)
	if resStatus != i18n.StatusOK {
		zap.L().Error("DeletePPEHandler getCurrentUser failed", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInternalError, ppe)
		return
	}
	ppe.Modifier.ID = operatorID
	// Delete
	resStatus, _ = ppe.Delete()
	// Response
	ResponseWithMsg(c, resStatus, ppe)
}

// Check PPE code handler
func CheckPPECodeExistHandler(c *gin.Context) {
	ppe := new(pg.PPE)
	err := c.ShouldBind(ppe)
	if err != nil {
		zap.L().Error("CheckPPECodeExistHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Check
	resStatus, _ := ppe.CheckCodeExist()
	// Response
	ResponseWithMsg(c, resStatus, ppe)
}

// Get front-end PPE cache handler
func GetPPECacheHandler(c *gin.Context) {
	ppec := new(pg.PPECache)
	err := c.ShouldBind(ppec)
	if err != nil {
		zap.L().Error("GetPPECacheHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}

	// Get latest ppe data
	resStatus, _ := ppec.GetPPEsCache()
	// Response
	ResponseWithMsg(c, resStatus, ppec)
}

// Batch delete PPE handler
func DeletePPEsHandler(c *gin.Context) {
	ppes := new([]pg.PPE)
	err := c.ShouldBind(ppes)
	if err != nil {
		zap.L().Error("DeletePPEsHandler invaid parms", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operator ID
	operatorID, resStatus := GetCurrentUser(c)
	if resStatus != i18n.StatusOK {
		zap.L().Error("DelUDCsHandler getCurrentUser failed", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInternalError, ppes)
		return
	}
	// Batch Delete
	statusCode, _ := pg.DeletePPEs(ppes, operatorID)
	// Response
	ResponseWithMsg(c, statusCode, ppes)
}
