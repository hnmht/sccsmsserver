package route

import (
	"sccsmsserver/handlers"
	"sccsmsserver/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoute(g *gin.RouterGroup) {
	userGroup := g.Group("/user", middleware.JWTAuthMiddleware())
	{
		//获取用户列表
		userGroup.POST("/list", handlers.GetUsersHandler)
		//获取权限菜单
		// userGroup.POST("/getmenu", handlers.GetMenuHandler)
		//删除用户
		userGroup.POST("/delete", handlers.DeleteUserHandler)
		//批量删除用户
		userGroup.POST("/deletemultiple", handlers.DeleteUsersHandler)
		//检查用户代码是否存在
		userGroup.POST("/validatecode", handlers.CheckUserCodeExistHandler)
		//检查用户名称是否存在
		userGroup.POST("/validatename", handlers.CheckUserNameExistHandler)
		//增加用户
		userGroup.POST("/add", handlers.AddUserHandler)
		//编辑用户
		userGroup.POST("/edit", handlers.EditUserHandler)
		//更换头像
		userGroup.POST("/changeAvatar", handlers.ChangeUserAvatarHandler)
		// 发送token获取用户信息
		userGroup.POST("/userInfo", handlers.UserInfoHandler)
		//通过用户中心修改用户信息
		userGroup.POST("/modifyprofile", handlers.ModifyProfileHandler)
	}
}
