package route

import (
	"sccsmsserver/handlers"
	"sccsmsserver/middleware"

	"github.com/gin-gonic/gin"
)

func PubRoute(g *gin.RouterGroup) {
	pubGroup := g.Group("/pub", middleware.CheckClientTypeMiddleware(), middleware.JWTAuthMiddleware())
	{
		// Public System Information
		pubGroup.POST("/sysinfo", handlers.PubSystemInformationHandler)
		// Generate frontend DBID
		pubGroup.POST("/addfrontdbid", handlers.GenerateFrontendDBID)
		// Get frontend dbid
		pubGroup.POST("/getfrontdbid", handlers.GetFrontendDBInfo)
	}
}
