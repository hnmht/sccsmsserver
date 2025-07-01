package main

import (
	"fmt"
	"sccsmsserver/cache"
	"sccsmsserver/db/pg"
	"sccsmsserver/logger"
	"sccsmsserver/pkg/minio"
	"sccsmsserver/pkg/mysf"
	"sccsmsserver/setting"

	"go.uber.org/zap"
)

func main() {
	// step 1: Read global configuration
	if err := setting.Init(); err != nil {
		fmt.Println("Read global configuration failed:", err)
		return
	}

	// step 2: Logger component initialization
	if err := logger.Init(setting.Conf.LogConfig, setting.Conf.Mode); err != nil {
		fmt.Println("logger component initialization failed:", err)
		return
	}
	defer zap.L().Sync()

	// step 3: Snowflake ID Generator initialization
	if err := mysf.Init(setting.Conf.StartTime, setting.Conf.MachineID); err != nil {
		zap.L().Error("snowflake Generator component initialization failed:", zap.Error(err))
		return
	}
	// step 4: Database connection initialization
	if err := pg.Init(setting.Conf.PqConfig); err != nil {
		zap.L().Error("Database connection initialization failed:", zap.Error(err))
		return
	}
	defer pg.Close()
	// step 5: Cache component initialization
	if err := cache.Init(setting.Conf.RedisConfig.Enabled); err != nil {
		zap.L().Error("Cache component initialization failed:", zap.Error(err))
		return
	}
	defer cache.Close()

	// setp 6: MINIO client initialization
	if err := minio.Init(setting.Conf.MinioConfig.Endpoint, setting.Conf.MinioConfig.AccessKeyID, setting.Conf.MinioConfig.SecretAccessKey,
		setting.Conf.MinioConfig.Secure, setting.Conf.SelfSigned, setting.Conf.MinioConfig.DefaultBucket); err != nil {
		zap.L().Error("minio Init failed:", zap.Error(err))
		return
	}
}
