package cache

import (
	"fmt"
	"sccsmsserver/db/localcache"
	"sccsmsserver/db/rediscache"
	"sccsmsserver/pub"
	"sccsmsserver/setting"
)

// Whether to use Redis cache
var redisEnabled = false

// Initalize cache
func Init(enabled bool) (err error) {
	redisEnabled = enabled
	if redisEnabled {
		err = rediscache.Init(setting.Conf.RedisConfig)
		return
	}
	err = localcache.Init(pub.CacheExpiration)
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

// Set Archive cache
func Set(docType pub.DocType, id int32, v []byte) (err error) {
	key := fmt.Sprintf("%s%s%d", docType, ":", id)
	if redisEnabled {
		err = rediscache.Set(key, v, pub.CacheExpiration)
		return
	}
	err = localcache.Set(key, v)
	return
}

// Get Archive cache
func Get(docType pub.DocType, id int32) (exist int32, v []byte, err error) {
	key := fmt.Sprintf("%s%s%d", docType, ":", id)
	if redisEnabled {
		exist, v, err = rediscache.Get(key)
		return
	}
	exist, v, err = localcache.Get(key)
	return
}

// Del Archive cache
func Del(docType pub.DocType, id int32) (err error) {
	key := fmt.Sprintf("%s%s%d", docType, ":", id)
	if redisEnabled {
		err = rediscache.Del(key)
		return
	}
	err = localcache.Del(key)
	return
}

// Set Other cache
func SetOther(key string, v []byte) (err error) {
	if redisEnabled {
		err = rediscache.Set(key, v, pub.CacheExpiration)
		return
	}
	err = localcache.Set(key, v)
	return
}

// Get Other cache
func GetOther(key string) (exist int32, v []byte, err error) {
	if redisEnabled {
		exist, v, err = rediscache.Get(key)
		return
	}
	exist, v, err = localcache.Get(key)
	return
}

// Del Other cache
func DelOther(key string) (err error) {
	if redisEnabled {
		err = rediscache.Del(key)
		return
	}
	err = localcache.Del(key)
	return
}
