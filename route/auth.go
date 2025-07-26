package route

import (
	"sccsmsserver/handlers"
	"sccsmsserver/middleware"

	"github.com/gin-gonic/gin"
)

func AuthRoute(g *gin.RouterGroup) {
	authGroup := g.Group("/auth", middleware.CheckClientTypeMiddleware())
	{
		// User Login
		authGroup.POST("/login", handlers.LoginHandler)
		// // Change user password
		// authGroup.POST("/changepwd", middleware.JWTAuthMiddleware(), control.ChangeUserPasswordHandler)
		// //注销登录
		// authGroup.POST("/logout", middleware.JWTAuthMiddleware(), control.LogoutHandler)
		//Rsa public key
		authGroup.POST("/publickey", handlers.GetPublicKeyHandler)
	}
}
