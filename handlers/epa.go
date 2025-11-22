package handlers

import (
	"sccsmsserver/db/pg"
	"sccsmsserver/i18n"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Get Execution Project List handler
func GetEPListHandler(c *gin.Context) {
	// Get EP list
	epas, resStatus, _ := pg.GetEPList()
	// Response
	ResponseWithMsg(c, resStatus, epas)
}

// Get Latest Execution Project Front-end Cache handler
func GetEPCacheHandler(c *gin.Context) {
	epac := new(pg.EPCache)
	err := c.ShouldBind(epac)
	if err != nil {
		zap.L().Error("GetEPCacheHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get front-end cache
	resStatus, _ := epac.GetEPCache()
	// Response
	ResponseWithMsg(c, resStatus, epac)
}

// Check Execution Project name exists handler
func CheckEPCodeExistHandler(c *gin.Context) {
	epa := new(pg.ExecutionProject)
	err := c.ShouldBind(epa)
	if err != nil {
		zap.L().Error("CheckEPCodeExistHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Check Code
	resStatus, _ := epa.CheckCodeExist()
	// Response
	ResponseWithMsg(c, resStatus, epa)
}

// Add Execution Project master data handler
func AddEPHandler(c *gin.Context) {
	// Parse the parameters
	epa := new(pg.ExecutionProject)
	err := c.ShouldBind(epa)
	if err != nil {
		zap.L().Error("AddEPHandler invaid parms", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operator ID
	opeartorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, epa)
		return
	}
	epa.Creator.ID = opeartorID
	// Add EP
	resStatus, _ = epa.Add()
	// Response
	ResponseWithMsg(c, resStatus, epa)
}

// Edit Execution Project master data handler
func EditEPHandler(c *gin.Context) {
	// Parse the parameters
	epa := new(pg.ExecutionProject)
	err := c.ShouldBind(epa)
	if err != nil {
		zap.L().Error("EditEPHandler invaid parms", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operator ID
	opeartorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, epa)
		return
	}
	epa.Modifier.ID = opeartorID

	// Modify
	resStats, _ := epa.Edit()
	// Response
	ResponseWithMsg(c, resStats, epa)
}

// Delete Execution Project master data handler
func DeleteEPHandler(c *gin.Context) {
	// Get the parameters
	epa := new(pg.ExecutionProject)
	err := c.ShouldBind(epa)
	if err != nil {
		zap.L().Error("DeleteEPHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get operator ID
	opeartorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, epa)
		return
	}
	epa.Modifier.ID = opeartorID
	// Delete
	resStatus, _ = epa.Delete()
	// Response
	ResponseWithMsg(c, resStatus, epa)
}

// Batch delete Execution Project master datas handler
func DeleteEPsHandler(c *gin.Context) {
	// Get the parameters
	epas := new([]pg.ExecutionProject)
	err := c.ShouldBind(epas)
	if err != nil {
		zap.L().Error("DeleteEPsHandler invaid parms", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
	}
	// Get operator ID
	opeartorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, epas)
		return
	}
	// Batch delete
	resStatus, _ = pg.DeleteEPs(epas, opeartorID)
	// Response
	ResponseWithMsg(c, resStatus, epas)
}
