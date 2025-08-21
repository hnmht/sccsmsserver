package handlers

import (
	"sccsmsserver/db/pg"
	"sccsmsserver/i18n"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Add UDC Handler
func AddUDCHandler(c *gin.Context) {
	udc := new(pg.UserDefineCategory)
	err := c.ShouldBind(udc)
	if err != nil {
		zap.L().Error("AddUDCHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get operator ID
	creatorID, resStatus := GetCurrentUser(c)
	if resStatus != i18n.StatusOK {
		zap.L().Error("AddUDCHandler getCurrentUser failed", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInternalError, udc)
		return
	}
	udc.Creator.ID = creatorID
	// Add
	resStatus, _ = udc.Add()
	ResponseWithMsg(c, resStatus, udc)
}

// Get User-define Category master data list handler
func GetUDCListHandler(c *gin.Context) {
	udcs, resStatus, _ := pg.GetUDCList()
	ResponseWithMsg(c, resStatus, udcs)
}

// Edit User-define Category handler
func EditUDCHandler(c *gin.Context) {
	udc := new(pg.UserDefineCategory)
	err := c.ShouldBind(udc)
	if err != nil {
		zap.L().Error("EditUDCHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get operator ID
	modifierID, resStatus := GetCurrentUser(c)
	if resStatus != i18n.StatusOK {
		zap.L().Error("EditUDCHandler getCurrentUser failed", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInternalError, udc)
		return
	}
	udc.Modifier.ID = modifierID
	// Modify
	resStatus, _ = udc.Edit()

	ResponseWithMsg(c, resStatus, udc)
}

// Delete UDC handler
func DeleteUDCHandler(c *gin.Context) {
	udc := new(pg.UserDefineCategory)
	err := c.ShouldBind(udc)
	if err != nil {
		zap.L().Error("DeleteUDCHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operator ID
	modifierID, resStatus := GetCurrentUser(c)
	if resStatus != i18n.StatusOK {
		zap.L().Error("DeleteUDCHandler getCurrentUser failed", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInternalError, udc)
		return
	}
	udc.Modifier.ID = modifierID
	// Delete
	resStatus, _ = udc.Delete()
	ResponseWithMsg(c, resStatus, udc)
}

// Check if the UDC name exists
func CheckUDCNameExistHandler(c *gin.Context) {
	udc := new(pg.UserDefineCategory)
	err := c.ShouldBind(udc)
	if err != nil {
		zap.L().Error("CheckUDCNameExistHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Check
	resStatus, _ := udc.CheckNameExist()
	ResponseWithMsg(c, resStatus, udc)
}

// Get latest UDC for front-end cache handler
func GetUDCsCacheHandler(c *gin.Context) {
	udcc := new(pg.UDCCache)
	err := c.ShouldBind(udcc)
	if err != nil {
		zap.L().Error("GetUDCsCacheHandler invalid param:", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get data
	resStatus, _ := udcc.GetUDCsCache()
	ResponseWithMsg(c, resStatus, udcc)
}

// Batch delete UDC handler
func DeleteUDCsHandler(c *gin.Context) {
	udcs := new([]pg.UserDefineCategory)
	err := c.ShouldBind(udcs)
	if err != nil {
		zap.L().Error("DelUDCsHandler invaid parms:", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, nil)
	}
	// Get operator ID
	modifierID, resStatus := GetCurrentUser(c)
	if resStatus != i18n.StatusOK {
		zap.L().Error("DelUDCsHandler getCurrentUser failed", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInternalError, udcs)
		return
	}
	// Batch Delete
	statusCode, _ := pg.DeleteUDCs(udcs, modifierID)
	// Resopnse
	ResponseWithMsg(c, statusCode, udcs)
}
