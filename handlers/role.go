package handlers

import (
	"sccsmsserver/db/pg"
	"sccsmsserver/i18n"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Get role list handler
func GetRolesHandler(c *gin.Context) {
	// Get data
	roles, resStatus, _ := pg.GetRoles()
	ResponseWithMsg(c, resStatus, roles)
}

// Check the Role name exists handler
func CheckRoleNameExistHandler(c *gin.Context) {
	r := new(pg.Role)
	err := c.ShouldBind(r)
	if err != nil {
		zap.L().Error("CheckRoleNameExistHandler with invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	resStatus, _ := r.CheckNameExist()
	// Response
	ResponseWithMsg(c, resStatus, r)
}

// Add Role Handler
func AddRoleHandler(c *gin.Context) {
	r := new(pg.Role)
	err := c.ShouldBind(r)
	if err != nil {
		zap.L().Error("AddRoleHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}

	// Get Current user ID
	creatorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, r)
		return
	}
	r.Creator.ID = creatorID
	// Add the role
	resStatus, _ = r.Add()
	// Response
	ResponseWithMsg(c, resStatus, r)
}

// Delete Role Handler
func DeleteRoleHandler(c *gin.Context) {
	r := new(pg.Role)
	err := c.ShouldBind(r)
	if err != nil {
		zap.L().Error("DeleteRoleHandler invalid params", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Current User ID
	modifierID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, r)
		return
	}
	r.Modifier.ID = modifierID
	resStatus, _ = r.Delete()

	// Response
	ResponseWithMsg(c, resStatus, r)
}

// Edit Role Handler
func EditRoleHandler(c *gin.Context) {
	r := new(pg.Role)
	err := c.ShouldBind(r)
	if err != nil {
		zap.L().Error("EditRoleHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get current user ID
	modifierId, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, r)
		return
	}
	r.Modifier.ID = modifierId
	// Edit role
	statusCode, _ := r.Edit()
	// Response
	ResponseWithMsg(c, statusCode, r)
}

// Batch delete roles Handler
func DeleteRolesHandler(c *gin.Context) {
	r := new([]pg.Role)
	err := c.ShouldBind(r)
	if err != nil {
		zap.L().Error("DeleteRolesHandler invalid params", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operator ID
	operatorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, r)
		return
	}
	resStatus, _ = pg.DeleteRoles(r, operatorID)
	// Response
	ResponseWithMsg(c, resStatus, r)
}

// Get a list of role permission Handler
func GetRoleMenusHandler(c *gin.Context) {
	r := new(pg.Role)
	err := c.ShouldBind(r)
	if err != nil {
		zap.L().Error("AddRoleHandler invalid param:", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get menu list
	menus, resStatus, _ := r.GetRoleMenus()
	// Response
	ResponseWithMsg(c, resStatus, menus)
}

// Modify role permission Handler
func UpdateRoleMenusHandler(c *gin.Context) {
	r := new(pg.ParamsRoleMenu)
	err := c.ShouldBind(r)
	if err != nil {
		zap.L().Error("AddRoleHandler invalid param:", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operator ID
	operatorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, r)
		return
	}
	r.Role.Modifier.ID = operatorID
	// Update
	resStatus, _ = r.RoleMenuUpdate()
	// Response
	ResponseWithMsg(c, resStatus, r)
}
