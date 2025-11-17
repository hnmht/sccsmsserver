package route

import (
	"sccsmsserver/handlers"
	"sccsmsserver/middleware"

	"github.com/gin-gonic/gin"
)

func LandPageRoute(g *gin.RouterGroup) {
	landGroup := g.Group("/land")
	{
		// Get Landing Page Info
		landGroup.POST("/get", handlers.GetLandingPageInfoHandler)
		// Modify Landing Page Info
		landGroup.POST("/edit", middleware.CheckClientTypeMiddleware(), middleware.JWTAuthMiddleware(), handlers.ModifyLandingPageInfoHandler)
	}
}
