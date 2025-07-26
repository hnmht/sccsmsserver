package handlers

import (
	"sccsmsserver/db/pg"
	"sccsmsserver/i18n"

	"github.com/gin-gonic/gin"
)

// Publish system public information
func PubSystemInformationHandler(c *gin.Context) {
	ResponseWithMsg(c, i18n.StatusOK, pg.ServerPubInfo)
}
