package route

import (
	"sccsmsserver/handlers"
	"sccsmsserver/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoute(g *gin.RouterGroup) {
	userGroup := g.Group("/user", middleware.JWTAuthMiddleware())
	{
		// Get User list
		userGroup.POST("/list", handlers.GetUsersHandler)
		// Get Menu List
		userGroup.POST("/getmenu", handlers.GetMenuHandler)
		// Delete User
		userGroup.POST("/del", handlers.DeleteUserHandler)
		// Batch Delete User
		userGroup.POST("/dels", handlers.DeleteUsersHandler)
		// Check if the user code exists
		userGroup.POST("/checkcode", handlers.CheckUserCodeExistHandler)
		// Check if the user name exists
		userGroup.POST("/checkname", handlers.CheckUserNameExistHandler)
		// Add User
		userGroup.POST("/add", handlers.AddUserHandler)
		// Edit User
		userGroup.POST("/edit", handlers.EditUserHandler)
		// Change user avatar
		userGroup.POST("/changeavatar", handlers.ChangeUserAvatarHandler)
		// Get user information based on token
		userGroup.POST("/info", handlers.UserInfoHandler)
		// User update VIA personal center
		userGroup.POST("/modifyprofile", handlers.ModifyProfileHandler)
	}
}
