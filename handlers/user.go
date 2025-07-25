package handlers

import (
	"sccsmsserver/db/pg"
	"sccsmsserver/i18n"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LoginHandler 用户登录
func LoginHandler(c *gin.Context) {
	//获取请求参数及校验
	p := new(pg.ParamLogin)
	if err := c.ShouldBind(p); err != nil {
		//请求参数有误，记录日志并返回响应
		zap.L().Error("Login with invalid param", zap.Error(err))

		ResponseWithMsg(c, i18n.CodeInvalidParm, nil)
		return
	}
	//补充ParamLogin内容
	p.ClientIP = c.ClientIP()
	p.ClientType = c.Request.Header.Get("XClientType")
	p.UserAgent = c.Request.UserAgent()

	//登录
	resStatus, token, _ := pg.Login(p)

	//返回响应
	ResponseWithMsg(c, resStatus, token)
}
