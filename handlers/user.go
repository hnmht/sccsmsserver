package handlers

import (
	"sccsmsserver/db/pg"
	"sccsmsserver/i18n"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// User login handler
func LoginHandler(c *gin.Context) {
	// Step 1: Get request parameters
	p := new(pg.ParamLogin)
	if err := c.ShouldBind(p); err != nil {
		zap.L().Error("Login with invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, nil)
		return
	}
	// Provde additional login information
	p.ClientIP = c.ClientIP()
	p.ClientType = c.Request.Header.Get("XClientType")
	p.UserAgent = c.Request.UserAgent()

	// User login validation
	resStatus, token, _ := pg.Login(p)

	// Respond to client request
	ResponseWithMsg(c, resStatus, token)
}
