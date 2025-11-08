package handlers

import (
	"sccsmsserver/db/pg"
	"sccsmsserver/i18n"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Get the list of Work Order awaiting execution handler
func GetWOReferHandler(c *gin.Context) {
	qp := new(pg.QueryParams)
	err := c.ShouldBind(qp)
	if err != nil {
		zap.L().Error("GetWOReferHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get list
	wors, resStatus, _ := pg.GetWORefer(qp.QueryString)
	// Response
	ResponseWithMsg(c, resStatus, wors)
}

// Get Work Order list handler
func GetWOListHanlder(c *gin.Context) {
	qp := new(pg.QueryParams)
	err := c.ShouldBind(qp)
	if err != nil {
		zap.L().Error("GetWOListHanlder invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get list
	wos, resStatus, _ := pg.GetWOList(qp.QueryString)
	// Response
	ResponseWithMsg(c, resStatus, wos)
}

// Get Work Order details handler
func GetWOInfoByIDHandler(c *gin.Context) {
	wo := new(pg.WorkOrder)
	err := c.ShouldBind(wo)
	if err != nil {
		zap.L().Error("GetWOInfoByIDHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Details
	resStatus, _ := wo.GetDetailByHID()
	// Resoponse
	ResponseWithMsg(c, resStatus, wo)
}

// Add Work Order handler
func AddWOHandler(c *gin.Context) {
	wo := new(pg.WorkOrder)
	err := c.ShouldBind(wo)
	if err != nil {
		zap.L().Error("AddWOHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operator ID
	operatorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, wo)
		return
	}
	wo.Creator.ID = operatorID
	// Add
	resStatus, _ = wo.Add()
	// Response
	ResponseWithMsg(c, resStatus, wo)
}

// Edit Work Order handler
func EditWOHandler(c *gin.Context) {
	wo := new(pg.WorkOrder)
	err := c.ShouldBind(wo)
	if err != nil {
		zap.L().Error("EditWOHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operator ID
	operatorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, wo)
		return
	}
	wo.Modifier.ID = operatorID
	// Modify
	resStatus, _ = wo.Edit()
	// Response
	ResponseWithMsg(c, resStatus, wo)
}

// Delete Work Order handler
func DeleteWOHandler(c *gin.Context) {
	wo := new(pg.WorkOrder)
	err := c.ShouldBind(wo)
	if err != nil {
		zap.L().Error("DeleteWOHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operator ID
	operatorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, wo)
		return
	}
	// Delete
	resStatus, _ = wo.Delete(operatorID)
	// Resoponse
	ResponseWithMsg(c, resStatus, wo)
}

// Batch delete Work Order handler
func DeleteWOsHandler(c *gin.Context) {
	wos := new([]pg.WorkOrder)
	err := c.ShouldBind(wos)
	if err != nil {
		zap.L().Error("DeleteWOsHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operator ID
	operatorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, wos)
		return
	}
	// Batch Delete
	resStatus, _ = pg.DeleteWOs(wos, operatorID)
	ResponseWithMsg(c, resStatus, wos)
}

// Confirm Work Order handler
func ConfirmWOHandler(c *gin.Context) {
	wo := new(pg.WorkOrder)
	err := c.ShouldBind(wo)
	if err != nil {
		zap.L().Error("ConfirmWOHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operator ID
	operatorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, wo)
		return
	}
	// Confirm
	resStatus, _ = wo.Confirm(operatorID)
	// Response
	ResponseWithMsg(c, resStatus, wo)
}

// Unconfirm Work Order handler
func UnConfirmWOHandler(c *gin.Context) {
	wo := new(pg.WorkOrder)
	err := c.ShouldBind(wo)
	if err != nil {
		zap.L().Error("UnConfirmWOHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operator ID
	operatorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, wo)
		return
	}
	// UnConfirm
	resStatus, _ = wo.UnConfirm(operatorID)
	// Response
	ResponseWithMsg(c, resStatus, wo)
}
