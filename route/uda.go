package route

import (
	"sccsmsserver/handlers"
	"sccsmsserver/middleware"

	"github.com/gin-gonic/gin"
)

func UDARoute(g *gin.RouterGroup) {
	UDAGroup := g.Group("/uda", middleware.JWTAuthMiddleware())
	{
		// Add UDA
		UDAGroup.POST("/add", handlers.AddUDAHandler)
		// Edit UDA
		UDAGroup.POST("/edit", handlers.EditUDAHandler)
		// Delete UDA
		UDAGroup.POST("/delete", handlers.DeleteUDAHandler)
		// Batch delete UDA
		UDAGroup.POST("/deleteudas", handlers.DeleteUDAsHandler)
		// Get UDA list under the UDC
		UDAGroup.POST("/list", handlers.GetUDAListHandler)
		// Get all UDA list
		UDAGroup.POST("/all", handlers.GetUDAAllHandler)
		// Check if the UDA code exist
		UDAGroup.POST("/checkcode", handlers.CheckUDACodeExistHandler)
		// Get front-end cache
		UDAGroup.POST("/cache", handlers.GetUDACacheHandler)
	}
}
