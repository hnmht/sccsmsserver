package route

import (
	"sccsmsserver/handlers"
	"sccsmsserver/middleware"

	"github.com/gin-gonic/gin"
)

// Document Route Group
func DocRoute(g *gin.RouterGroup) {
	DocGroup := g.Group("/doc", middleware.CheckClientTypeMiddleware(), middleware.JWTAuthMiddleware())
	{
		// Add Document
		DocGroup.POST("/add", handlers.AddDocumentHandler)
		// Modify Document
		DocGroup.POST("/edit", handlers.EditDocumentHandler)
		// Get Document list pagination
		DocGroup.POST("/list", handlers.GetDocumentPagingListHanlder)
		// Delete Document
		DocGroup.POST("/del", handlers.DeleteDocumentHandler)
		// Batch Delte Document
		DocGroup.POST("/dels", handlers.DeleteDocumentsHandler)
		// Get Document Report
		DocGroup.POST("/rep", handlers.GetDocumentReportHandler)
	}
}
