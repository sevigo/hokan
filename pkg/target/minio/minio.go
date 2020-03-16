package minio

import (
	"context"
	"errors"
	"path"
	"strconv"

	"github.com/minio/minio-go"
	log "github.com/sirupsen/logrus"

	"github.com/sevigo/hokan/pkg/core"
	filestore "github.com/sevigo/hokan/pkg/store/file"
	"github.com/sevigo/hokan/pkg/target/utils"
)

const TargetName = "minio"

type minioStore struct {
	client     core.MinioWrapper
	fileStore  core.FileStore
	bucketName string
}

func DefaultConfig() *core.TargetConfig {
	return &core.TargetConfig{
		Active:      false,
		Name:        "minio",
		Description: "open source cloud object storage server compatible with Amazon S3",
		Settings: map[string]string{
			"MINIO_HOST":        "",
			"MINIO_ACCESS_KEY":  "",
			"MINIO_SECRET_KEY":  "",
			"MINIO_USE_SSL":     "",
			"MINIO_BUCKET_NAME": "",
		},
	}
}

func New(ctx context.Context, fs core.FileStore, conf core.TargetConfig) (core.TargetStorage, error) {
	logger := log.WithFields(log.Fields{
		"target": TargetName,
	})
	if !conf.Active {
		return nil, core.ErrTargetNotActive
	}

	bucketName := conf.Settings["MINIO_BUCKET_NAME"]
	useSSL, err := strconv.ParseBool(conf.Settings["MINIO_USE_SSL"])
	if err != nil {
		logger.WithError(err).Errorf("can't convert the value of MINIO_USE_SSL=%q to bool", conf.Settings["MINIO_USE_SSL"])
		useSSL = false
	}

	log.WithFields(log.Fields{
		"target": TargetName,
	}).Info("Starting new target storage")

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

func (s *minioStore) Save(ctx context.Context, file *core.File) error {
	logger := log.WithFields(log.Fields{
		"target": TargetName,
		"file":   file.Path,
	})

	storedFile, err := s.fileStore.Find(ctx, TargetName, file.Path)
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
		return s.fileStore.Save(ctx, TargetName, file)
	}
	logger.Info("the file has not changed")
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
