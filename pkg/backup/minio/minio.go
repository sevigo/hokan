package minio

import (
	"context"
	"fmt"
	"path"
	"regexp"

	"github.com/minio/minio-go"
	log "github.com/sirupsen/logrus"

	"github.com/sevigo/hokan/pkg/core"
)

var bucketNameRegexp = regexp.MustCompile("^[a-zA-Z0-9_.]+$")

const name = "minio"

type minioStore struct {
	client     core.MinioWrapper
	fileStore  core.FileStore
	bucketName string
}

func New(ctx context.Context, fs core.FileStore, conf *core.BackupOptions) (core.Backup, error) {
	log.WithFields(log.Fields{
		"backup": name,
	}).Info("Starting new target storage")

	minioClient, err := NewMinioWrapper(&core.MinioConfig{
		Endpoint:        conf.TargetURL,
		AccessKeyID:     conf.AccessKeyID,
		SecretAccessKey: conf.SecretAccessKey,
		UseSSL:          conf.UseSSL,
		Bucket:          name,
	})
	if err != nil {
		return nil, err
	}

	return &minioStore{
		bucketName: name,
		client:     minioClient,
		fileStore:  fs,
	}, nil
}

func (s *minioStore) Name() string {
	return name
}

func (s *minioStore) Save(ctx context.Context, result chan core.BackupResult, file *core.File, opt *core.BackupOperationOptions) {
	logger := log.WithFields(log.Fields{
		"backup": name,
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
		result <- core.BackupResult{
			Success: false,
			Error:   err,
			Message: fmt.Sprintf("can't save file %q to %q", file.Path, s.bucketName),
		}
		return
	}
	logger.Infof("Successfully uploaded %s of size %d", objectName, n)
	saveErr := s.fileStore.Save(ctx, name, file)
	if saveErr != nil {
		result <- core.BackupResult{
			Success: false,
			Error:   saveErr,
			Message: fmt.Sprintf("can't save backup info for the file %q", file.Path),
		}
	} else {
		result <- core.BackupResult{
			Success: true,
			Message: core.BackupSuccessMessage,
		}
	}
}

func (s *minioStore) Restore(ctx context.Context, result chan core.BackupResult, files []*core.File, opt *core.BackupOperationOptions) {
	result <- core.BackupResult{
		Success: false,
		Error:   core.ErrNotImplemented,
	}
}

func (s *minioStore) Delete(ctx context.Context, result chan core.BackupResult, files []*core.File, opt *core.BackupOperationOptions) {
	log.WithField("backup", name).Print("Delete")
	result <- core.BackupResult{
		Success: false,
		Error:   core.ErrNotImplemented,
	}
}

func (s *minioStore) Ping(ctx context.Context) error {
	log.WithField("backup", name).Print("target")
	_, err := s.client.BucketExists(s.bucketName)
	return err
}
