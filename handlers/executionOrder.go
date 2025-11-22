package handlers

import (
	"sccsmsserver/db/pg"
	"sccsmsserver/i18n"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Get the list of  Execution Orders to be referenced handler
func GetReferEOHandler(c *gin.Context) {
	qp := new(pg.QueryParams)
	err := c.ShouldBind(qp)
	if err != nil {
		zap.L().Error("GetReferEOHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get list
	reos, resStatus, _ := pg.GetReferEOs(qp.QueryString)
	// Response
	ResponseWithMsg(c, resStatus, reos)
}

// Add Execution Order handler
func AddEOHandler(c *gin.Context) {
	eo := new(pg.ExecutionOrder)
	err := c.ShouldBind(eo)
	if err != nil {
		zap.L().Error("AddEOHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operator ID
	operatorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, eo)
		return
	}
	eo.Creator.ID = operatorID
	// Add
	resStatus, _ = eo.Add()
	// Resoponse
	ResponseWithMsg(c, resStatus, eo)
}

// Edit Execution Order handler
func EditEOHandler(c *gin.Context) {
	eo := new(pg.ExecutionOrder)
	err := c.ShouldBind(eo)
	if err != nil {
		zap.L().Error("EditEOHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operator ID
	operatorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, eo)
		return
	}
	eo.Modifier.ID = operatorID
	// Modify
	resStatus, _ = eo.Edit()
	// Response
	ResponseWithMsg(c, resStatus, eo)
}

// Delete Execution Order handler
func DeleteEOHandler(c *gin.Context) {
	eo := new(pg.ExecutionOrder)
	err := c.ShouldBind(eo)
	if err != nil {
		zap.L().Error("DeleteEOHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operator ID
	operatorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, eo)
		return
	}
	// Delete
	resStatus, _ = eo.Delete(operatorID)
	// Response
	ResponseWithMsg(c, resStatus, eo)
}

// Confirm Execution Order handler
func ConfirmEOHandler(c *gin.Context) {
	eo := new(pg.ExecutionOrder)
	err := c.ShouldBind(eo)
	if err != nil {
		zap.L().Error("ConfirmEOHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operator ID
	operatorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, eo)
		return
	}
	// Confirm
	resStatus, _ = eo.Confirm(operatorID)
	// Response
	ResponseWithMsg(c, resStatus, eo)
}

// UnConfirm Execution Order handler
func CancelConfirmEOHandler(c *gin.Context) {
	eo := new(pg.ExecutionOrder)
	err := c.ShouldBind(eo)
	if err != nil {
		zap.L().Error("CancelConfirmEOHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operator ID
	operatorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, eo)
		return
	}
	// Un-Confirm
	resStatus, _ = eo.UnConfirm(operatorID)
	// Response
	ResponseWithMsg(c, resStatus, eo)
}

// Get the list of Execution Orders handler
func GetEOListHandler(c *gin.Context) {
	qp := new(pg.QueryParams)
	err := c.ShouldBind(qp)
	if err != nil {
		zap.L().Error("GetEOListHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get list
	eos, resStatus, _ := pg.GetEOList(qp.QueryString)
	// Response
	ResponseWithMsg(c, resStatus, eos)
}

// Get Execution Order List for Review handler
func GetEOReviewListHandler(c *gin.Context) {
	qp := new(pg.QueryParams)
	err := c.ShouldBind(qp)
	if err != nil {
		zap.L().Error("GetEOReviewListHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operator ID
	operatorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, qp)
		return
	}
	// Get List
	eos, resStatus, _ := pg.GetEOReviewList(qp.QueryString, operatorID)
	// Response
	ResponseWithMsg(c, resStatus, eos)
}

// Get the list of Execution Orders to be reviewed by pagination handler
func GetEOReviewListPaginationHandler(c *gin.Context) {
	pqp := new(pg.PagingQueryParams)
	err := c.ShouldBind(pqp)
	if err != nil {
		zap.L().Error("GetEOReviewListPaginationHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operator ID
	operatorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, pqp)
		return
	}
	// Get list
	eosp, resStatus, _ := pg.GetEOReviewListPagination(*pqp, operatorID)
	// Response
	ResponseWithMsg(c, resStatus, eosp)
}

// Get Execution Order details by HID handler
func GetEOInfoByHIDHandler(c *gin.Context) {
	eo := new(pg.ExecutionOrder)
	err := c.ShouldBind(eo)
	if err != nil {
		zap.L().Error("GetEOInfoByHIDHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Details
	resStatus, _ := eo.GetDetailByHID()
	// Resopnse
	ResponseWithMsg(c, resStatus, eo)
}

// Add Execution Order Comment handler
func AddCommentHandler(c *gin.Context) {
	eoc := new(pg.ExecutionOrderComment)
	err := c.ShouldBind(eoc)
	if err != nil {
		zap.L().Error("AddCommentHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operator ID
	operatorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, eoc)
		return
	}
	eoc.Creator.ID = operatorID
	// Add
	resStatus, _ = eoc.Add()
	// Response
	ResponseWithMsg(c, resStatus, eoc)
}

// Add Execution Order Review Record handler
func AddReviewHandler(c *gin.Context) {
	eor := new(pg.ExecutionOrderReview)
	err := c.ShouldBind(eor)
	if err != nil {
		zap.L().Error("AddReviewHandler invalid param", zap.Error(err))

		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operator ID
	operatorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, eor)
		return
	}
	eor.Creator.ID = operatorID
	// Add
	resStatus, _ = eor.Add()
	// Resopnse
	ResponseWithMsg(c, resStatus, eor)
}

// Get Execution Order Comments handler
func GetEOCommentsHandler(c *gin.Context) {
	cs := new(pg.EOCommentsParams)
	err := c.ShouldBind(cs)
	if err != nil {
		zap.L().Error("GetEOCommentsHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get comments
	resStatus, _ := cs.Get()
	// Responses
	ResponseWithMsg(c, resStatus, cs)
}

// Get Execution Order Review Records handler
func GetEOReviewsHandler(c *gin.Context) {
	rs := new(pg.EOReviewsParams)
	err := c.ShouldBind(rs)
	if err != nil {
		zap.L().Error("GetEOReviewsHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get records
	resStatus, _ := rs.Get()
	// Response
	ResponseWithMsg(c, resStatus, rs)
}
