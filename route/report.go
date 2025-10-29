package route

import (
	"sccsmsserver/handlers"
	"sccsmsserver/middleware"

	"github.com/gin-gonic/gin"
)

func RepRoute(g *gin.RouterGroup) {
	REPGroup := g.Group("/rep", middleware.CheckClientTypeMiddleware(), middleware.JWTAuthMiddleware())
	{
		// Work Order Status Report
		REPGroup.POST("/wor", handlers.GetWoReportHandler)
		// Execution Order status Report
		REPGroup.POST("/edr", handlers.GetEoReportHandler)
		// Issue Resolution Form Report
		REPGroup.POST("/ddr", handlers.GetIRFReportHandler)
	}
}
