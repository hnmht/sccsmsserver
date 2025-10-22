package route

import (
	"sccsmsserver/handlers"
	"sccsmsserver/middleware"

	"github.com/gin-gonic/gin"
)

func DashboardRoute(g *gin.RouterGroup) {
	DashboardGroup := g.Group("/da", middleware.CheckClientTypeMiddleware(), middleware.JWTAuthMiddleware())
	{
		// Get Dashboard data
		DashboardGroup.POST("/data", handlers.GetDashboardDataHandler)
		// Get Risk Trends data
		DashboardGroup.POST("/risktrend", handlers.GetRiskTrendDataHandler)
	}
}
