package handlers

import (
	"sccsmsserver/db/pg"
	"sccsmsserver/i18n"

	"github.com/gin-gonic/gin"
)

// GetMenuHandler 获取树状菜单列表
func GetMenuHandler(c *gin.Context) {
	ResponseWithMsg(c, i18n.StatusOK, pg.SysFunctionList)
}
