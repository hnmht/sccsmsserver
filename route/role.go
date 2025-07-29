package route

import (
	"sccsmsserver/handlers"
	"sccsmsserver/middleware"

	"github.com/gin-gonic/gin"
)

func RoleRoute(g *gin.RouterGroup) {
	roleGroup := g.Group("/role")
	{
		//获取角色列表
		roleGroup.POST("/list", middleware.JWTAuthMiddleware(), handlers.GetRolesHandler)
		/* //查询角色名称是否存在
		roleGroup.POST("/validatename", middleware.JWTAuthMiddleware(), handlers.CheckRoleNameExistHandler)
		//增加角色
		roleGroup.POST("/add", middleware.JWTAuthMiddleware(), handlers.AddRoleHandler)
		//删除角色
		roleGroup.POST("/delete", middleware.JWTAuthMiddleware(), handlers.DeleteRoleHandler)
		//批量删除角色
		roleGroup.POST("/deleteroles", middleware.JWTAuthMiddleware(), handlers.DeleteRolesHandler)
		//更新角色
		roleGroup.POST("/edit", middleware.JWTAuthMiddleware(), handlers.EditRoleHandler)
		//查询角色权限菜单
		roleGroup.POST("/getmenus", middleware.JWTAuthMiddleware(), handlers.GetRoleMenusHandler)
		//更新用户角色权限
		roleGroup.POST("/updaterolemenus", middleware.JWTAuthMiddleware(), handlers.UpdateRoleMenusHandler) */
	}
}
