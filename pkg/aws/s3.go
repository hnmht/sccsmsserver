package aws

import (
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

var (
	s3Client      *s3.Client
	presignClient *s3.PresignClient
	currentBucket string
	location      = "cn-north-1"
)

// Init initializes the S3-compatible client
func Init(endpoint string, accessKey string, secretKey string, secure bool, selfSigned bool, defaultBucket string) error {
	// Build HTTP transport with optional insecure TLS
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	if secure && selfSigned {
		customTransport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}
	// S3 client without deprecated global endpoint resolver
	s3Client = s3.New(s3.Options{
		Region:       location,
		Credentials:  credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		BaseEndpoint: aws.String(endpoint), // Recommended way (MinIO/RustFS/AWS)
		UsePathStyle: true,                 // Required for MinIO/RustFS
		HTTPClient: &http.Client{
			Transport: customTransport,
		},
	})
	presignClient = s3.NewPresignClient(s3Client)

	// Check if bucket exists
	_, err := s3Client.HeadBucket(context.TODO(), &s3.HeadBucketInput{
		Bucket: aws.String(defaultBucket),
	})
	if err != nil {
		// Create bucket
		_, err = s3Client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
			Bucket: aws.String(defaultBucket),
			CreateBucketConfiguration: &s3types.CreateBucketConfiguration{
				LocationConstraint: s3types.BucketLocationConstraint(location),
			},
		})
		if err != nil {
			zap.L().Error("S3 Init CreateBucket failed:", zap.Error(err))
			return err
		}
	}
	currentBucket = defaultBucket

	zap.L().Info("S3 client initialized successfully.")
	return nil
}

// UploadFile uploads an object to S3
func UploadFile(objectName string, reader io.Reader, objectSize int64) (bool, error) {
	_, err := s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:        aws.String(currentBucket),
		Key:           aws.String(objectName),
		Body:          reader,
		ContentLength: aws.Int64(objectSize),
		ContentType:   aws.String("application/octet-stream"),
	})
	if err != nil {
		zap.L().Error("S3 UploadFile PutObject failed:", zap.Error(err))
		return false, err
	}
	return true, nil
}

// GetFileUrl returns a presigned download URL for the file
func GetFileUrl(fileName string, expires time.Duration) (string, error) {
	presignedReq, err := presignClient.PresignGetObject(
		context.TODO(),
		&s3.GetObjectInput{
			Bucket: aws.String(currentBucket),
			Key:    aws.String(fileName),
		},
		func(o *s3.PresignOptions) {
			o.Expires = expires
		},
	)

	if err != nil {
		zap.L().Error("S3 GetFileUrl PresignGetObject failed:", zap.Error(err))
		return "", err
	}

	return presignedReq.URL, nil
}
