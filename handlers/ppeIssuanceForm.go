package handlers

import (
	"sccsmsserver/db/pg"
	"sccsmsserver/i18n"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Add PPE Issuance Form handler
func AddPPEIFHandler(c *gin.Context) {
	pif := new(pg.PPEIssuanceForm)
	err := c.ShouldBind(pif)
	if err != nil {
		zap.L().Error("AddPPEIFHandler invalid params:", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	//Get Current Operator ID
	operatorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, pif)
		return
	}
	pif.Creator.ID = operatorID
	// Add
	resStatus, _ = pif.Add()
	// Response
	ResponseWithMsg(c, resStatus, pif)
}

// Wiard Add PPE Issuance Form handler
func WiardAddPPEIFHandler(c *gin.Context) {
	lidw := new(pg.PPEIssuanceFormWizard)
	err := c.ShouldBind(lidw)
	if err != nil {
		zap.L().Error("WiardAddPPEIFHandler invalid params:", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Current Operator ID
	operatorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, lidw)
		return
	}
	lidw.Params.Creator.ID = operatorID

	// Generate
	resStatus, _ = lidw.Generate()
	// Response
	ResponseWithMsg(c, resStatus, lidw)
}

// Get PPE Issuance Form List handler
func GetPPEIFListHandler(c *gin.Context) {
	qp := new(pg.QueryParams)
	err := c.ShouldBind(qp)
	if err != nil {
		zap.L().Error("GetPPEIFListHandler invalid params:", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get List
	pls, resStatus, _ := pg.GetPPEIFList(qp.QueryString)
	ResponseWithMsg(c, resStatus, pls)
}

// Get PPE Issuance Form Info by HID handler
func GetPPEIFInfoByHIDHandler(c *gin.Context) {
	pif := new(pg.PPEIssuanceForm)
	err := c.ShouldBind(pif)
	if err != nil {
		zap.L().Error("GetPPEIFInfoByHIDHandler invalid params:", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Detail
	resStatus, _ := pif.GetDetailByHID()
	// Response
	ResponseWithMsg(c, resStatus, pif)
}

// Edit PPE Issuance Form handler
func EditPPEIFHandler(c *gin.Context) {
	pif := new(pg.PPEIssuanceForm)
	err := c.ShouldBind(pif)
	if err != nil {
		zap.L().Error("EditPPEIFHandler invalid params:", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	//Get Current Operator ID
	operatorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, pif)
		return
	}
	pif.Modifier.ID = operatorID
	// Modify
	resStatus, _ = pif.Edit()
	// Response
	ResponseWithMsg(c, resStatus, pif)
}

// Delete PPE Issuance Form handler
func DeletePPEIFHandler(c *gin.Context) {
	pif := new(pg.PPEIssuanceForm)
	err := c.ShouldBind(pif)
	if err != nil {
		zap.L().Error("DeletePPEIFHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	//Get Current Operator ID
	operatorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, pif)
		return
	}
	// Delete
	resStatus, _ = pif.Delete(operatorID)
	// Response
	ResponseWithMsg(c, resStatus, pif)
}

// Confirm PPE Issuance Form handler
func ConfirmPPEIFHandler(c *gin.Context) {
	pif := new(pg.PPEIssuanceForm)
	err := c.ShouldBind(pif)
	if err != nil {
		zap.L().Error("ConfirmPPEIFHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	//Get Current Operator ID
	confirmUserID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, pif)
		return
	}
	// Confirm
	resStatus, _ = pif.Confirm(confirmUserID)
	// Response
	ResponseWithMsg(c, resStatus, pif)
}

// Unconfirm PPE Issuance Form handler
func UnconfirmPPEIFHandler(c *gin.Context) {
	pif := new(pg.PPEIssuanceForm)
	err := c.ShouldBind(pif)
	if err != nil {
		zap.L().Error("UnconfirmPPEIFHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Current Operator ID
	confirmUserID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, pif)
		return
	}
	// Unconfirm
	resStatus, _ = pif.Unconfirm(confirmUserID)
	// Response
	ResponseWithMsg(c, resStatus, pif)
}

// Get PPE Issuance Form handler
func GetPPEIFReportHandler(c *gin.Context) {
	qp := new(pg.QueryParams)
	err := c.ShouldBind(qp)
	if err != nil {
		zap.L().Error("GetPPEIFReportHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Report
	ldrs, resStatus, _ := pg.GetPPEIFReport(qp.QueryString)
	// Response
	ResponseWithMsg(c, resStatus, ldrs)
}
