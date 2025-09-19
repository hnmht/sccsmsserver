package route

import (
	"sccsmsserver/handlers"
	"sccsmsserver/middleware"

	"github.com/gin-gonic/gin"
)

func TCRoute(g *gin.RouterGroup) {
	TCGroup := g.Group("/tc", middleware.CheckClientTypeMiddleware(), middleware.JWTAuthMiddleware())
	{
		// Get Training Course List
		TCGroup.POST("/list", handlers.GetTCListHandler)
		// Get Training Course Frontend Cache
		TCGroup.POST("/cache", handlers.GetTCCacheHandler)
		// Check Training Course Name Exist
		TCGroup.POST("/checkname", handlers.CheckTCNameExistHandler)
		// Add
		TCGroup.POST("/add", handlers.AddTCHandler)
		// Edit
		TCGroup.POST("/edit", handlers.EditTCHandler)
		// Delete
		TCGroup.POST("/del", handlers.DeleteTCHandler)
		// Batch Delete
		TCGroup.POST("/dels", handlers.DeleteTCsHandler)
	}
}
