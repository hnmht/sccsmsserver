package handlers

import (
	"sccsmsserver/db/pg"

	"github.com/gin-gonic/gin"
)

// Get role list handler
func GetRolesHandler(c *gin.Context) {
	// Get data
	roles, resStatus, _ := pg.GetRoles()

	ResponseWithMsg(c, resStatus, roles)
}

/* //CheckRoleNameExistHandler 检查角色名称是否存在
func CheckRoleNameExistHandler(c *gin.Context) {
	r := new(pg.Role)
	if err := c.ShouldBind(r); err != nil {
		//请求参数有误，记录日志并返回响应
		zap.L().Error("CheckRoleNameExistHandler with invalid param", zap.Error(err))
		//判断err是不是validator.validationErrors类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, pub.CodeInvalidParm)
			return
		}
		ResponseErrorWithMsg(c, pub.CodeInvalidParm, removeTopStruct(errs.Translate(trans)))
		return
	}
	//检查
	resStatus, _ := r.CheckNameExist()
	//返回
	ResponseSuccess(c, resStatus, r)
}

//AddRoleHandler 增加角色
func AddRoleHandler(c *gin.Context) {
	r := new(pg.Role)
	err := c.ShouldBind(r)
	if err != nil {
		//请求参数有误，记录日志并返回响应
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
	createUserId, err := GetCurrentUser(c)
	if err != nil {
		zap.L().Error("AddRoleHandler getCurrentUser failed", zap.Error(err))
		ResponseErrorWithMsg(c, pub.CodeInternalError, r)
		return
	}
	r.CreateUser.UserID = createUserId
	//向数据库添加角色
	resStatus, _ := r.Add()
	//返回
	ResponseSuccess(c, resStatus, r)
}

//DeleteRoleHandler 删除角色
func DeleteRoleHandler(c *gin.Context) {
	r := new(pg.Role)
	err := c.ShouldBind(r)
	if err != nil {
		zap.L().Error("DeleteRoleHandler invalid params", zap.Error(err))
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
		zap.L().Error("DeleteRoleHandler getCurrentUser failed", zap.Error(err))
		ResponseErrorWithMsg(c, pub.CodeInternalError, r)
		return
	}
	r.ModifyUser.UserID = modifyUserId

	statusCode, _ := r.Delete()

	//返回
	ResponseSuccess(c, statusCode, r)
}

//DeleteRolesHandler 批量删除角色
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

//EditRoleHandler 修改角色
func EditRoleHandler(c *gin.Context) {
	r := new(pg.Role)
	err := c.ShouldBind(r)
	if err != nil {
		zap.L().Error("EditRoleHandler invalid param", zap.Error(err))
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
		zap.L().Error("EditRoleHandler getCurrentUser failed", zap.Error(err))
		ResponseErrorWithMsg(c, pub.CodeInternalError, r)
		return
	}
	r.ModifyUser.UserID = modifyUserId

	//修改角色
	statusCode, _ := r.Edit()
	//返回响应
	ResponseSuccess(c, statusCode, r)
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
