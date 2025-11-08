package handlers

import (
	"sccsmsserver/db/pg"
	"sccsmsserver/i18n"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Get User UnRead Comments handler
func GetUserUnReadCommentsHandler(c *gin.Context) {
	// Get Operator ID
	opeartorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, i18n.CodeInternalError, 0)
		return
	}
	// Get UnRead Comments list
	comments, resStatus, _ := pg.GetUserUnReadComments(opeartorID)
	// Response
	ResponseWithMsg(c, resStatus, comments)
}

// Get User Read Comments handler
func GetUserReadCommentsHandler(c *gin.Context) {
	qp := new(pg.QueryParams)
	err := c.ShouldBind(qp)
	if err != nil {
		zap.L().Error("GetUserReadCommentsHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operator ID
	opeartorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, i18n.CodeInternalError, 0)
		return
	}
	// Get Comments
	comments, resStatus, _ := pg.GetUserReadComments(opeartorID, qp.QueryString)
	// Response
	ResponseWithMsg(c, resStatus, comments)
}

// Get User Work Orders awaiting execution
func GetUserWORefsHandler(c *gin.Context) {
	//Get Operator ID
	opeartorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, i18n.CodeInternalError, 0)
		return
	}
	// Get Work Orders
	wors, resStatus, _ := pg.GetUserWORefs(opeartorID)
	// Response
	ResponseWithMsg(c, resStatus, wors)
}

// Get User To-Do Issues handler
func GetUserEORefsHandler(c *gin.Context) {
	//Get Operator ID
	opeartorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, i18n.CodeInternalError, 0)
		return
	}
	// Get Execution Order list
	eors, resStatus, _ := pg.GetUserEORefs(opeartorID)
	// Response
	ResponseWithMsg(c, resStatus, eors)
}

// Read Comment handler
func ReadCommentMessageHandler(c *gin.Context) {
	cm := new(pg.CommentMessage)
	err := c.ShouldBind(cm)
	if err != nil {
		zap.L().Error("ReadCommentMessageHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	//Get Operator ID
	operatorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, i18n.CodeInternalError, 0)
		return
	}
	cm.Modifier.ID = operatorID
	// Read
	resStatus, _ = cm.Read()
	// Response
	ResponseWithMsg(c, resStatus, cm)
}
