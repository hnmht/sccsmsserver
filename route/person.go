package route

import (
	"sccsmsserver/handlers"
	"sccsmsserver/middleware"

	"github.com/gin-gonic/gin"
)

func PersonRoute(g *gin.RouterGroup) {
	personGroup := g.Group("/person", middleware.CheckClientTypeMiddleware(), middleware.JWTAuthMiddleware())
	{
		// Get Person Master Data list
		personGroup.POST("/list", handlers.GetPersonsHandler)
		// Get Latest Person Master data for front-end caching
		personGroup.POST("/cache", handlers.GetPersonsCacheHandler)
	}
}
