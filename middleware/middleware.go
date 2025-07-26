package middleware

import (
	"sccsmsserver/db/pg"
	"sccsmsserver/handlers"
	"sccsmsserver/i18n"
	"sccsmsserver/pkg/jwt"
	"sccsmsserver/pub"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Check if the clientType is in the list of valid values
func clientTypeValid(clientType string) (valid bool) {
	valid = false
	for _, v := range pub.ValidClientTypes {
		if v == clientType {
			valid = true
			return
		}
	}
	return
}

// CheckClientTypeMiddleware Check the client type.
// This aims to reduce server resources consumed by network hack scans.
func CheckClientTypeMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// Frontend requests require the "XClientType" custom header.
		clientType := c.Request.Header.Get("XClientType")
		// Verify if the request includes the custom “XClientType” header.
		if clientType == "" {
			handlers.ResponseWithMsg(c, i18n.CodeClientEmpty, nil)
			c.Abort()
			return
		}

		if !clientTypeValid(clientType) {
			zap.L().Info("CheckClientTypeMiddleware  ClientType invalid")
			handlers.ResponseWithMsg(c, i18n.CodeClientUnknown, nil)
			c.Abort()
			return
		}
		c.Set(pub.CTXClientType, clientType)

		c.Next()
	}
}

// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// startTime := time.Now()
		//约定请求JWTToken放在请求Header的Authorization中，并使用Bearer开头 : Authorization: Bearer ****
		authHeader := c.Request.Header.Get("Authorization")
		clientType := c.Request.Header.Get("XClientType")

		if authHeader == "" {
			handlers.ResponseWithMsg(c, i18n.CodeNeedLogin, nil)
			c.Abort()
			return
		}

		// 判断格式：解析请求头按照空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			handlers.ResponseWithMsg(c, i18n.CodeInvalidToken, nil)
			c.Abort()
			return
		}

		//parts[1]是获取到的tokenString, 使用JWT解析函数进行解析
		mc, resStaus := jwt.ParseToken(parts[1])
		if resStaus != i18n.StatusOK {
			handlers.ResponseWithMsg(c, resStaus, nil)
			c.Abort()
			return
		}

		//从缓存中获取当前用户
		var ou pg.OnlineUser
		ou.User.ID = mc.UserID
		ou.ClientType = clientType
		exist, _, err := ou.Get()
		if err != nil {
			handlers.ResponseWithMsg(c, i18n.CodeInternalError, nil)
			c.Abort()
			return
		}

		//如果token没有存在于在线用户缓存中，说明用户被管理员踢出系统
		if exist == 0 {
			handlers.ResponseWithMsg(c, i18n.CodeTokenDestroy, nil)
			c.Abort()
			return
		}

		//如果token存在但是和tokenID不一致，说明用户已经在其他终端登录
		if ou.TokenID != mc.Id {
			handlers.ResponseWithMsg(c, i18n.CodeLoginOther, nil)
			c.Abort()
			return
		}

		//将当前请求的usrename信息保存到请求上下文
		c.Set(pub.CTXUserCode, mc.UserCode)
		c.Set(pub.CTXUserID, mc.UserID)
		c.Set(pub.CTXTokenID, mc.Id)
		c.Next() //后续的处理请求函数中,可以通过c.GET(CTXUsername和CTXUserID)获取当前用户信息
	}
}
