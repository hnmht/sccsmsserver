package route

import (
	"sccsmsserver/handlers"
	"sccsmsserver/middleware"

	"github.com/gin-gonic/gin"
)

func IRFRoute(g *gin.RouterGroup) {
	IRFGroup := g.Group("/irf", middleware.CheckClientTypeMiddleware(), middleware.JWTAuthMiddleware())
	{
		// Add Issue Resolution Form
		IRFGroup.POST("/add", handlers.AddIRFHandler)
		// Modify Issue Resolution Form
		IRFGroup.POST("/edit", handlers.EditIRFHandler)
		// Delete Issue Resolution Form
		IRFGroup.POST("/del", handlers.DeleteIRFhandler)
		// Confirm Issue Resolution Form
		IRFGroup.POST("/confirm", handlers.ConfirmIRFhandler)
		// UnConfirm Issue Resolution Form
		IRFGroup.POST("/unconfirm", handlers.UnConfirmIRFhandler)
		// Get Issue Resolution Form List
		IRFGroup.POST("/list", handlers.GetIRFListHandler)
	}
}
