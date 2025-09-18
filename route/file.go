package route

import (
	"sccsmsserver/handlers"
	"sccsmsserver/middleware"

	"github.com/gin-gonic/gin"
)

func FileRoute(g *gin.RouterGroup) {
	FileGroup := g.Group("/file", middleware.CheckClientTypeMiddleware(), middleware.JWTAuthMiddleware())
	{
		// Receive client files upload
		FileGroup.POST("/receive", handlers.RecieveFilesHandler)
		// Get file information by file ha
		FileGroup.POST("/getfilebyhash", handlers.GetFileInfoByHashHandler)
		// Get files information by file hash array
		FileGroup.POST("/getfilesbyhash", handlers.GetFilesByHashHandler)
	}
}
