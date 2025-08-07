package handlers

import (
	"sccsmsserver/db/pg"
	"sccsmsserver/i18n"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Get Department list handler
func GetDeptsHandler(c *gin.Context) {
	// Get departments list
	depts, resStatus, _ := pg.GetDepts()
	// Response
	ResponseWithMsg(c, resStatus, depts)
}

// Get simplify department list handler
func GetSimpDeptsHandler(c *gin.Context) {
	// Get simpDepts list
	simpDepts, resStatus, _ := pg.GetSimpDepts()
	// Response
	ResponseWithMsg(c, resStatus, simpDepts)
}

// Get simplify department front cache handler
func GetSimpDeptsCacheHandler(c *gin.Context) {
	dc := new(pg.SimpDeptCache)
	err := c.ShouldBind(dc)
	if err != nil {
		zap.L().Error("GetSimpDeptsCacheHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, nil)
		return
	}
	resStatus, _ := dc.GetLatestSimpDepts()

	ResponseWithMsg(c, resStatus, dc)
}

// Check department code exists handler
func CheckDeptCodeExistHandler(c *gin.Context) {
	d := new(pg.Department)
	err := c.ShouldBind(d)
	if err != nil {
		zap.L().Error("CheckDeptCodeExistHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, nil)
		return
	}
	resStatus, _ := d.CheckDeptCodeExist()
	ResponseWithMsg(c, resStatus, d)
}

// Add department
func AddDeptHandler(c *gin.Context) {
	d := new(pg.Department)
	err := c.ShouldBind(d)
	if err != nil {
		zap.L().Error("AddDeptHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, nil)
		return
	}
	// Get Operator ID
	operatorID, resStatus := GetCurrentUser(c)
	if resStatus != i18n.StatusOK {
		zap.L().Error("AddDeptHandler getCurrentUser failed: " + resStatus.String())
		ResponseWithMsg(c, i18n.CodeInternalError, d)
		return
	}
	d.Creator.ID = operatorID
	// Add department
	resStatus, _ = d.AddDept()
	// Resopnse
	ResponseWithMsg(c, resStatus, d)
}

// Edit department
func EditDeptHandler(c *gin.Context) {
	d := new(pg.Department)
	err := c.ShouldBind(d)
	if err != nil {
		zap.L().Error("EditDeptHandler invaid parms", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, nil)
	}
	// Get Operator ID
	operatorID, resStatus := GetCurrentUser(c)
	if resStatus != i18n.StatusOK {
		zap.L().Error("EditDeptHandler getCurrentUser failed: " + resStatus.String())
		ResponseWithMsg(c, i18n.CodeInternalError, d)
		return
	}
	d.Modifier.ID = operatorID

	// edit department
	resStatus, _ = d.Edit()
	// Response
	ResponseWithMsg(c, resStatus, d)
}

// Delete department
func DelDeptHandler(c *gin.Context) {
	d := new(pg.Department)
	err := c.ShouldBind(d)
	if err != nil {
		zap.L().Error("DelDeptHandler invaid parms", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, nil)
	}
	// Get Operator ID
	operatorID, resStatus := GetCurrentUser(c)
	if resStatus != i18n.StatusOK {
		zap.L().Error("DelDeptHandler getCurrentUser failed: " + resStatus.String())
		ResponseWithMsg(c, i18n.CodeInternalError, d)
		return
	}
	d.Modifier.ID = operatorID

	// Delete
	statusCode, _ := d.Delete()
	// Resopnse
	ResponseWithMsg(c, statusCode, d)
}

// Batch delete department
func DelDeptsHandler(c *gin.Context) {
	depts := new([]pg.Department)
	err := c.ShouldBind(depts)
	if err != nil {
		zap.L().Error("DelDeptsHandler invaid parms", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, nil)
	}
	// Get Operator ID
	operatorID, resStatus := GetCurrentUser(c)
	if resStatus != i18n.StatusOK {
		zap.L().Error("DelDeptsHandler getCurrentUser failed: " + resStatus.String())
		ResponseWithMsg(c, i18n.CodeInternalError, depts)
		return
	}
	// Delete departments
	resStatus, _ = pg.DeleteDepts(depts, operatorID)
	// Resopnse
	ResponseWithMsg(c, resStatus, depts)
}
