package route

import (
	"sccsmsserver/handlers"

	"github.com/gin-gonic/gin"
)

func AuthRoute(g *gin.RouterGroup) {
	authGroup := g.Group("/auth")
	{
		/* //登录
		authGroup.POST("/login", handlers.LoginHandler)
		//更改用户密码
		authGroup.POST("/changepwd", middleware.JWTAuthMiddleware(), control.ChangeUserPasswordHandler)
		//注销登录
		authGroup.POST("/logout", middleware.JWTAuthMiddleware(), control.LogoutHandler) */
		//Rsa public key
		authGroup.POST("/publickey", handlers.GetPublicKeyHandler)
	}
}
