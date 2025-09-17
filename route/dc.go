package route

import (
	"sccsmsserver/handlers"
	"sccsmsserver/middleware"

	"github.com/gin-gonic/gin"
)

func DCRoute(g *gin.RouterGroup) {
	DCGroup := g.Group("/dc", middleware.JWTAuthMiddleware())
	{
		// Get Document Categories List
		DCGroup.POST("/list", handlers.GetDCListHandler)
		// Get Simple Document Categories List
		DCGroup.POST("/simplist", handlers.GetSimpDCListHandler)
		// Get Simple Document Categories front-end Cache
		DCGroup.POST("/cache", handlers.GetSimpDCCacheHandler)
		// Check Document Category Name Exist
		DCGroup.POST("/checkname", handlers.CheckDCNameExistHandler)
		// Add Document Category
		DCGroup.POST("/add", handlers.AddDCHandler)
		// Edit Document Category
		DCGroup.POST("/edit", handlers.EditDCHandler)
		// Delete Document Category
		DCGroup.POST("/del", handlers.DeleteDCHandler)
		// Delete Multiple Document Categories
		DCGroup.POST("/dels", handlers.DeleteDCsHandler)
	}
}
