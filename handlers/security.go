package handlers

import (
	"sccsmsserver/i18n"
	"sccsmsserver/pkg/jwt"
	"sccsmsserver/pkg/security"
	"sccsmsserver/pub"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Publish public key
func GetPublicKeyHandler(c *gin.Context) {
	publicKey := security.ScRsa.GetPublicKey()
	ResponseWithMsg(c, i18n.StatusOK, publicKey)
}

// Check Token
func ValidateToken(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")

	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		ResponseWithMsg(c, i18n.CodeInvalidToken, false)
	}
	// Parse Token
	mc, resStaus := jwt.ParseToken(parts[1])
	if resStaus != i18n.StatusOK {
		ResponseWithMsg(c, resStaus, false)
	}

	expireAt := time.Unix(mc.ExpiresAt, 0)
	d := time.Since(expireAt)
	if (d.Seconds() + pub.TokenAboutToExpirtSeconds) > 0 {
		ResponseWithMsg(c, i18n.CodeAboutToExpireToken, false)
	}

	ResponseWithMsg(c, i18n.StatusOK, true)
}
