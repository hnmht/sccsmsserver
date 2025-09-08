package route

import (
	"sccsmsserver/handlers"
	"sccsmsserver/middleware"

	"github.com/gin-gonic/gin"
)

func EPARoute(g *gin.RouterGroup) {
	EPAGroup := g.Group("/epa", middleware.CheckClientTypeMiddleware(), middleware.JWTAuthMiddleware())
	{
		// Get EP list
		EPAGroup.POST("/list", handlers.GetEPListHandler)
		// Get latest front-end cache
		EPAGroup.POST("/cache", handlers.GetEPCacheHandler)
		// Add EP
		EPAGroup.POST("/add", handlers.AddEPHandler)
		// Modify EP
		EPAGroup.POST("/edit", handlers.EditEPHandler)
		// Delete EP
		EPAGroup.POST("/del", handlers.DeleteEPHandler)
		// Batche Delete EP
		EPAGroup.POST("/dels", handlers.DeleteEPsHandler)
		// Check if the EP's code exists
		EPAGroup.POST("/checkcode", handlers.CheckEPCodeExistHandler)
	}
}
