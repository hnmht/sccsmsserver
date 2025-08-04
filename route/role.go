package route

import (
	"sccsmsserver/handlers"
	"sccsmsserver/middleware"

	"github.com/gin-gonic/gin"
)

func RoleRoute(g *gin.RouterGroup) {
	roleGroup := g.Group("/role")
	{
		// Get Role list
		roleGroup.POST("/list", middleware.JWTAuthMiddleware(), handlers.GetRolesHandler)
		// Check the Role name exists
		roleGroup.POST("/validatename", middleware.JWTAuthMiddleware(), handlers.CheckRoleNameExistHandler)
		// Add Role
		roleGroup.POST("/add", middleware.JWTAuthMiddleware(), handlers.AddRoleHandler)
		// Edit Role
		roleGroup.POST("/edit", middleware.JWTAuthMiddleware(), handlers.EditRoleHandler)
		// Delete Role
		roleGroup.POST("/delete", middleware.JWTAuthMiddleware(), handlers.DeleteRoleHandler)
		/*//批量删除角色
		roleGroup.POST("/deleteroles", middleware.JWTAuthMiddleware(), handlers.DeleteRolesHandler)
		//查询角色权限菜单
		roleGroup.POST("/getmenus", middleware.JWTAuthMiddleware(), handlers.GetRoleMenusHandler)
		//更新用户角色权限
		roleGroup.POST("/updaterolemenus", middleware.JWTAuthMiddleware(), handlers.UpdateRoleMenusHandler) */
	}
}
