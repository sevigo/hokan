package minio

import (
	"context"
	"errors"
	"path"

	"github.com/minio/minio-go"
	"github.com/sevigo/hokan/pkg/core"
	log "github.com/sirupsen/logrus"

	filestore "github.com/sevigo/hokan/pkg/store/file"
)

const TargetName = "minio"

const endpointDefault = "192.168.0.141:9000"
const accessKeyIDDefault = "minio"
const secretAccessKeyDefault = "miniostorage"
const useSSL = false

type minioStore struct {
	client     *minio.Client
	fs         core.FileStore
	bucketName string
}

func New(ctx context.Context, fs core.FileStore) (core.TargetStorage, error) {
	// Initialize minio client object.
	minioClient, err := minio.New(endpointDefault, accessKeyIDDefault, secretAccessKeyDefault, useSSL)
	if err != nil {
		log.WithError(err).Error("Can't create new minio client")
		return nil, err
	}

	// Make a new bucket with mchine name (must be from the config)
	bucketName := "osaka"
	err = minioClient.MakeBucket(bucketName, "")
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			return nil, err
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}

	return &minioStore{
		bucketName: bucketName,
		client:     minioClient,
		fs:         fs,
	}, nil
}

func fileHasChanged(newFile, storedFile *core.File) bool {
	if storedFile == nil {
		return true
	}
	if newFile.Checksum != storedFile.Checksum {
		return true
	}
	return false
}

func (s *minioStore) Save(ctx context.Context, file *core.File) error {
	storedFile, err := s.fs.Find(ctx, TargetName, file.Path)

	// logger := log.Debug().Str("target", TargetName).Str("file", file.Path).Str("op", "save")
	logger := log.WithFields(log.Fields{
		"target": TargetName,
		"file":   file.Path,
	})

	if errors.Is(err, filestore.ErrFileEntryNotFound) || fileHasChanged(file, storedFile) {
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
