package route

import (
	"sccsmsserver/handlers"
	"sccsmsserver/middleware"

	"github.com/gin-gonic/gin"
)

func EORoute(g *gin.RouterGroup) {
	EOGroup := g.Group("/eo", middleware.CheckClientTypeMiddleware(), middleware.JWTAuthMiddleware())
	{
		// Get Execution Order List
		EOGroup.POST("/list", handlers.GetEOListHandler)
		// Get the list of Execution Orders to be reviewed by pagination
		EOGroup.POST("/listpage", handlers.GetEOReviewListPaginationHandler)
		// Add Execution Order
		EOGroup.POST("/add", handlers.AddEOHandler)
		// Edit Execution Order
		EOGroup.POST("/edit", handlers.EditEOHandler)
		// Delete Execution Order
		EOGroup.POST("/del", handlers.DeleteEOHandler)
		// Confirm Execution Order
		EOGroup.POST("/confirm", handlers.ConfirmEOHandler)
		// Un-Confirm Execution Order
		EOGroup.POST("/unconfirm", handlers.CancelConfirmEOHandler)
		// Get the Execution Order details
		EOGroup.POST("/detail", handlers.GetEOInfoByHIDHandler)
		// Get the list of Execution Orders to be referenced
		EOGroup.POST("/refer", handlers.GetReferEOHandler)
		// Add Execution Order comment
		EOGroup.POST("/addcomment", handlers.AddCommentHandler)
		// Add Execution Order Review Record
		EOGroup.POST("/addreview", handlers.AddReviewHandler)
		// Get Execution Order Review Records list
		EOGroup.POST("/reviews", handlers.GetEOReviewsHandler)
		// Get Execution Order Comments list
		EOGroup.POST("/comments", handlers.GetEOCommentsHandler)
	}
}
