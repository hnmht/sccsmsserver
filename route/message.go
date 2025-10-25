package route

import (
	"sccsmsserver/handlers"
	"sccsmsserver/middleware"

	"github.com/gin-gonic/gin"
)

func MsgRoute(g *gin.RouterGroup) {
	MSGGroup := g.Group("/msg", middleware.CheckClientTypeMiddleware(), middleware.JWTAuthMiddleware())
	{
		// Get user's UnRead message
		MSGGroup.POST("/unread", handlers.GetUserUnReadCommentsHandler)
		// Get user's read message
		MSGGroup.POST("/read", handlers.GetUserReadCommentsHandler)
		// Get user Work Orders awaiting Execution
		MSGGroup.POST("/wos", handlers.GetUserWORefsHandler)
		// Get user Execution Order issues
		MSGGroup.POST("/eos", handlers.GetUserEORefsHandler)
		// Read message
		MSGGroup.POST("/toread", handlers.ReadCommentMessageHandler)
	}
}
