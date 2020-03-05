package minio

import (
	"context"
	"path"

	"github.com/minio/minio-go"
	"github.com/rs/zerolog/log"

	"github.com/sevigo/hokan/pkg/core"
)

const TargetName = "minio"

const endpointDefault = "192.168.0.141:9000"
const accessKeyIDDefault = "minio"
const secretAccessKeyDefault = "miniostorage"
const useSSL = false

type minioStore struct {
	// endpoint        string //:= "play.min.io"
	// accessKeyID     string //:= "Q3AM3UQ867SPQQA43P2F"
	// secretAccessKey string //:= "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG"

	client     *minio.Client
	fs         core.FileStore
	bucketName string
}

func New(fs core.FileStore) (core.TargetStorage, error) {

	// Initialize minio client object.
	minioClient, err := minio.New(endpointDefault, accessKeyIDDefault, secretAccessKeyDefault, useSSL)
	if err != nil {
		log.Fatal().Err(err).Msg("can't create new minio client")
		return nil, err
	}

	// Make a new bucket with mchine name (must be from the config)
	bucketName := "osaka"
	location := "us-east-1"

	err = minioClient.MakeBucket(bucketName, location)
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Fatal().Err(err).Msg("can't create new minio client")
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

func (s *minioStore) Save(ctx context.Context, file *core.File) error {
	f, err := s.fs.Find(ctx, TargetName, file.Path)
	if err != nil { // TODO: not fount error
		log.Printf(">>> [%s] save a new file [%v]", TargetName, file)

		// Upload the file
		objectName := path.Clean(file.Path)

		// Upload the file with FPutObject
		options := minio.PutObjectOptions{
			UserMetadata: map[string]string{
				"path":     file.Path,
				"info":     file.Info,
				"checksum": file.Checksum,
			},
			// Progress                io.Reader
			// ContentType             string
			// ContentEncoding         string
			// ContentDisposition      string
			// ContentLanguage         string
			// CacheControl            string
			// ServerSideEncryption    encrypt.ServerSide
			// NumThreads              uint
			// StorageClass            string
			// WebsiteRedirectLocation string
		}
		n, err := s.client.FPutObject(s.bucketName, objectName, file.Path, options)
		if err != nil {
			return err
		}
		log.Printf("Successfully uploaded %s of size %d\n", objectName, n)
		return s.fs.Save(ctx, TargetName, file)
	}
	if f != nil && f.Checksum == file.Checksum {
		log.Printf("!!! [%s] ignore saving, file [%v] already stored", TargetName, file)
		return nil
	}
	if f != nil && f.Checksum != file.Checksum {
		log.Printf("!!! [%s] save changed file [%v]", TargetName, file)
		return s.fs.Save(ctx, TargetName, file)
	}

	return nil
}

func (s *minioStore) List(context.Context) ([]*core.File, error) {
	log.Printf("[minio] list\n")
	return nil, nil
}

func (s *minioStore) Find(ctx context.Context, q string) (*core.File, error) {
	log.Printf("[minio] find %q\n", q)
	return nil, nil
}

func (s *minioStore) Delete(ctx context.Context, file *core.File) error {
	log.Printf("[minio] save %#v\n", file)

	return nil
}
