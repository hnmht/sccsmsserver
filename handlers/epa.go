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
	eids, resStatus, _ := pg.GetEPList()
	// Response
	ResponseWithMsg(c, resStatus, eids)
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

// CheckEPCodeExistHandler 检查执行项目类别档案名称是否存在
func CheckEPCodeExistHandler(c *gin.Context) {
	eid := new(pg.ExecutionProject)
	err := c.ShouldBind(eid)
	if err != nil {
		zap.L().Error("CheckEPCodeExistHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Check Code
	resStatus, _ := eid.CheckCodeExist()
	// Response
	ResponseWithMsg(c, resStatus, eid)
}

// Add Execution Project master data handler
func AddEPHandler(c *gin.Context) {
	// Parse the parameters
	eid := new(pg.ExecutionProject)
	err := c.ShouldBind(eid)
	if err != nil {
		zap.L().Error("AddEPHandler invaid parms", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operator ID
	opeartorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, eid)
		return
	}
	eid.Creator.ID = opeartorID
	// Add EP
	resStatus, _ = eid.Add()
	// Response
	ResponseWithMsg(c, resStatus, eid)
}

// Edit Execution Project master data handler
func EditEPHandler(c *gin.Context) {
	// Parse the parameters
	eid := new(pg.ExecutionProject)
	err := c.ShouldBind(eid)
	if err != nil {
		zap.L().Error("EditEPHandler invaid parms", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operator ID
	opeartorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, eid)
		return
	}
	eid.Modifier.ID = opeartorID

	// Modify
	resStats, _ := eid.Edit()
	// Response
	ResponseWithMsg(c, resStats, eid)
}

// Delete Execution Project master data handler
func DeleteEPHandler(c *gin.Context) {
	// Get the parameters
	eid := new(pg.ExecutionProject)
	err := c.ShouldBind(eid)
	if err != nil {
		zap.L().Error("DeleteEPHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get operator ID
	opeartorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, eid)
		return
	}
	eid.Modifier.ID = opeartorID
	// Delete
	resStatus, _ = eid.Delete()
	// Response
	ResponseWithMsg(c, resStatus, eid)
}

// Batch delete Execution Project master datas handler
func DeleteEPsHandler(c *gin.Context) {
	// Get the parameters
	eids := new([]pg.ExecutionProject)
	err := c.ShouldBind(eids)
	if err != nil {
		zap.L().Error("DeleteEPsHandler invaid parms", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
	}
	// Get operator ID
	opeartorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, eids)
		return
	}
	// Batch delete
	resStatus, _ = pg.DeleteEPs(eids, opeartorID)
	// Response
	ResponseWithMsg(c, resStatus, eids)
}
