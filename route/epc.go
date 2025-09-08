package route

import (
	"sccsmsserver/handlers"
	"sccsmsserver/middleware"

	"github.com/gin-gonic/gin"
)

func EPCRoute(g *gin.RouterGroup) {
	EPCGroup := g.Group("/epc", middleware.CheckClientTypeMiddleware(), middleware.JWTAuthMiddleware())
	{
		// Get EPC list
		EPCGroup.POST("/list", handlers.GetEPCListHandler)
		// Get SimpEPC list
		EPCGroup.POST("/simplist", handlers.GetSimpEPCListHandler)
		// Get SimpEPC front-end cache
		EPCGroup.POST("/simpcache", handlers.GetSimpEPCCacheHandler)
		// Check if the EPC name exists
		EPCGroup.POST("/checkname", handlers.CheckEPCNameExistHandler)
		// Add EPC
		EPCGroup.POST("/add", handlers.AddEPCHandler)
		// Modify EPC
		EPCGroup.POST("/edit", handlers.EditEPCHandler)
		// Delete EPC
		EPCGroup.POST("/del", handlers.DeleteEPCHandler)
		// Batch delete EPC
		EPCGroup.POST("/dels", handlers.DeleteEPCsHandler)
	}
}
