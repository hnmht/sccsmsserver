package main

import (
	"fmt"
	"sccsmsserver/db/cache"
	"sccsmsserver/db/pg"
	"sccsmsserver/logger"
	"sccsmsserver/pkg/minio"
	"sccsmsserver/pkg/mysf"
	"sccsmsserver/setting"

	"go.uber.org/zap"
)

func main() {
	// step 1: Global variable initialzation
	if err := setting.Init(); err != nil {
		fmt.Println("ConfigInit settings failed:", err)
		return
	}

	// step 2: Initialize system log
	if err := logger.Init(setting.Conf.LogConfig, setting.Conf.Mode); err != nil {
		fmt.Println("logger Init failed:", err)
		return
	}
	defer zap.L().Sync()
	zap.L().Info("logger Init success...")

	// step 3: Initialize Snowflake ID Generator.
	if err := mysf.Init(setting.Conf.StartTime, setting.Conf.MachineID); err != nil {
		zap.L().Error("snowflake Init failed:", zap.Error(err))
		return
	}
	// step 4: Initialize Database connection
	if err := pg.Init(setting.Conf.PqConfig); err != nil {
		zap.L().Error("init postgresql failed:", zap.Error(err))
		return
	}
	defer pg.Close()
	// step 5: Initialize cache
	if err := cache.Init(setting.Conf.RedisConfig.Enabled); err != nil {
		zap.L().Error("cache Init failed:", zap.Error(err))
		return
	}
	defer cache.Close()

	// setp 6: Initialize MINIO client
	if err := minio.Init(setting.Conf.MinioConfig.Endpoint, setting.Conf.MinioConfig.AccessKeyID, setting.Conf.MinioConfig.SecretAccessKey,
		setting.Conf.MinioConfig.Secure, setting.Conf.SelfSigned, setting.Conf.MinioConfig.DefaultBucket); err != nil {
		zap.L().Error("minio Init failed:", zap.Error(err))
		return
	}
}
