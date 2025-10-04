package route

import (
	"sccsmsserver/handlers"
	"sccsmsserver/middleware"

	"github.com/gin-gonic/gin"
)

func WORoute(g *gin.RouterGroup) {
	WOGroup := g.Group("/wo", middleware.CheckClientTypeMiddleware(), middleware.JWTAuthMiddleware())
	{
		// Get Work Order list
		WOGroup.POST("/list", handlers.GetWOListHanlder)
		// Get Work Order details
		WOGroup.POST("/detail", handlers.GetWOInfoByIDHandler)
		// Add Work Order
		WOGroup.POST("/add", handlers.AddWOHandler)
		// Edit Work Order
		WOGroup.POST("/edit", handlers.EditWOHandler)
		// Delete Work Order
		WOGroup.POST("/del", handlers.DeleteWOHandler)
		// Batch delete Work Order
		WOGroup.POST("/dels", handlers.DeleteWOsHandler)
		// Confirm Work Order
		WOGroup.POST("/confirm", handlers.ConfirmWOHandler)
		// Unconfirm Work Order
		WOGroup.POST("/unconfirm", handlers.UnConfirmWOHandler)
		// Get the list of Work Order awaiting execution
		WOGroup.POST("/refer", handlers.GetWOReferHandler)
	}
}
