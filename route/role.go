package route

import (
	"sccsmsserver/handlers"
	"sccsmsserver/middleware"

	"github.com/gin-gonic/gin"
)

func RoleRoute(g *gin.RouterGroup) {
	roleGroup := g.Group("/role", middleware.CheckClientTypeMiddleware(), middleware.JWTAuthMiddleware())
	{
		// Get role list
		roleGroup.POST("/list", handlers.GetRolesHandler)
		// Check the Role name exists
		roleGroup.POST("/checkname", handlers.CheckRoleNameExistHandler)
		// Add role
		roleGroup.POST("/add", handlers.AddRoleHandler)
		// Edit role
		roleGroup.POST("/edit", handlers.EditRoleHandler)
		// Delete Role
		roleGroup.POST("/del", handlers.DeleteRoleHandler)
		// Batch delete roles
		roleGroup.POST("/dels", handlers.DeleteRolesHandler)
		// Get role permissions
		roleGroup.POST("/getmenu", handlers.GetRoleMenusHandler)
		// Modify role permissions
		roleGroup.POST("/updaterolemenus", handlers.UpdateRoleMenusHandler)
	}
}
