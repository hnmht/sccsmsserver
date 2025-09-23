package handlers

import (
	"sccsmsserver/db/pg"
	"sccsmsserver/i18n"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Add a new execution project template handler
func AddEPTHandler(c *gin.Context) {
	// Get parameters
	ept := new(pg.EPT)
	err := c.ShouldBind(ept)
	if err != nil {
		zap.L().Error("AddEPTHandler invalid params", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
	}
	// Get the operator id
	operatorID, resStatus := GetCurrentUser(c)
	if resStatus != i18n.StatusOK {
		zap.L().Error("AddEPTHandler getCurrentUser failed", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInternalError, ept)
		return
	}
	ept.Creator.ID = operatorID
	// Add
	resStatus, _ = ept.Add()
	// Response
	ResponseWithMsg(c, resStatus, ept)
}

// Edit Execution Project Template Handler
func EditEPTHandler(c *gin.Context) {
	// Get parameters
	ept := new(pg.EPT)
	err := c.ShouldBind(ept)
	if err != nil {
		zap.L().Error("EditEPTHandler invalid params", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
	}
	// Get the operator id
	operatorID, resStatus := GetCurrentUser(c)
	if resStatus != i18n.StatusOK {
		zap.L().Error("EditEPTHandler getCurrentUser failed", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInternalError, ept)
		return
	}
	ept.Modifier.ID = operatorID
	// Edit
	resStatus, _ = ept.Edit()
	// Response
	ResponseWithMsg(c, resStatus, ept)
}

// Delete Execution Project Template Handler
func DeleteEPTHandler(c *gin.Context) {
	// Get parameters
	ept := new(pg.EPT)
	err := c.ShouldBind(ept)
	if err != nil {
		zap.L().Error("DeleteEPTHandler invalid params", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
	}
	// Get the operator id
	operatorID, resStatus := GetCurrentUser(c)
	if resStatus != i18n.StatusOK {
		zap.L().Error("DeleteEPTHandler getCurrentUser failed", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInternalError, ept)
		return
	}
	ept.Modifier.ID = operatorID
	// Delete
	resStatus, _ = ept.Delete()
	// Response
	ResponseWithMsg(c, resStatus, ept)
}

// Check if the execution project template code exists
func CheckEPTCodeExistHandler(c *gin.Context) {
	// Get parameters
	ept := new(pg.EPT)
	err := c.ShouldBind(ept)
	if err != nil {
		zap.L().Error("CheckEPTCodeExistHandler invalid params", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
	}
	// Check if the code exists
	resStatus, _ := ept.CheckCodeExist()
	// Response
	ResponseWithMsg(c, resStatus, ept)
}

// Get Execution Project Template List Handler
func GetEPTListHandler(c *gin.Context) {
	// Get List
	epts, resStatus, _ := pg.GetEPTList()
	// Response
	ResponseWithMsg(c, resStatus, epts)
}

// Get Execution Project Template front-end cache handler
func GetEPTCacheHandler(c *gin.Context) {
	// Get parameters
	eptc := new(pg.EPTCache)
	err := c.ShouldBind(eptc)
	if err != nil {
		zap.L().Error("GetEPTCacheHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get latest Execution Project Template for front-end cache
	resStatus, _ := eptc.GetEPTCahce()
	// Response
	ResponseWithMsg(c, resStatus, eptc)
}

// Delete multiple Execution Project Templates Handler
func DeleteEPTsHandler(c *gin.Context) {
	// Get parameters
	epts := new([]pg.EPT)
	err := c.ShouldBind(epts)
	if err != nil {
		zap.L().Error("DeleteEPTsHandler invalid params", zap.Error(err))

		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
	}
	// Get the operator id
	operatorID, resStatus := GetCurrentUser(c)
	if resStatus != i18n.StatusOK {
		zap.L().Error("DeleteEPTsHandler getCurrentUser failed", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInternalError, epts)
		return
	}
	// Delete
	resStatus, _ = pg.DeleteEPTs(epts, operatorID)
	// Response
	ResponseWithMsg(c, resStatus, epts)
}
