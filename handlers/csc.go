package handlers

import (
	"sccsmsserver/db/pg"
	"sccsmsserver/i18n"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Get CSC list handler
func GetCSCListHandler(c *gin.Context) {
	cscs, resStatus, _ := pg.GetCSCList()
	ResponseWithMsg(c, resStatus, cscs)
}

// Get Simple CSC list handler
func GetSimpCSCListHandler(c *gin.Context) {
	ssics, resStatus, _ := pg.GetSimpCSCList()
	ResponseWithMsg(c, resStatus, ssics)
}

// Get Simple CSC front-end cache handler
func GetSimpCSCCacheHandler(c *gin.Context) {
	scscc := new(pg.SimpSICCache)
	err := c.ShouldBind(scscc)
	if err != nil {
		zap.L().Error("GetSimpCSCCacheHandler invalid param:", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, nil)
		return
	}

	// Get latest Simple CSC for front-end cache
	resStatus, _ := scscc.GetSimpCSCCache()
	ResponseWithMsg(c, resStatus, scscc)
}

// Check if the CSC name exists handler.
func CheckCSCNameExistHandler(c *gin.Context) {
	csc := new(pg.CSC)
	err := c.ShouldBind(csc)
	if err != nil {
		zap.L().Error("CheckCSCNameExistHandler invalid param:", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, nil)
		return
	}
	// Check
	resStatus, _ := csc.CheckNameExist()
	ResponseWithMsg(c, resStatus, csc)
}

// Add CSC handler
func AddCSCHandler(c *gin.Context) {
	csc := new(pg.CSC)
	err := c.ShouldBind(csc)
	if err != nil {
		zap.L().Error("AddCSCHandler invalid param:", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get current operator ID
	operatorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, csc)
		return
	}
	csc.Creator.ID = operatorID
	// Add
	resStatus, _ = csc.Add()
	// Response
	ResponseWithMsg(c, resStatus, csc)
}

// Edit CSC hanlder
func EditCSCHandler(c *gin.Context) {
	csc := new(pg.CSC)
	err := c.ShouldBind(csc)
	if err != nil {
		zap.L().Error("EditCSCHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get current operator ID
	operatorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, csc)
		return
	}
	csc.Modifier.ID = operatorID

	// Edit
	resStatus, _ = csc.Edit()
	ResponseWithMsg(c, resStatus, csc)
}

// Delete CSC handler
func DeleteCSCHandler(c *gin.Context) {
	csc := new(pg.CSC)
	err := c.ShouldBind(csc)
	if err != nil {
		zap.L().Error("DeleteCSCHandler invalid param:", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get current opeartor ID
	operatorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, csc)
		return
	}
	csc.Modifier.ID = operatorID
	// Delete
	resStatus, _ = csc.Delete()
	// Response
	ResponseWithMsg(c, resStatus, csc)
}

// Batch delete CSCs handler
func DeleteCSCsHandler(c *gin.Context) {
	cscs := new([]pg.CSC)
	err := c.ShouldBind(cscs)
	if err != nil {
		zap.L().Error("DeleteCSCsHandler invaid parms", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
	}
	// Get current operator ID
	operatorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, cscs)
		return
	}
	// Batch delete
	resStatus, _ = pg.DeleteCSCs(cscs, operatorID)
	// Response
	ResponseWithMsg(c, resStatus, cscs)
}
