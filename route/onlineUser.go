package route

import (
	"sccsmsserver/handlers"
	"sccsmsserver/middleware"

	"github.com/gin-gonic/gin"
)

func OuRoute(g *gin.RouterGroup) {
	ouGroup := g.Group("/ou", middleware.CheckClientTypeMiddleware(), middleware.JWTAuthMiddleware())
	{
		// Get online user list
		ouGroup.POST("/list", handlers.GetOnlineUserHandler)
		// Remove online user
		ouGroup.POST("/remove", handlers.RemoveOnlineUserHandler)
	}
}
