package route

import (
	"sccsmsserver/handlers"
	"sccsmsserver/middleware"

	"github.com/gin-gonic/gin"
)

func UDCRoute(g *gin.RouterGroup) {
	UDCGroup := g.Group("/udc", middleware.CheckClientTypeMiddleware(), middleware.JWTAuthMiddleware())
	{
		// Add UDC
		UDCGroup.POST("/add", handlers.AddUDCHandler)
		// Get UDC list
		UDCGroup.POST("/list", handlers.GetUDCListHandler)
		// Edit UDC
		UDCGroup.POST("/edit", handlers.EditUDCHandler)
		// Delete UDC
		UDCGroup.POST("/delete", handlers.DeleteUDCHandler)
		// Batch delete UDC
		UDCGroup.POST("/deleteudcs", handlers.DeleteUDCsHandler)
		// Check if the UDC name exists
		UDCGroup.POST("/checkname", handlers.CheckUDCNameExistHandler)
		// Get latest UDC front-end cache
		UDCGroup.POST("/cache", handlers.GetUDCsCacheHandler)
	}
}
