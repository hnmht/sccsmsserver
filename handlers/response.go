package handlers

import (
	"net/http"
	"sccsmsserver/i18n"

	"github.com/gin-gonic/gin"
)

type ResponseData struct {
	ResKey i18n.ResKey `json:"resKey"`
	Msg    interface{} `json:"msg"`
	Data   interface{} `json:"data"`
}

func ResponseWithMsg(c *gin.Context, key i18n.ResKey, data interface{}) {
	lang := c.GetHeader("Accept-Language")
	msg := key.Msg(lang)
	c.JSON(http.StatusOK, &ResponseData{
		ResKey: key,
		Msg:    msg,
		Data:   data,
	})
}
