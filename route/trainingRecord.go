package route

import (
	"sccsmsserver/handlers"
	"sccsmsserver/middleware"

	"github.com/gin-gonic/gin"
)

func TRRoute(g *gin.RouterGroup) {
	TRGroup := g.Group("/tr", middleware.CheckClientTypeMiddleware(), middleware.JWTAuthMiddleware())
	{
		// Add Training Record
		TRGroup.POST("/add", handlers.AddTRHandler)
		// Get Training Record list
		TRGroup.POST("/list", handlers.GetTRListHandler)
		// Get Training Record details
		TRGroup.POST("/detail", handlers.GetTRInfoByHIDHandler)
		// Modify Training Record
		TRGroup.POST("/edit", handlers.EditTRHandler)
		// Delete Training Record
		TRGroup.POST("/del", handlers.DeleteTRHandler)
		// Confirm Training Record
		TRGroup.POST("/confirm", handlers.ConfirmTRHandler)
		// UnConfirm Training Record
		TRGroup.POST("/unconfirm", handlers.UnConfirmTRHandler)
		// Get Taught Lessons Report
		TRGroup.POST("/tlrep", handlers.GetTaughtLessonsReportHandler)
		// Get Recieved Training Report
		TRGroup.POST("/rtrep", handlers.GetRecivedTrainingReportHandler)
	}
}
