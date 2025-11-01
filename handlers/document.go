package handlers

import (
	"sccsmsserver/db/pg"
	"sccsmsserver/i18n"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Add Document handler
func AddDocumentHandler(c *gin.Context) {
	d := new(pg.Document)
	err := c.ShouldBind(d)
	if err != nil {
		zap.L().Error("AddDocumentHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operator ID
	operatorID, resStatus := GetCurrentUser(c)
	if resStatus != i18n.StatusOK {
		zap.L().Error("AddDocumentHandler getCurrentUser failed", zap.Error(err))
		ResponseWithMsg(c, resStatus, d)
		return
	}
	d.Creator.ID = operatorID
	// Add
	resStatus, _ = d.Add()
	// Response
	ResponseWithMsg(c, resStatus, d)
}

// Edit Document Handler
func EditDocumentHandler(c *gin.Context) {
	d := new(pg.Document)
	err := c.ShouldBind(d)
	if err != nil {
		zap.L().Error("EditDocumentHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operator ID
	operaorID, resStatus := GetCurrentUser(c)
	if resStatus != i18n.StatusOK {
		zap.L().Error("EditDocumentHandler getCurrentUser failed", zap.Error(err))
		ResponseWithMsg(c, resStatus, d)
		return
	}
	d.Modifier.ID = operaorID
	// Modify
	resStatus, _ = d.Edit()
	// Response
	ResponseWithMsg(c, resStatus, d)
}

// Get Document pagination list handler
func GetDocumentPagingListHanlder(c *gin.Context) {
	dpp := new(pg.DCPagingParams)
	err := c.ShouldBind(dpp)
	if err != nil {
		zap.L().Error("GetDocumentPagingListHanlder invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get document list
	resStatus, _ := dpp.Get()
	// Response
	ResponseWithMsg(c, resStatus, dpp)
}

// Delete Document Handler
func DeleteDocumentHandler(c *gin.Context) {
	d := new(pg.Document)
	err := c.ShouldBind(d)
	if err != nil {
		zap.L().Error("DeleteDocumentHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operator ID
	operaorID, resStatus := GetCurrentUser(c)
	if resStatus != i18n.StatusOK {
		zap.L().Error("DeleteDocumentHandler getCurrentUser failed", zap.Error(err))
		ResponseWithMsg(c, resStatus, d)
		return
	}
	// Delete
	resStatus, _ = d.Delete(operaorID)
	// Response
	ResponseWithMsg(c, resStatus, d)
}

// Batch Delete Document handler
func DeleteDocumentsHandler(c *gin.Context) {
	docs := new([]pg.Document)
	err := c.ShouldBind(docs)
	if err != nil {
		zap.L().Error("DeleteDocumentsHandler invaid parms", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
	}
	// Get Operator ID
	operaorID, resStatus := GetCurrentUser(c)
	if resStatus != i18n.StatusOK {
		zap.L().Error("DeleteDocumentsHandler getCurrentUser failed", zap.Error(err))
		ResponseWithMsg(c, resStatus, docs)
		return
	}
	// Batch Delete
	resStatus, _ = pg.DeleteDocuments(docs, operaorID)
	// Response
	ResponseWithMsg(c, resStatus, docs)
}

// Get Document Report handler
func GetDocumentReportHandler(c *gin.Context) {
	qp := new(pg.QueryParams)
	err := c.ShouldBind(qp)
	if err != nil {
		zap.L().Error("GetDocumentReportHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Report
	rep, resStatus, _ := pg.GetQueryDocumentReport(qp.QueryString)
	// Response
	ResponseWithMsg(c, resStatus, rep)
}
