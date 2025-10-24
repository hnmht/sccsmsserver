package route

import (
	"sccsmsserver/handlers"
	"sccsmsserver/middleware"

	"github.com/gin-gonic/gin"
)

func EventRoute(g *gin.RouterGroup) {
	EventGroup := g.Group("/event", middleware.CheckClientTypeMiddleware(), middleware.JWTAuthMiddleware())
	{
		// Get User Events
		EventGroup.POST("/list", handlers.GetEventsHandler)
	}
}
