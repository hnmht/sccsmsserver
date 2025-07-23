package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sccsmsserver/cache"
	"sccsmsserver/db/pg"
	"sccsmsserver/logger"
	"sccsmsserver/pkg/environment"
	"sccsmsserver/pkg/minio"
	"sccsmsserver/pkg/mysf"
	"sccsmsserver/route"
	"sccsmsserver/setting"
	"strconv"
	"syscall"
	"time"

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

	//9.注册路由
	r := route.Setup(setting.Conf.Mode)

	//10.配置服务
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", setting.Conf.Port),
		Handler: r,
	}

	var isTLS = setting.Conf.TLS
	if setting.Conf.TLS {
		var certExist bool = true
		var keyExist bool = true
		//检查证书文件是否存在
		_, errCert := os.Lstat(setting.Conf.CertificateFile)
		if errCert != nil {
			certExist = false
			zap.L().Error("get CertificateFile failed", zap.Error(errCert))
		}
		//检查私玥文件是否存在
		_, errKey := os.Lstat(setting.Conf.PrivateKeyFile)
		if errKey != nil {
			keyExist = false
			zap.L().Error("get PrivateKeyFile failed", zap.Error(errKey))
		}

		if !certExist || !keyExist {
			isTLS = false
			zap.L().Info("没有找到证书文件或私钥文件,系统将以http方式提供服务!")
		}
	}

	ipList := environment.GetLocalIPs()
	var protocol string = "http"
	if isTLS {
		protocol = "https"
	}
	portStr := strconv.FormatInt(int64(setting.Conf.Port), 10)

	if setting.Conf.Mode == "release" {
		fmt.Println("...............................................................")
		fmt.Println(" ")
		fmt.Println("SeaCloud现场管理系统")
		fmt.Println("在浏览器中输入以下地址访问系统:")
		for _, ip := range ipList {
			if ip.To4() != nil {
				fmt.Println(protocol + "://" + ip.String() + ":" + portStr)
			}
		}
		fmt.Println(" ")
		fmt.Println("...............................................................")
		fmt.Println(" ")
		fmt.Println("周口市海云信息技术有限公司")
		fmt.Println("https://www.zkseacloud.cn")
		fmt.Println(" ")
		fmt.Println("...............................................................")
		fmt.Println(" ")
		fmt.Println("SeaCloud 现场管理系统正在运行,请勿关闭此窗口...")
		fmt.Println("SeaCloud Scene Management System Backend Services running, Don't close this window...")
		fmt.Println(" ")
		fmt.Println("...............................................................")
	}

	go func() {
		var err error
		if isTLS {
			err = srv.ListenAndServeTLS(setting.Conf.CertificateFile, setting.Conf.PrivateKeyFile)
		} else {
			err = srv.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			zap.L().Error("srv.ListenAndServe failed", zap.Error(err))
		}
	}()
	//11.关机
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.L().Info("Shutdown Server...")
	//创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//等待5秒，将未处理完的请求处理完成再关闭服务，超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Error("Server Shutdown failed:", zap.Error(err))
	}
	zap.L().Info("Server exiting")
}
