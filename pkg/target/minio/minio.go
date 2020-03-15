package minio

import (
	"context"
	"errors"
	"path"

	"github.com/minio/minio-go"
	log "github.com/sirupsen/logrus"

	"github.com/sevigo/hokan/pkg/core"
	filestore "github.com/sevigo/hokan/pkg/store/file"
	"github.com/sevigo/hokan/pkg/target/utils"
)

const TargetName = "minio"

const endpointDefault = "192.168.0.141:9000"
const accessKeyIDDefault = "minio"
const secretAccessKeyDefault = "miniostorage"
const useSSL = false

type minioStore struct {
	client     core.MinioWrapper
	fs         core.FileStore
	bucketName string
}

func New(ctx context.Context, fs core.FileStore) (core.TargetStorage, error) {
	bucketName := "osaka"
	minioClient, err := NewMinioWrapper(&core.MinioConfig{
		Endpoint:        endpointDefault,
		AccessKeyID:     accessKeyIDDefault,
		SecretAccessKey: secretAccessKeyDefault,
		Bucket:          bucketName,
		UseSSL:          useSSL,
	})
	if err != nil {
		return nil, err
	}

	return &minioStore{
		bucketName: bucketName,
		client:     minioClient,
		fs:         fs,
	}, nil
}

func (s *minioStore) Save(ctx context.Context, file *core.File) error {
	logger := log.WithFields(log.Fields{
		"target": TargetName,
		"file":   file.Path,
	})

	storedFile, err := s.fs.Find(ctx, TargetName, file.Path)
	if errors.Is(err, filestore.ErrFileEntryNotFound) || utils.FileHasChanged(file, storedFile) {
		logger.Debug("saving file")
		objectName := path.Clean(file.Path)
		options := minio.PutObjectOptions{
			UserMetadata: map[string]string{
				"path":     file.Path,
				"info":     file.Info,
				"checksum": file.Checksum,
			},
			// TODO: we can use Progress for the reporting the progress back to the client
		}
		n, err := s.client.FPutObjectWithContext(ctx, s.bucketName, objectName, file.Path, options)
		if err != nil {
			return err
		}
		logger.Infof("Successfully uploaded %s of size %d\n", objectName, n)
		return s.fs.Save(ctx, TargetName, file)
	}
	return nil
}

func (s *minioStore) List(context.Context) ([]*core.File, error) {
	log.Printf("[minio] list\n")
	return nil, errors.New("not implemented")
}

func (s *minioStore) Find(ctx context.Context, q string) (*core.File, error) {
	log.Printf("[minio] find %q\n", q)
	return nil, errors.New("not implemented")
}

func (s *minioStore) Delete(ctx context.Context, file *core.File) error {
	log.Printf("[minio] save %#v\n", file)
	return errors.New("not implemented")
}

func (s *minioStore) Ping(ctx context.Context) error {
	_, err := s.client.BucketExists(s.bucketName)
	return err
}
