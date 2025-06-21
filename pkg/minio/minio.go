package minio

import (
	"crypto/tls"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/minio/minio-go"
	_ "github.com/minio/minio-go/pkg/encrypt"
	"go.uber.org/zap"
)

// Minio Client
var minioClient *minio.Client

// Bucket region
var location = "cn-north-1"

// Current Bucket
var currentBucket string

// Initialize Minio client
func Init(endpoint string, accessKeyID string, secretAccessKey string, secure bool, selfSigned bool, defaultBucket string) (err error) {
	minioClient, err = minio.New(endpoint, accessKeyID, secretAccessKey, secure)
	if err != nil {
		zap.L().Error("minio Init failed:", zap.Error(err))
		return
	}

	// Support Minio server with self-signed certificate
	if secure && selfSigned {
		minioClient.SetCustomTransport(&http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}})
	}

	// Check if the default bucket exists
	found, err := minioClient.BucketExists(defaultBucket)
	if err != nil {
		zap.L().Error("minio BucketExists failed:", zap.Error(err))
		return
	}

	// Create a bucket if the default bucket does not exist
	if !found {
		err = minioClient.MakeBucket(defaultBucket, location)
		if err != nil {
			zap.L().Error("minio MakeBucket failed:", zap.Error(err))
			return
		}
	}
	// Assign a value to currentBucket
	currentBucket = defaultBucket
	zap.L().Info("minio Init success...")
	return nil
}

// Upload a file to the minio server
func UploadFile(objectName string, reader io.Reader, objectSize int64) (ok bool, err error) {
	_, err = minioClient.PutObject(currentBucket, objectName, reader, objectSize, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		zap.L().Error("minio UploadFile failed:", zap.Error(err))
		return false, err
	}

	return true, nil
}

// Get the file access URL from the Minio server
func GetFileUrl(fileName string, expirs time.Duration) (fileurl string, err error) {
	reqParams := make(url.Values)
	presignedURL, err := minioClient.PresignedGetObject(currentBucket, fileName, expirs, reqParams)
	if err != nil {
		zap.L().Error("minio GetFileUrl failed:", zap.Error(err))
		return "", err
	}
	fileurl = presignedURL.String()
	return
}
