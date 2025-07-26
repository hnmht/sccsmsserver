package route

import (
	"sccsmsserver/handlers"
	"sccsmsserver/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoute(g *gin.RouterGroup) {
	userGroup := g.Group("/user", middleware.JWTAuthMiddleware())
	{
		// //获取用户列表
		// userGroup.POST("/list", control.GetUsersHandler)
		// //获取权限菜单
		// userGroup.POST("/getmenu", control.GetMenuHandler)
		// //删除用户
		// userGroup.POST("/delete", control.DeleteUserHandler)
		// //批量删除用户
		// userGroup.POST("/deletemultiple", control.DeleteUsersHandler)
		// //检查用户代码是否存在
		// userGroup.POST("/validatUserCode", control.CheckUserCodeExistHandler)
		// //检查用户名称是否存在
		// userGroup.POST("/validatUserName", control.CheckUserNameExistHandler)
		// //增加用户
		// userGroup.POST("/add", control.AddUserHandler)
		// //编辑用户
		// userGroup.POST("/edit", control.EditUserHandler)
		// //更换头像
		// userGroup.POST("/changeAvatar", control.ChangeUserAvatarHandler)
		//发送token获取用户信息
		userGroup.POST("/userInfo", handlers.UserInfoHandler)
		// //通过用户中心修改用户信息
		// userGroup.POST("/modifyprofile", control.ModifyProfileHandler)
	}
}
