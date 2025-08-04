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
		ResponseWithMsg(c, i18n.CodeInvalidParm, nil)
		return
	}
	resStatus, _ := r.CheckNameExist()
	//返回
	ResponseWithMsg(c, resStatus, r)
}

// Add Role Handler
func AddRoleHandler(c *gin.Context) {
	r := new(pg.Role)
	err := c.ShouldBind(r)
	if err != nil {
		zap.L().Error("AddRoleHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, nil)
		return
	}

	// Get Current user ID
	creatorID, resStatus := GetCurrentUser(c)
	if resStatus != i18n.StatusOK {
		zap.L().Error("AddRoleHandler getCurrentUser failed", zap.Error(err))
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
		ResponseWithMsg(c, i18n.CodeInvalidParm, nil)
		return
	}
	// Get Current User ID
	modifierID, resStatus := GetCurrentUser(c)
	if resStatus != i18n.StatusOK {
		zap.L().Error("DeleteRoleHandler getCurrentUser failed", zap.Error(err))
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
		ResponseWithMsg(c, i18n.CodeInvalidParm, nil)
		return
	}
	// Get current user ID
	modifierId, resStatus := GetCurrentUser(c)
	if resStatus != i18n.StatusOK {
		zap.L().Error("EditRoleHandler getCurrentUser failed", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInternalError, r)
		return
	}
	r.Modifier.ID = modifierId

	// Edit role
	statusCode, _ := r.Edit()
	// Response
	ResponseWithMsg(c, statusCode, r)
}

/*//DeleteRolesHandler 批量删除角色
func DeleteRolesHandler(c *gin.Context) {
	r := new([]pg.Role)
	err := c.ShouldBind(r)
	if err != nil {
		zap.L().Error("DeleteRolesHandler invalid params", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, pub.CodeInvalidParm)
			return
		}
		ResponseErrorWithMsg(c, pub.CodeInvalidParm, removeTopStruct(errs.Translate(trans)))
		return
	}
	//获取操作用户id
	modifyUserId, err := GetCurrentUser(c)
	if err != nil {
		zap.L().Error("DeleteRoleHandler getCurrentUser failed", zap.Error(err))
		ResponseErrorWithMsg(c, pub.CodeInternalError, r)
		return
	}

	statusCode, _ := pg.DeleteRoles(r, modifyUserId)
	//返回
	ResponseSuccess(c, statusCode, r)
}

//GetRoleMenusHandler 获取角色权限列表
func GetRoleMenusHandler(c *gin.Context) {
	r := new(pg.Role)
	err := c.ShouldBind(r)
	if err != nil {
		zap.L().Error("AddRoleHandler invalid param", zap.Error(err))
		//判断err是不是validator.validationErrors类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, pub.CodeInvalidParm)
			return
		}
		ResponseErrorWithMsg(c, pub.CodeInternalError, removeTopStruct(errs.Translate(trans)))
		return
	}
	//从数据库提取角色权限列表
	menus, resStatus, _ := r.GetRoleMenus()

	ResponseSuccess(c, resStatus, menus)
}



//UpdateRoleMenus 更新角色权限
func UpdateRoleMenusHandler(c *gin.Context) {
	r := new(pg.ParamsRoleMenu)
	err := c.ShouldBind(r)
	if err != nil {
		zap.L().Error("AddRoleHandler invalid param", zap.Error(err))
		//判断err是不是validator.validationErrors类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, pub.CodeInvalidParm)
			return
		}
		ResponseErrorWithMsg(c, pub.CodeInvalidParm, removeTopStruct(errs.Translate(trans)))
		return
	}
	//获取操作用户id
	modifyUserId, err := GetCurrentUser(c)
	if err != nil {
		zap.L().Error("UpdateRoleMenus getCurrentUser failed", zap.Error(err))
		ResponseErrorWithMsg(c, pub.CodeInternalError, r)
		return
	}
	r.Role.ModifyUser.UserID = modifyUserId

	statusCode, _ := r.RoleMenuUpdate()

	ResponseSuccess(c, statusCode, r)
}
*/
