package route

import (
	"sccsmsserver/handlers"
	"sccsmsserver/middleware"

	"github.com/gin-gonic/gin"
)

func PositionRoute(g *gin.RouterGroup) {
	PositionGroup := g.Group("/position", middleware.CheckClientTypeMiddleware(), middleware.JWTAuthMiddleware())
	{
		// Add position
		PositionGroup.POST("/add", handlers.AddPositionHandler)
		// Get position list
		PositionGroup.POST("/list", handlers.GetPositionListHandler)
		// Edit position
		PositionGroup.POST("/edit", handlers.EditPositionHandler)
		// Delete position
		PositionGroup.POST("/delete", handlers.DeletePositionHandler)
		// Batch delete positions
		PositionGroup.POST("/batchdelete", handlers.DeletePositionsHandler)
		// Check name exists
		PositionGroup.POST("/checkname", handlers.CheckPositionNameExistHandler)
		// Get position master data front-end cache
		PositionGroup.POST("/cache", handlers.GetPositionCacheHandler)
	}
}
