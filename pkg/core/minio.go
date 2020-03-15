package core

import (
	"context"

	"github.com/minio/minio-go"
)

type MinioConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	Bucket          string
	UseSSL          bool
}

type MinioWrapper interface {
	FPutObjectWithContext(ctx context.Context, bucketName, objectName, filePath string, opts minio.PutObjectOptions) (int64, error)
	BucketExists(bucketName string) (bool, error)
}
