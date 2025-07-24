package route

import (
	"sccsmsserver/handlers"

	"github.com/gin-gonic/gin"
)

func PubRoute(g *gin.RouterGroup) {
	pubGroup := g.Group("/pub")
	{
		//
		pubGroup.POST("/sysinfo", handlers.PubSystemInformationHandler)

	}
}
