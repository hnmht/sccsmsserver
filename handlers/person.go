package handlers

import (
	"sccsmsserver/db/pg"
	"sccsmsserver/i18n"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Get Person list Master data handler
func GetPersonsHandler(c *gin.Context) {
	persons, resStatus, _ := pg.GetPersons()
	ResponseWithMsg(c, resStatus, persons)
}

// Get latest Person master data for front-end cache
func GetPersonsCacheHandler(c *gin.Context) {
	pc := new(pg.PersonCache)
	err := c.ShouldBind(pc)
	if err != nil {
		zap.L().Error("GetPersonsCacheHandler invalid param", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, nil)
		return
	}
	resStatus, _ := pc.GetLatestPersons()
	ResponseWithMsg(c, resStatus, pc)
}
