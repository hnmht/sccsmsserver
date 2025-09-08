package route

import (
	"sccsmsserver/handlers"
	"sccsmsserver/middleware"

	"github.com/gin-gonic/gin"
)

func CSARoute(g *gin.RouterGroup) {
	CSAGroup := g.Group("/csa", middleware.CheckClientTypeMiddleware(), middleware.JWTAuthMiddleware())
	{
		// Add CSA
		CSAGroup.POST("/add", handlers.AddCSHandler)
		// Modify CSA
		CSAGroup.POST("/edit", handlers.EditCSHandler)
		// Check if the CSA code exists
		CSAGroup.POST("/checkcode", handlers.CheckCSCodeExistHandler)
		// Delete CSA
		CSAGroup.POST("/del", handlers.DeleteCSHandler)
		// Batch delete CSA
		CSAGroup.POST("/dels", handlers.DeleteCSsHandler)
		// Get CSA list
		CSAGroup.POST("/list", handlers.GetCSsHandler)
		// Get CSA front-end cache
		CSAGroup.POST("/cache", handlers.GetCSCacheHandler)
	}
}
