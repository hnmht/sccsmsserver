package handlers

import (
	"sccsmsserver/db/pg"
	"sccsmsserver/i18n"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Add Issue Resolution Form handler
func AddIRFHandler(c *gin.Context) {
	irf := new(pg.IssueResolutionForm)
	err := c.ShouldBind(irf)
	if err != nil {
		zap.L().Error("AddIRFHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operation ID
	operatorID, resStatus := GetCurrentUser(c)
	if resStatus != i18n.StatusOK {
		zap.L().Error("AddIRFHandler getCurrentUser failed", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInternalError, irf)
		return
	}
	irf.Creator.ID = operatorID
	// Add
	resStatus, _ = irf.Add()
	// Resopnse
	ResponseWithMsg(c, resStatus, irf)
}

// Edit Issue Resolution Form handler
func EditIRFHandler(c *gin.Context) {
	irf := new(pg.IssueResolutionForm)
	err := c.ShouldBind(irf)
	if err != nil {
		zap.L().Error("EditIRFHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operation ID
	operatorID, resStatus := GetCurrentUser(c)
	if resStatus != i18n.StatusOK {
		zap.L().Error("EditIRFHandler getCurrentUser failed", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInternalError, irf)
		return
	}
	irf.Modifier.ID = operatorID
	// Modify
	resStatus, _ = irf.Edit()
	// Resopnse
	ResponseWithMsg(c, resStatus, irf)
}

// Delete Issue Resolution Form handler
func DeleteIRFhandler(c *gin.Context) {
	irf := new(pg.IssueResolutionForm)
	err := c.ShouldBind(irf)
	if err != nil {
		zap.L().Error("DeleteIRFhandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	//Get Operation ID
	operatorID, resStatus := GetCurrentUser(c)
	if resStatus != i18n.StatusOK {
		zap.L().Error("DeleteIRFhandler getCurrentUser failed", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInternalError, irf)
		return
	}
	// Delete
	resStatus, _ = irf.Delete(operatorID)
	// Resopnse
	ResponseWithMsg(c, resStatus, irf)
}

// Confirm Issue Resolution Form handler
func ConfirmIRFhandler(c *gin.Context) {
	irf := new(pg.IssueResolutionForm)
	err := c.ShouldBind(irf)
	if err != nil {
		zap.L().Error("ConfirmIRFhandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	//Get Operation ID
	operationID, resStatus := GetCurrentUser(c)
	if resStatus != i18n.StatusOK {
		zap.L().Error("ConfirmIRFhandler getCurrentUser failed", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInternalError, irf)
		return
	}
	// Confirm
	resStatus, _ = irf.Confirm(operationID)
	// Response
	ResponseWithMsg(c, resStatus, irf)
}

// UnConfirm Issue Resolution Form handler
func UnConfirmIRFhandler(c *gin.Context) {
	irf := new(pg.IssueResolutionForm)
	err := c.ShouldBind(irf)
	if err != nil {
		zap.L().Error("UnConfirmIRFhandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	//Get Operation ID
	operationID, resStatus := GetCurrentUser(c)
	if resStatus != i18n.StatusOK {
		zap.L().Error("UnConfirmIRFhandler getCurrentUser failed", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInternalError, irf)
		return
	}
	// UnConfirm
	resStatus, _ = irf.UnConfirm(operationID)
	// Response
	ResponseWithMsg(c, resStatus, irf)
}

// Get Issue Resolution Form list handler
func GetIRFListHandler(c *gin.Context) {
	qp := new(pg.QueryParams)
	err := c.ShouldBind(qp)
	if err != nil {
		zap.L().Error("GetIRFListHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get List
	irfs, resStauts, _ := pg.GetIRFList(qp.QueryString)
	// Response
	ResponseWithMsg(c, resStauts, irfs)
}
