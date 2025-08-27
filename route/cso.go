package route

import (
	"sccsmsserver/handlers"
	"sccsmsserver/middleware"

	"github.com/gin-gonic/gin"
)

func CSORoute(g *gin.RouterGroup) {
	CSOGroup := g.Group("/cso", middleware.CheckClientTypeMiddleware(), middleware.JWTAuthMiddleware())
	{
		// Get CSO list
		CSOGroup.POST("/options", handlers.GetCSOsHandler)
		// Modify CSO
		CSOGroup.POST("/editoption", handlers.EditCSOHandler)
		// Get CSO front-end cache
		CSOGroup.POST("/optioncache", handlers.GetCSOCacheHandler)
	}
}
