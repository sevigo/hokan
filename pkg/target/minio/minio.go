package minio

import (
	"context"
	"errors"
	"fmt"
	"path"
	"regexp"
	"strconv"

	"github.com/minio/minio-go"
	log "github.com/sirupsen/logrus"

	"github.com/sevigo/hokan/pkg/core"
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
		Settings: core.TargetSettings{
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

	s := &minioStore{}
	if ok, err := s.ValidateSettings(conf.Settings); !ok {
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

	s.bucketName = bucketName
	s.client = minioClient
	s.fileStore = fs

	return s, nil
}

func (s *minioStore) Save(ctx context.Context, file *core.File, opt *core.TargetStorageSaveOpt) error {
	logger := log.WithFields(log.Fields{
		"target": TargetName,
		"file":   file.Path,
	})

	storedFile, err := s.fileStore.Find(ctx, &core.FileSearchOptions{
		FilePath:   file.Path,
		TargetName: TargetName,
	})
	if errors.Is(err, core.ErrFileNotFound) || utils.FileHasChanged(file, storedFile) {
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
			return err
		}
		logger.Infof("Successfully uploaded %s of size %d", objectName, n)
		return s.fileStore.Save(ctx, TargetName, file)
	}
	logger.Info("file hasn't change")
	return nil
}

func (s *minioStore) List(ctx context.Context, opt *core.TargetStorageListOpt) ([]*core.File, error) {
	log.Printf("[minio] list")
	return nil, errors.New("not implemented")
}

func (s *minioStore) Find(ctx context.Context, opt *core.TargetStorageFindOpt) (*core.File, error) {
	log.Printf("[minio] find %q", opt.Query)
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

func (s *minioStore) Info(ctx context.Context) core.TargetInfo {
	return core.TargetInfo{}
}

var bucketNameRegexp = regexp.MustCompile("^[a-zA-Z0-9_.]+$")

func (s *minioStore) ValidateSettings(settings core.TargetSettings) (bool, error) {
	logger := log.WithField("target", TargetName)
	logger.Infof("ValidateSettings(): %+v", settings)

	for name := range DefaultConfig().Settings {
		value, ok := settings[name]
		if !ok {
			return false, fmt.Errorf("%q key is missing", name)
		}
		if value == "" {
			return false, fmt.Errorf("%q value is mepty", name)
		}
	}

	_, err := strconv.ParseBool(settings["MINIO_USE_SSL"])
	if err != nil {
		return false, fmt.Errorf("can't convert the value of MINIO_USE_SSL=%q to bool", settings["MINIO_USE_SSL"])
	}

	bucket := settings["MINIO_BUCKET_NAME"]
	match := bucketNameRegexp.MatchString(bucket)
	if !match {
		return false, fmt.Errorf("bucket name contains illegal characters")
	}

	return true, nil
}
