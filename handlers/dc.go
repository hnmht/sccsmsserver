package handlers

import (
	"sccsmsserver/db/pg"
	"sccsmsserver/i18n"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Get Document Categories List Handler
func GetDCListHandler(c *gin.Context) {
	dcs, resStatus, _ := pg.GetDCList()
	ResponseWithMsg(c, resStatus, dcs)
}

// Get Simple Document Categories List Handler
func GetSimpDCListHandler(c *gin.Context) {
	sdcs, resStatus, _ := pg.GetSimpDCList()
	ResponseWithMsg(c, resStatus, sdcs)
}

// Get Simple Document Categories front-end Cache Handler
func GetSimpDCCacheHandler(c *gin.Context) {
	sdcc := new(pg.SimpDCCache)
	err := c.ShouldBind(sdcc)
	if err != nil {
		zap.L().Error("GetSimpDCCacheHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Cache
	resStatus, _ := sdcc.GetSimpDCCache()
	ResponseWithMsg(c, resStatus, sdcc)
}

// Check Document Category Name Exist Handler
func CheckDCNameExistHandler(c *gin.Context) {
	dc := new(pg.DC)
	err := c.ShouldBind(dc)
	if err != nil {
		zap.L().Error("CheckDCNameExistHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Check Name Exist
	resStatus, _ := dc.CheckNameExist()
	ResponseWithMsg(c, resStatus, dc)
}

// Add Document Category Handler
func AddDCHandler(c *gin.Context) {
	dc := new(pg.DC)
	err := c.ShouldBind(dc)
	if err != nil {
		zap.L().Error("AddDCHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operator ID
	operatorID, resStatus := GetCurrentUser(c)
	if resStatus != i18n.StatusOK {
		zap.L().Error("AddDCHandler getCurrentUser failed", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInternalError, dc)
		return
	}
	dc.Creator.ID = operatorID
	// Add Document Category
	resStatus, _ = dc.Add()
	ResponseWithMsg(c, resStatus, dc)
}

// Edit Document Category Handler
func EditDCHandler(c *gin.Context) {
	dc := new(pg.DC)
	err := c.ShouldBind(dc)
	if err != nil {
		zap.L().Error("EditDCHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}

	// Get Operator ID
	operatorID, resStatus := GetCurrentUser(c)
	if resStatus != i18n.StatusOK {
		zap.L().Error("EditDCHandler getCurrentUser failed", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInternalError, dc)
		return
	}
	dc.Modifier.ID = operatorID

	// Edit Document Category
	resStatus, _ = dc.Edit()
	ResponseWithMsg(c, resStatus, dc)
}

// Delete Document Category Handler
func DeleteDCHandler(c *gin.Context) {
	// Get Parameters
	dc := new(pg.DC)
	err := c.ShouldBind(dc)
	if err != nil {
		zap.L().Error("DeleteDCHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operator ID
	operatorID, resStatus := GetCurrentUser(c)
	if resStatus != i18n.StatusOK {
		zap.L().Error("DeleteEICHandler getCurrentUser failed", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInternalError, dc)
		return
	}
	dc.Modifier.ID = operatorID

	// Delete Document Category
	resStatus, _ = dc.Delete()
	ResponseWithMsg(c, resStatus, dc)
}

// Delete Multiple Document Categories Handler
func DeleteDCsHandler(c *gin.Context) {
	// Get Parameters
	dcs := new([]pg.DC)
	err := c.ShouldBind(dcs)
	if err != nil {
		zap.L().Error("DeleteDCsHandler invaid parms", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operator ID
	operatorID, resStatus := GetCurrentUser(c)
	if resStatus != i18n.StatusOK {
		zap.L().Error("DeleteDCsHandler getCurrentUser failed", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInternalError, dcs)
		return
	}
	// Delete Document Categories
	resStatus, _ = pg.DeleteDCs(dcs, operatorID)
	// R
	ResponseWithMsg(c, resStatus, dcs)
}
