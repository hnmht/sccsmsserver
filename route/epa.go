package route

import (
	"sccsmsserver/handlers"
	"sccsmsserver/middleware"

	"github.com/gin-gonic/gin"
)

func EPARoute(g *gin.RouterGroup) {
	EPAGroup := g.Group("/epa", middleware.JWTAuthMiddleware())
	{
		// Get EP list
		EPAGroup.POST("/list", handlers.GetEPListHandler)
		// Get front-end cache
		EPAGroup.POST("/cache", handlers.GetEPCacheHandler)
		// Add EP
		EPAGroup.POST("/add", handlers.AddEPHandler)
		// Modify EP
		EPAGroup.POST("/edit", handlers.EditEPHandler)
		// Delete EP
		EPAGroup.POST("/delete", handlers.DeleteEPHandler)
		// Batche Delete EP
		EPAGroup.POST("/deletes", handlers.DeleteEPsHandler)
		// Check if the EP's code exists
		EPAGroup.POST("/checkcode", handlers.CheckEPCodeExistHandler)
	}
}
