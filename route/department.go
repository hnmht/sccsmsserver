package route

import (
	"sccsmsserver/handlers"
	"sccsmsserver/middleware"

	"github.com/gin-gonic/gin"
)

func DeptRoute(g *gin.RouterGroup) {
	deptGroup := g.Group("/dept", middleware.CheckClientTypeMiddleware(), middleware.JWTAuthMiddleware())
	{
		// Get department list
		deptGroup.POST("/list", handlers.GetDeptsHandler)
		// Get simplify department list
		deptGroup.POST("/simplist", handlers.GetSimpDeptsHandler)
		// Check if the department code eists
		deptGroup.POST("/checkcode", handlers.CheckDeptCodeExistHandler)
		// Add department
		deptGroup.POST("/add", handlers.AddDeptHandler)
		//  Get Simple Department latest front-end cache
		deptGroup.POST("/simpcache", handlers.GetSimpDeptsCacheHandler)
		// Modify department
		deptGroup.POST("/edit", handlers.EditDeptHandler)
		// Delete department
		deptGroup.POST("/del", handlers.DelDeptHandler)
		// Batch department
		deptGroup.POST("/dels", handlers.DelDeptsHandler)
	}
}
