package main

import (
	"fmt"
	"sccsmsserver/logger"
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
		fmt.Println("init logger failed:", err)
		return
	}
	defer zap.L().Sync()
	zap.L().Info("logger init success...")

	// step 3: Initialize Snowflake ID Generator.
	if err := mysf.Init(setting.Conf.StartTime, setting.Conf.MachineID); err != nil {
		zap.L().Error("snowflake Init failed:", zap.Error(err))
		return
	}

}
