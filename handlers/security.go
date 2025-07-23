package handlers

import (
	"sccsmsserver/i18n"
	"sccsmsserver/pkg/security"

	"github.com/gin-gonic/gin"
)

// Publish public key
func GetPublicKeyHandler(c *gin.Context) {
	publicKey := security.ScRsa.GetPublicKey()
	ResponseWithMsg(c, i18n.CodeSuccess, publicKey)
}
