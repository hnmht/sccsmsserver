package handlers

import (
	"sccsmsserver/db/pg"
	"sccsmsserver/i18n"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Get Work Order status report handler
func GetWoReportHandler(c *gin.Context) {
	qp := new(pg.QueryParams)
	err := c.ShouldBind(qp)
	if err != nil {
		zap.L().Error("GetWoReportHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Report
	wors, resStatus, _ := pg.GetWorkOrderReport(qp.QueryString)
	// Response
	ResponseWithMsg(c, resStatus, wors)
}

// Get Execution Order Status Report handler
func GetEoReportHandler(c *gin.Context) {
	qp := new(pg.QueryParams)
	err := c.ShouldBind(qp)
	if err != nil {
		zap.L().Error("GetEoReportHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get report
	edrs, resStatus, _ := pg.GetExecutionOrderReport(qp.QueryString)
	// Response
	ResponseWithMsg(c, resStatus, edrs)
}

// Get Issue Resolution Form Report handler
func GetIRFReportHandler(c *gin.Context) {
	qp := new(pg.QueryParams)
	err := c.ShouldBind(qp)
	if err != nil {
		zap.L().Error("GetIRFReportHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Report
	ddrs, resStatus, _ := pg.GetIssueResolutionFormReport(qp.QueryString)
	// Response
	ResponseWithMsg(c, resStatus, ddrs)
}
