package rediscache

import (
	"context"
	"fmt"
	"sccsmsserver/setting"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var ctx = context.Background()
var rdb *redis.Client

// Initialize Redis cache
func Init(cfg *setting.RedisConfig) (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.Db,
		PoolSize: cfg.PoolSize,
	})

	_, err = rdb.Ping(ctx).Result()

	if err != nil {
		zap.L().Error("redis ping failed", zap.Error(err))
		return
	}
	zap.L().Info("redis cache connected successful.")
	return
}

// Close Redis connection
func Close() {
	_ = rdb.Close()
}

// Set a value in the Redis cache
func Set(key string, v []byte, durtion time.Duration) (err error) {
	p := string(v)
	err = rdb.SetEx(ctx, key, p, durtion).Err()
	if err != nil {
		msg := fmt.Sprintf("%s%s", key, " Set redis rdb.SetEx falied:")
		zap.L().Error(msg, zap.Error(err))
	}

	return
}

// Get a value from the Redis cache
func Get(key string) (exist int32, v []byte, err error) {
	p, err := rdb.Get(ctx, key).Result()
	if err != nil {
		exist = 0
		p = ""
		if err == redis.Nil {
			v = []byte(p)
			return exist, v, nil
		}
		msg := fmt.Sprintf("%s%s", key, " Get redis rdb.Get falied: ")
		zap.L().Error(msg, zap.Error(err))
		return
	}
	v = []byte(p)
	exist = 1
	return
}

// Delete a value from the Redis cache
func Del(key string) (err error) {
	_, err = rdb.Del(ctx, key).Result()
	if err != nil {
		msg := fmt.Sprintf("%s%s", key, " Del redis rdb.Del failed: ")
		zap.L().Error(msg, zap.Error(err))
	}
	return
}
