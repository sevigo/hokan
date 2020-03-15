package minio

import (
	"context"

	"github.com/minio/minio-go"
	"github.com/sevigo/hokan/pkg/core"
	log "github.com/sirupsen/logrus"
)

type wrapper struct {
	client *minio.Client
}

func NewMinioWrapper(config *core.MinioConfig) (core.MinioWrapper, error) {
	minioClient, err := minio.New(config.Endpoint, config.AccessKeyID, config.SecretAccessKey, config.UseSSL)
	if err != nil {
		log.WithError(err).Error("Can't create new minio client")
		return nil, err
	}

	// Make a new bucket with mchine name (must be from the config)
	err = minioClient.MakeBucket(config.Bucket, "")
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(config.Bucket)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", config.Bucket)
		} else {
			return nil, err
		}
	} else {
		log.Printf("Successfully created %s\n", config.Bucket)
	}

	return &wrapper{
		client: minioClient,
	}, nil
}

func (w *wrapper) FPutObjectWithContext(ctx context.Context, bucketName, objectName, filePath string, opts minio.PutObjectOptions) (int64, error) {
	return w.client.FPutObjectWithContext(ctx, bucketName, objectName, filePath, opts)
}

func (w *wrapper) BucketExists(bucketName string) (bool, error) {
	return w.client.BucketExists(bucketName)
}
