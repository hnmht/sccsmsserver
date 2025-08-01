package route

import (
	"sccsmsserver/handlers"
	"sccsmsserver/middleware"

	"github.com/gin-gonic/gin"
)

func PersonRoute(g *gin.RouterGroup) {
	personGroup := g.Group("/person")
	{
		// Get Person Master Data list
		personGroup.POST("/list", middleware.JWTAuthMiddleware(), handlers.GetPersonsHandler)
		// Get Latest Person Master data for front-end caching
		personGroup.POST("/cache", middleware.JWTAuthMiddleware(), handlers.GetPersonsCacheHandler)
	}
}
