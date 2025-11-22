package route

import (
	"encoding/json"
	"fmt"
	"math"
	"sccsmsserver/cache"
	"sccsmsserver/db/pg"
	"sccsmsserver/handlers"
	"sccsmsserver/i18n"
	"sccsmsserver/pub"
	"sccsmsserver/setting"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// IP Address Blacklist
func IpBlackListMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the client's IP address
		clientIp := c.ClientIP()
		// Check if the IP address is on the blacklist
		key := fmt.Sprintf("%s%s%s", pub.IPBlack, ":", clientIp)
		exist, v, err := cache.GetOther(key)
		if err != nil {
			zap.L().Error("IpBlackListMiddleWare cache.GetOther failed:", zap.Error(err))
			handlers.ResponseWithMsg(c, i18n.StatusInternalError, nil)
			c.Abort()
			return
		}
		// If the IP Address is in the blacklist
		if exist == 1 {
			var ipLock pg.IpLock
			err1 := json.Unmarshal(v, &ipLock)
			if err1 != nil {
				zap.L().Error("IpBlackListMiddleWare json.Unmarshal failed:", zap.Error(err))
				handlers.ResponseWithMsg(c, i18n.StatusInternalError, nil)
				c.Abort()
				return
			}
			// Check if the IP lock has expired
			lockedTime := ipLock.StartTime.Add(time.Minute * time.Duration(setting.Conf.IpLockedMinutes))
			intervalMinute := int32(math.Ceil(time.Since(lockedTime).Minutes()))
			if intervalMinute < 0 {
				handlers.ResponseWithMsg(c, i18n.StatusResReject, nil, intervalMinute*-1)
				c.Abort()
				return
			}
			// Delete IP from cache
			_ = cache.DelOther(key)
		}
		c.Next()
	}
}
