package route

import (
	"sccsmsserver/handlers"
	"sccsmsserver/middleware"

	"github.com/gin-gonic/gin"
)

func RoleRoute(g *gin.RouterGroup) {
	roleGroup := g.Group("/role")
	{
		// Get role list
		roleGroup.POST("/list", middleware.JWTAuthMiddleware(), handlers.GetRolesHandler)
		// Check the Role name exists
		roleGroup.POST("/checkname", middleware.JWTAuthMiddleware(), handlers.CheckRoleNameExistHandler)
		// Add role
		roleGroup.POST("/add", middleware.JWTAuthMiddleware(), handlers.AddRoleHandler)
		// Edit role
		roleGroup.POST("/edit", middleware.JWTAuthMiddleware(), handlers.EditRoleHandler)
		// Delete Role
		roleGroup.POST("/del", middleware.JWTAuthMiddleware(), handlers.DeleteRoleHandler)
		// Batch delete roles
		roleGroup.POST("/dels", middleware.JWTAuthMiddleware(), handlers.DeleteRolesHandler)
		// Get role permissions
		roleGroup.POST("/getmenus", middleware.JWTAuthMiddleware(), handlers.GetRoleMenusHandler)
		// Modify role permissions
		roleGroup.POST("/updaterolemenus", middleware.JWTAuthMiddleware(), handlers.UpdateRoleMenusHandler)
	}
}
