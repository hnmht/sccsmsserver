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
	"sccsmsserver/pkg/aws"
	"sccsmsserver/pkg/environment"
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

	// setp 6: aws s3 client initialization
	if err := aws.Init(setting.Conf.S3Storage.Endpoint, setting.Conf.S3Storage.AccessKeyID, setting.Conf.S3Storage.SecretAccessKey,
		setting.Conf.S3Storage.Secure, setting.Conf.SelfSigned, setting.Conf.S3Storage.DefaultBucket, setting.Conf.S3Storage.Location); err != nil {
		zap.L().Error("S3 Object Storage Init failed:", zap.Error(err))
		return
	}

	// Step 7: Route Setup
	r := route.Setup(setting.Conf.Mode)

	// Step 8: Start HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", setting.Conf.Port),
		Handler: r,
	}
	var isTLS = setting.Conf.TLS
	if setting.Conf.TLS {
		var certExist bool = true
		var keyExist bool = true
		// Check if the certificate file exists
		_, errCert := os.Lstat(setting.Conf.CertificateFile)
		if errCert != nil {
			certExist = false
			zap.L().Error("get CertificateFile failed", zap.Error(errCert))
		}
		// Check if the private key file exists
		_, errKey := os.Lstat(setting.Conf.PrivateKeyFile)
		if errKey != nil {
			keyExist = false
			zap.L().Error("get PrivateKeyFile failed", zap.Error(errKey))
		}

		if !certExist || !keyExist {
			isTLS = false
			zap.L().Info("The certificate or private key file was not found. the system will serve using HTTP.")
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
		fmt.Println("Sea&Cloud Consruction Site Management System")
		fmt.Println("Enter the following address in your browser to access the system:")
		for _, ip := range ipList {
			if ip.To4() != nil {
				fmt.Println(protocol + "://" + ip.String() + ":" + portStr)
			}
		}
		fmt.Println(" ")
		fmt.Println("...............................................................")
		fmt.Println(" ")
		fmt.Println("Author: Haitao Meng")
		fmt.Println("https://github.com/hnmht")
		fmt.Println(" ")
		fmt.Println("...............................................................")
		fmt.Println(" ")
		fmt.Println("Sea&Cloud Construction Site Management System Backend Services running, Don't close this window...")
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
	// Shut down the service
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.L().Info("Shutdown Server...")
	// Create a context with a 5-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Wait for 5 seconds to complete pending requests before shutting down the service.
	// Exit withe a timeout if it exceeds 5 seconds.
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Error("Server Shutdown failed:", zap.Error(err))
	}
	zap.L().Info("Server exiting")
}
