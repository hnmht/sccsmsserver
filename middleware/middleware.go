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

// JWT authentication middleware
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		clientType := c.Request.Header.Get("XClientType")

		if authHeader == "" {
			handlers.ResponseWithMsg(c, i18n.CodeNeedLogin, nil)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			handlers.ResponseWithMsg(c, i18n.CodeInvalidToken, nil)
			c.Abort()
			return
		}

		mc, resStaus := jwt.ParseToken(parts[1])
		if resStaus != i18n.StatusOK {
			handlers.ResponseWithMsg(c, resStaus, nil)
			c.Abort()
			return
		}

		// Get Current User from cache
		var ou pg.OnlineUser
		ou.User.ID = mc.UserID
		ou.ClientType = clientType
		exist, _, err := ou.Get()
		if err != nil {
			handlers.ResponseWithMsg(c, i18n.CodeInternalError, nil)
			c.Abort()
			return
		}

		// If the user token is not found in the cache,
		// it indicate that the token has been invalidated bye the administrator.
		if exist == 0 {
			handlers.ResponseWithMsg(c, i18n.CodeTokenDestroy, nil)
			c.Abort()
			return
		}

		// If the token exists but the ID is mismatched,
		// it means the user has already logged in another device
		if ou.TokenID != mc.Id {
			handlers.ResponseWithMsg(c, i18n.CodeLoginOther, nil)
			c.Abort()
			return
		}

		// Save the current request's user information to the context
		c.Set(pub.CTXUserCode, mc.UserCode)
		c.Set(pub.CTXUserID, mc.UserID)
		c.Set(pub.CTXTokenID, mc.Id)
		c.Next()
	}
}
