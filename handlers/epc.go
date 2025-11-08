package handlers

import (
	"sccsmsserver/db/pg"
	"sccsmsserver/i18n"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Get Execution Project Category List handler
func GetEPCListHandler(c *gin.Context) {
	epcs, resStatus, _ := pg.GetEPCList()
	ResponseWithMsg(c, resStatus, epcs)
}

// Get Simple Execution Project Category list handler
func GetSimpEPCListHandler(c *gin.Context) {
	sepcs, resStatus, _ := pg.GetSimpEPCList()
	ResponseWithMsg(c, resStatus, sepcs)
}

// Get Simple Execution Project Category front-end cache handler
func GetSimpEPCCacheHandler(c *gin.Context) {
	sepcc := new(pg.SimpEPCCache)
	err := c.ShouldBind(sepcc)
	if err != nil {
		zap.L().Error("GetSimpEPCCacheHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get data
	resStatus, _ := sepcc.GetSimpEPCCache()
	// Response
	ResponseWithMsg(c, resStatus, sepcc)
}

// Check if the EPC name exists handler
func CheckEPCNameExistHandler(c *gin.Context) {
	epc := new(pg.EPC)
	err := c.ShouldBind(epc)
	if err != nil {
		zap.L().Error("CheckEPCNameExistHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Check
	resStatus, _ := epc.CheckNameExist()
	ResponseWithMsg(c, resStatus, epc)
}

// Add EPC handler
func AddEPCHandler(c *gin.Context) {
	epc := new(pg.EPC)
	err := c.ShouldBind(epc)
	if err != nil {
		zap.L().Error("AddEPCHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operator ID
	operatorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, epc)
		return
	}
	epc.Creator.ID = operatorID
	// Add
	resStatus, _ = epc.Add()
	// Response
	ResponseWithMsg(c, resStatus, epc)
}

// Edit EPC handler
func EditEPCHandler(c *gin.Context) {
	epc := new(pg.EPC)
	err := c.ShouldBind(epc)
	if err != nil {
		zap.L().Error("EditEPCHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operator ID
	operatorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, epc)
		return
	}
	epc.Modifier.ID = operatorID
	// Modify
	resStatus, _ = epc.Edit()
	// Response
	ResponseWithMsg(c, resStatus, epc)
}

// Delete EPC handler
func DeleteEPCHandler(c *gin.Context) {
	// Get parameter
	epc := new(pg.EPC)
	err := c.ShouldBind(epc)
	if err != nil {
		zap.L().Error("DeleteEPCHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operator ID
	operatorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, epc)
		return
	}
	epc.Modifier.ID = operatorID
	// Delete
	resStatus, _ = epc.Delete()
	// Response
	ResponseWithMsg(c, resStatus, epc)
}

// Batch delete EPC handler
func DeleteEPCsHandler(c *gin.Context) {
	// Get parameter
	epcs := new([]pg.EPC)
	err := c.ShouldBind(epcs)
	if err != nil {
		zap.L().Error("DeleteEPCsHandler invaid parms", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operator ID
	operatorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, epcs)
		return
	}
	// Batch delete
	resStatus, _ = pg.DeleteEPCs(epcs, operatorID)
	// Response
	ResponseWithMsg(c, resStatus, epcs)
}
