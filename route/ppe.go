package route

import (
	"sccsmsserver/handlers"
	"sccsmsserver/middleware"

	"github.com/gin-gonic/gin"
)

func PPERoute(g *gin.RouterGroup) {
	PPEGroup := g.Group("/ppe", middleware.CheckClientTypeMiddleware(), middleware.JWTAuthMiddleware())
	{
		// Add PPE
		PPEGroup.POST("/add", handlers.AddPPEHandler)
		// Get PPE list
		PPEGroup.POST("/list", handlers.GetPPEListHandler)
		// Modify PPE
		PPEGroup.POST("/edit", handlers.EditPPEHandler)
		// Delete PPE
		PPEGroup.POST("/del", handlers.DeletePPEHandler)
		// Batch Delete PPE
		PPEGroup.POST("/dels", handlers.DeletePPEsHandler)
		// Check if the PPE code exists
		PPEGroup.POST("/checkcode", handlers.CheckPPECodeExistHandler)
		// Get latest PPE front-end cache
		PPEGroup.POST("/cache", handlers.GetPPECacheHandler)
	}
}
