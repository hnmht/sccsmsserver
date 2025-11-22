package handlers

import (
	"sccsmsserver/db/pg"
	"sccsmsserver/i18n"

	"github.com/gin-gonic/gin"
)

// Get Menulist handler
func GetMenuHandler(c *gin.Context) {
	ResponseWithMsg(c, i18n.StatusOK, pg.SysFunctionList)
}
