package cache

import (
	"sccsmsserver/db/localcache"
	"sccsmsserver/db/rediscache"
	"sccsmsserver/setting"
	"time"
)

// cache expiration time
const durtion = 2 * time.Hour

// Whether to use Redis cache
var redisEnabled = false

// Initalize cache
func Init(enabled bool) (err error) {
	redisEnabled = enabled
	if redisEnabled {
		err = rediscache.Init(setting.Conf.RedisConfig)
		return
	}
	err = localcache.Init(durtion)
	return
}

// Close cache
func Close() {
	if redisEnabled {
		rediscache.Close()
	} else {
		localcache.Close()
	}
}
