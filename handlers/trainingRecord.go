package handlers

import (
	"sccsmsserver/db/pg"
	"sccsmsserver/i18n"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Add Training Record handler
func AddTRHandler(c *gin.Context) {
	tr := new(pg.TrainingRecord)
	err := c.ShouldBind(tr)
	if err != nil {
		zap.L().Error("AddTRHandler invalid params:", zap.Error(err))

		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	//Get Operator ID
	operatorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, tr)
		return
	}
	tr.Creator.ID = operatorID
	// Add
	resStatus, _ = tr.Add()
	// Response
	ResponseWithMsg(c, resStatus, tr)
}

// Get Training Record list handler
func GetTRListHandler(c *gin.Context) {
	qp := new(pg.QueryParams)
	err := c.ShouldBind(qp)
	if err != nil {
		zap.L().Error("GetTRListHandler invalid params:", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get List
	trs, resStatus, _ := pg.GetTRList(qp.QueryString)
	// Response
	ResponseWithMsg(c, resStatus, trs)
}

// Get Training Recorde details by HID
func GetTRInfoByHIDHandler(c *gin.Context) {
	tr := new(pg.TrainingRecord)
	err := c.ShouldBind(tr)
	if err != nil {
		zap.L().Error("GetTRInfoByHIDHandler invalid params:", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Training Record List
	resStatus, _ := tr.GetDetailByHID()
	// Response
	ResponseWithMsg(c, resStatus, tr)
}

// Edit Training Recrod handler
func EditTRHandler(c *gin.Context) {
	tr := new(pg.TrainingRecord)
	err := c.ShouldBind(tr)
	if err != nil {
		zap.L().Error("EditTRHandler invalid params:", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	//Get Operator ID
	operatorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, tr)
		return
	}
	tr.Modifier.ID = operatorID
	// Modify
	resStatus, _ = tr.Edit()
	// Response
	ResponseWithMsg(c, resStatus, tr)
}

// Delete Training Record handler
func DeleteTRHandler(c *gin.Context) {
	tr := new(pg.TrainingRecord)
	err := c.ShouldBind(tr)
	if err != nil {
		zap.L().Error("DeleteTRHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	//Get Operator ID
	operatorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, tr)
		return
	}
	// Delete
	resStatus, _ = tr.Delete(operatorID)
	// Response
	ResponseWithMsg(c, resStatus, tr)
}

// Confirm Training Record handler
func ConfirmTRHandler(c *gin.Context) {
	tr := new(pg.TrainingRecord)
	err := c.ShouldBind(tr)
	if err != nil {
		zap.L().Error("ConfirmTRHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	//Get Operator ID
	operatorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, tr)
		return
	}
	// Confirm
	resStatus, _ = tr.Confirm(operatorID)
	// Response
	ResponseWithMsg(c, resStatus, tr)
}

// UnConfirm Training Record handler
func UnConfirmTRHandler(c *gin.Context) {
	tr := new(pg.TrainingRecord)
	err := c.ShouldBind(tr)
	if err != nil {
		zap.L().Error("UnConfirmTRHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operator ID
	operatorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, tr)
		return
	}
	// UnConfirm
	resStatus, _ = tr.UnConfirm(operatorID)
	// Response
	ResponseWithMsg(c, resStatus, tr)
}

// Get Taught Lessons Report handler
func GetTaughtLessonsReportHandler(c *gin.Context) {
	qp := new(pg.QueryParams)
	err := c.ShouldBind(qp)
	if err != nil {
		zap.L().Error("GetTaughtLessonsReportHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get report
	glrs, resStatus, _ := pg.GetTaughtLessonsReport(qp.QueryString)
	// Response
	ResponseWithMsg(c, resStatus, glrs)
}

// Get Recieved Training Report handler
func GetRecivedTrainingReportHandler(c *gin.Context) {
	qp := new(pg.QueryParams)
	err := c.ShouldBind(qp)
	if err != nil {
		zap.L().Error("GetRecivedTrainingReportHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Report
	rtrs, resStatus, _ := pg.GetRecivedTrainingReport(qp.QueryString)
	// Response
	ResponseWithMsg(c, resStatus, rtrs)
}
