package minio

import (
	"context"
	"fmt"
	"path"
	"regexp"
	"strconv"

	"github.com/minio/minio-go"
	log "github.com/sirupsen/logrus"

	"github.com/sevigo/hokan/pkg/core"
)

var bucketNameRegexp = regexp.MustCompile("^[a-zA-Z0-9_.]+$")

const TargetName = "minio"

type minioStore struct {
	client     core.MinioWrapper
	fileStore  core.FileStore
	bucketName string
}

func New(ctx context.Context, fs core.FileStore, conf core.TargetConfig) (core.TargetStorage, error) {
	logger := log.WithFields(log.Fields{
		"target": TargetName,
	})
	if !conf.Active {
		return nil, core.ErrTargetNotActive
	}

	configurator := NewConfigurator()
	if ok, err := configurator.ValidateSettings(conf.Settings); !ok {
		return nil, err
	}

	bucketName := conf.Settings["MINIO_BUCKET_NAME"]
	useSSL, _ := strconv.ParseBool(conf.Settings["MINIO_USE_SSL"])

	logger.Info("Starting new target storage")

	minioClient, err := NewMinioWrapper(&core.MinioConfig{
		Endpoint:        conf.Settings["MINIO_HOST"],
		AccessKeyID:     conf.Settings["MINIO_ACCESS_KEY"],
		SecretAccessKey: conf.Settings["MINIO_SECRET_KEY"],
		UseSSL:          useSSL,
		Bucket:          bucketName,
	})
	if err != nil {
		return nil, err
	}

	return &minioStore{
		bucketName: bucketName,
		client:     minioClient,
		fileStore:  fs,
	}, nil
}

func (s *minioStore) Name() string {
	return TargetName
}

func (s *minioStore) Save(ctx context.Context, result chan core.TargetOperationResult, file *core.File, opt *core.TargetStorageSaveOpt) {
	logger := log.WithFields(log.Fields{
		"target": TargetName,
		"file":   file.Path,
	})

	logger.Debug("saving file")
	objectName := path.Clean(file.Path)
	options := minio.PutObjectOptions{
		UserMetadata: map[string]string{
			"path":      file.Path,
			"size":      fmt.Sprintf("%d", file.Info.Size()),
			"name":      file.Info.Name(),
			"mode-time": file.Info.ModTime().String(),
			"checksum":  file.Checksum,
		},
		// TODO: we can use Progress for the reporting the progress back to the client
	}
	n, err := s.client.FPutObjectWithContext(ctx, s.bucketName, objectName, file.Path, options)
	if err != nil {
		result <- core.TargetOperationResultError(err)
		return
	}
	logger.Infof("Successfully uploaded %s of size %d", objectName, n)
	saveErr := s.fileStore.Save(ctx, TargetName, file)
	result <- core.TargetOperationResultError(saveErr)
}

func (s *minioStore) Restore(ctx context.Context, result chan core.TargetOperationResult, files []*core.File, opt *core.TargetStorageRestoreOpt) {
	result <- core.TargetOperationResultError(core.ErrNotImplemented)
}

func (s *minioStore) Delete(ctx context.Context, result chan core.TargetOperationResult, files []*core.File, opt *core.TargetStorageDeleteOpt) {
	log.WithField("target", TargetName).Print("Delete")
	result <- core.TargetOperationResultError(core.ErrNotImplemented)
}

func (s *minioStore) Ping(ctx context.Context) error {
	log.WithField("target", TargetName).Print("target")
	_, err := s.client.BucketExists(s.bucketName)
	return err
}

func (s *minioStore) Info(ctx context.Context) core.TargetInfo {
	return core.TargetInfo{}
}
