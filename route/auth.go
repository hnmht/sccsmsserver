package route

import (
	"sccsmsserver/handlers"
	"sccsmsserver/middleware"

	"github.com/gin-gonic/gin"
)

func AuthRoute(g *gin.RouterGroup) {
	authGroup := g.Group("/auth", middleware.CheckClientTypeMiddleware())
	{
		//Rsa public key
		authGroup.POST("/publickey", handlers.GetPublicKeyHandler)
		// Validate token
		authGroup.POST("/validatetoken", middleware.JWTAuthMiddleware(), handlers.ValidateToken)
		// User Login
		authGroup.POST("/login", handlers.LoginHandler)
		// Change user password
		authGroup.POST("/changepwd", middleware.JWTAuthMiddleware(), handlers.ChangeUserPasswordHandler)
		// Logout
		authGroup.POST("/logout", middleware.JWTAuthMiddleware(), handlers.LogoutHandler)
	}
}
