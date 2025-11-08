package handlers

import (
	"sccsmsserver/db/pg"
	"sccsmsserver/i18n"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Add User-defined Archive handler
func AddUDAHandler(c *gin.Context) {
	uda := new(pg.UserDefinedArchive)
	err := c.ShouldBind(uda)
	if err != nil {
		zap.L().Error("AddUDAHandler invalid param", zap.Error(err))

		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operator ID
	creatorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, uda)
		return
	}
	uda.Creator.ID = creatorID
	// Add
	resStatus, _ = uda.Add()
	ResponseWithMsg(c, resStatus, uda)
}

// Edit User-defined Archive handler
func EditUDAHandler(c *gin.Context) {
	uda := new(pg.UserDefinedArchive)
	err := c.ShouldBind(uda)
	if err != nil {
		zap.L().Error("EditUDAHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
	}
	//Get operator ID
	modifierID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, uda)
		return
	}
	uda.Modifier.ID = modifierID
	// Modify
	resStatus, _ = uda.Edit()
	ResponseWithMsg(c, resStatus, uda)
}

// Delete User-defined Archive handler
func DeleteUDAHandler(c *gin.Context) {
	uda := new(pg.UserDefinedArchive)
	err := c.ShouldBind(uda)
	if err != nil {
		zap.L().Error("DeleteUDAHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get operator ID
	modifierID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, uda)
		return
	}
	uda.Modifier.ID = modifierID
	// Delete
	resStatus, _ = uda.Delete()
	// Response
	ResponseWithMsg(c, resStatus, uda)
}

// Batch delete UDAs handler
func DeleteUDAsHandler(c *gin.Context) {
	udas := new([]pg.UserDefinedArchive)
	err := c.ShouldBind(udas)
	if err != nil {
		zap.L().Error("DeleteUDAsHandler invalid params:", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
	}
	// Get operator ID
	modifyUserId, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, udas)
		return
	}
	// Batch delete
	resStatus, _ = pg.DeleteUDAs(udas, modifyUserId)
	// Response
	ResponseWithMsg(c, resStatus, udas)
}

// Get User-defined Archive list under the UDC handler
func GetUDAListHandler(c *gin.Context) {
	udc := new(pg.UserDefineCategory)
	err := c.ShouldBind(udc)
	if err != nil {
		zap.L().Error("GetUDAsHandler invalid params:", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get list
	udas, restatus, _ := pg.GetUDAList(udc)
	// Response
	ResponseWithMsg(c, restatus, udas)
}

// Check if the UDA code exist handler
func CheckUDACodeExistHandler(c *gin.Context) {
	uda := new(pg.UserDefinedArchive)
	err := c.ShouldBind(uda)
	if err != nil {
		zap.L().Error("CheckUDACodeExistHandler invalid params:", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
	}

	// Check
	resStatus, _ := uda.CheckCodeExist()
	// Resopnse
	ResponseWithMsg(c, resStatus, uda)
}

// Get All UDA master data list handler
func GetUDAAllHandler(c *gin.Context) {
	udaAll, resStatus, _ := pg.GetUDAAll()
	ResponseWithMsg(c, resStatus, udaAll)
}

// Get latest UDA master data for front-end cache handler
func GetUDACacheHandler(c *gin.Context) {
	uddc := new(pg.UDACache)
	err := c.ShouldBind(uddc)
	if err != nil {
		zap.L().Error("GetUDACacheHandler invalid params:", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get front-end cache
	resStatus, _ := uddc.GetUDACache()
	// Response
	ResponseWithMsg(c, resStatus, uddc)
}
