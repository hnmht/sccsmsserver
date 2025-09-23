package route

import (
	"sccsmsserver/handlers"
	"sccsmsserver/middleware"

	"github.com/gin-gonic/gin"
)

func EPTRoute(g *gin.RouterGroup) {
	EPTGroup := g.Group("/ept", middleware.CheckClientTypeMiddleware(), middleware.JWTAuthMiddleware())
	{
		// Get execution project template list
		EPTGroup.POST("/list", handlers.GetEPTListHandler)
		// Get latest Execution Project Template for fontend cache
		EPTGroup.POST("/cache", handlers.GetEPTCacheHandler)
		// Add Execution Project Template
		EPTGroup.POST("/add", handlers.AddEPTHandler)
		// Edit Execution Project Template
		EPTGroup.POST("/edit", handlers.EditEPTHandler)
		// Delete Execution Project Template
		EPTGroup.POST("/del", handlers.DeleteEPTHandler)
		// Batch delete Execution Project Template
		EPTGroup.POST("dels", handlers.DeleteEPTsHandler)
		// Check if the execution project template code exists
		EPTGroup.POST("/checkcode", handlers.CheckEPTCodeExistHandler)
	}
}
