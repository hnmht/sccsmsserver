package handlers

import (
	"sccsmsserver/db/pg"
	"sccsmsserver/i18n"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Get Personal Protective Equipment Quota List handler
func GetPPEQuotaListHandler(c *gin.Context) {
	qp := new(pg.QueryParams)
	err := c.ShouldBind(qp)
	if err != nil {
		zap.L().Error("GetPPEQuotaListHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get PPEQuota List
	pqs, resStatus, _ := pg.GetPQList(qp.QueryString)
	// Response
	ResponseWithMsg(c, resStatus, pqs)
}

// Get Personal Protective Qtuipment Quota Details by HID handler
func GetPPEQuotaInfoByHIDHandler(c *gin.Context) {
	pq := new(pg.PPEQuota)
	err := c.ShouldBind(pq)
	if err != nil {
		zap.L().Error("GetPPEQuotaInfoByHIDHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Details
	resStatus, _ := pq.GetDetailByHID()
	// Response
	ResponseWithMsg(c, resStatus, pq)
}

// Add Personal Protective Equipment Quota Handler
func AddPPEQuotaHandler(c *gin.Context) {
	pq := new(pg.PPEQuota)
	err := c.ShouldBind(pq)
	if err != nil {
		zap.L().Error("AddPPEQuotaHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operator ID
	operatorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, pq)
		return
	}
	pq.Creator.ID = operatorID
	// Add
	resStatus, _ = pq.Add()
	// Response
	ResponseWithMsg(c, resStatus, pq)
}

// Edit Personal Protective Equipment Quota handler
func EditPPEQuotaHandler(c *gin.Context) {
	pq := new(pg.PPEQuota)
	err := c.ShouldBind(pq)
	if err != nil {
		zap.L().Error("EditPPEQuotaHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operator ID
	operatorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, pq)
		return
	}
	pq.Modifier.ID = operatorID
	// Modify
	resStatus, _ = pq.Edit()
	// Response
	ResponseWithMsg(c, resStatus, pq)
}

// Delete Personal Protective Equipment Quota Handler
func DeletePPEQuotaHandler(c *gin.Context) {
	pq := new(pg.PPEQuota)
	err := c.ShouldBind(pq)
	if err != nil {
		zap.L().Error("DeletePPEQuotaHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operator ID
	operatorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, pq)
		return
	}
	// Delete
	resStatus, _ = pq.Delete(operatorID)
	// Response
	ResponseWithMsg(c, resStatus, pq)
}

// Confirm Personal Protective Equipment Quota Handler
func ConfirmPPEQuotaHandler(c *gin.Context) {
	pq := new(pg.PPEQuota)
	err := c.ShouldBind(pq)
	if err != nil {
		zap.L().Error("ConfirmPPEQuotaHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operator ID
	operatorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, pq)
		return
	}
	// Confirm
	resStatus, _ = pq.Confirm(operatorID)
	// Response
	ResponseWithMsg(c, resStatus, pq)
}

// Unconfirm Personal Protective Equipment Quota handler
func UnconfirmPPEQuotaHandler(c *gin.Context) {
	pq := new(pg.PPEQuota)
	err := c.ShouldBind(pq)
	if err != nil {
		zap.L().Error("UnconfirmPPEQuotaHandler invalid param", zap.Error(err))

		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operator ID
	operatorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, pq)
		return
	}
	// Unconfirm
	resStatus, _ = pq.Unconfirm(operatorID)
	// Response
	ResponseWithMsg(c, resStatus, pq)
}

// Check if a PPE Position Quota for the same period handler
func CheckPPEQuotaExistHandler(c *gin.Context) {
	pq := new(pg.PPEQuota)
	err := c.ShouldBind(pq)
	if err != nil {
		zap.L().Error("CheckOPNameExistHandler invalid param", zap.Error(err))

		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Check
	resStatus, _ := pq.CheckExist()
	// Response
	ResponseWithMsg(c, resStatus, pq)
}

// Get the list of all position that have PPE Quotas within the same period handler
func GetPPEPositionsPeriodHandler(c *gin.Context) {
	pps := new(pg.PPEPositionsParams)
	err := c.ShouldBind(pps)
	if err != nil {
		zap.L().Error("GetOpsHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Retrieve
	resStatus, _ := pps.Get()
	// Response
	ResponseWithMsg(c, resStatus, pps)
}
