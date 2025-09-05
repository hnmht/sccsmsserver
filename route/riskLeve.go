package route

import (
	"sccsmsserver/handlers"
	"sccsmsserver/middleware"

	"github.com/gin-gonic/gin"
)

func RLRoute(g *gin.RouterGroup) {
	RLGroup := g.Group("/rl", middleware.CheckClientTypeMiddleware(), middleware.JWTAuthMiddleware())
	{
		// Add Risk Level
		RLGroup.POST("/add", handlers.AddRLHandler)
		// Get Risk Level List
		RLGroup.POST("/list", handlers.GetRLListHandler)
		// Modify Risk Level
		RLGroup.POST("/edit", handlers.EditRLHandler)
		// Delete Risk Level
		RLGroup.POST("/delete", handlers.DeleteRLHandler)
		// Batch datele Risk Level
		RLGroup.POST("/deleterls", handlers.DeleteRLsHandler)
		// Check if the Risk Level name exists
		RLGroup.POST("/checkname", handlers.CheckRLNameExistHandler)
		// Get Risk Level front-end cache
		RLGroup.POST("/cache", handlers.GetRLsCacheHandler)
	}
}
