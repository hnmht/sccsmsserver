package localcache

import (
	"context"
	"fmt"
	"time"

	"github.com/allegro/bigcache/v3"
	"go.uber.org/zap"
)

var localCache *bigcache.BigCache
var ctx = context.Background()

// Initialize local cache
func Init(durtion time.Duration) (err error) {
	localCache, err = bigcache.New(ctx, bigcache.DefaultConfig(durtion))
	if err != nil {
		return
	}
	zap.L().Info("Local cache component initialized successfully.")
	return
}

// Close local cache
func Close() {
	_ = localCache.Close()
}

// Set a value in the local cache
func Set(key string, v []byte) (err error) {
	err = localCache.Set(key, v)
	if err != nil {
		msg := fmt.Sprintf("%s%s", key, "Set LocalCache.set falied:")
		zap.L().Error(msg, zap.Error(err))
	}
	return
}

// Get a value from the local cache
func Get(key string) (exist int32, v []byte, err error) {
	v, err = localCache.Get(key)
	if err != nil {
		exist = 0
		if err == bigcache.ErrEntryNotFound {
			return exist, v, nil
		}
		msg := fmt.Sprintf("%s%s", key, "Get LocalCache.get falied: ")
		zap.L().Error(msg, zap.Error(err))
		return
	}
	exist = 1
	return
}

// Delete a value from the local cache
func Del(key string) (err error) {
	err = localCache.Delete(key)
	if err != nil {
		msg := fmt.Sprintf("%s%s", key, " Del localCache.Del failed: ")
		zap.L().Error(msg, zap.Error(err))
	}
	return
}
