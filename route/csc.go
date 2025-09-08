package route

import (
	"sccsmsserver/handlers"
	"sccsmsserver/middleware"

	"github.com/gin-gonic/gin"
)

func CSCRoute(g *gin.RouterGroup) {
	CSCGroup := g.Group("/csc", middleware.CheckClientTypeMiddleware(), middleware.JWTAuthMiddleware())
	{
		// Get CSC list
		CSCGroup.POST("/list", handlers.GetCSCListHandler)
		// Get Simple CSC list
		CSCGroup.POST("/simplist", handlers.GetSimpCSCListHandler)
		// Get front-end cache
		CSCGroup.POST("/simpcache", handlers.GetSimpCSCCacheHandler)
		// Check if the csc name exists
		CSCGroup.POST("/checkname", handlers.CheckCSCNameExistHandler)
		// Add CSC
		CSCGroup.POST("/add", handlers.AddCSCHandler)
		// Edit CSC
		CSCGroup.POST("/edit", handlers.EditCSCHandler)
		// Delete CSC
		CSCGroup.POST("/del", handlers.DeleteCSCHandler)
		// Batch delete CSC
		CSCGroup.POST("/dels", handlers.DeleteCSCsHandler)
	}
}
