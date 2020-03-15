package minio

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sevigo/hokan/mocks"
	"github.com/sevigo/hokan/pkg/core"
	filestore "github.com/sevigo/hokan/pkg/store/file"
	"github.com/stretchr/testify/assert"
)

var testFilePath = "/test/file.txt"
var testBucket = "test"

func Test_minioStore_SaveNewFile(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	file := &core.File{
		Path:     testFilePath,
		Checksum: "abc",
		Targets:  []string{"minio"},
	}

	fileStore := mocks.NewMockFileStore(controller)
	fileStore.EXPECT().Find(context.TODO(), TargetName, testFilePath).Return(nil, filestore.ErrFileEntryNotFound)
	fileStore.EXPECT().Save(context.TODO(), TargetName, file).Return(nil)

	minioClient := mocks.NewMockMinioWrapper(controller)
	minioClient.EXPECT().FPutObjectWithContext(context.TODO(), testBucket, testFilePath, testFilePath, gomock.Any()).Return(int64(64), nil)

	store := &minioStore{
		bucketName: testBucket,
		fs:         fileStore,
		client:     minioClient,
	}

	err := store.Save(context.TODO(), file)
	assert.NoError(t, err)
}

func Test_minioStore_SaveFileChange(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	fileA := &core.File{
		Path:     testFilePath,
		Checksum: "abc",
		Targets:  []string{"minio"},
	}
	fileB := &core.File{
		Path:     testFilePath,
		Checksum: "abX",
		Targets:  []string{"minio"},
	}

	fileStore := mocks.NewMockFileStore(controller)
	fileStore.EXPECT().Find(context.TODO(), TargetName, testFilePath).Return(fileA, nil)
	fileStore.EXPECT().Save(context.TODO(), TargetName, fileB).Return(nil)

	minioClient := mocks.NewMockMinioWrapper(controller)
	minioClient.EXPECT().FPutObjectWithContext(context.TODO(), testBucket, testFilePath, testFilePath, gomock.Any()).Return(int64(64), nil)

	store := &minioStore{
		bucketName: testBucket,
		fs:         fileStore,
		client:     minioClient,
	}

	err := store.Save(context.TODO(), fileB)
	assert.NoError(t, err)
}

func Test_minioStore_NoSave(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	fileA := &core.File{
		Path:     testFilePath,
		Checksum: "abc",
		Targets:  []string{"minio"},
	}

	fileStore := mocks.NewMockFileStore(controller)
	fileStore.EXPECT().Find(context.TODO(), TargetName, testFilePath).Return(fileA, nil)
	minioClient := mocks.NewMockMinioWrapper(controller)

	store := &minioStore{
		bucketName: testBucket,
		fs:         fileStore,
		client:     minioClient,
	}

	err := store.Save(context.TODO(), fileA)
	assert.NoError(t, err)
}

func Test_minioStore_ErrorNoSave(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	file := &core.File{
		Path:     testFilePath,
		Checksum: "abc",
		Targets:  []string{"minio"},
	}

	fileStore := mocks.NewMockFileStore(controller)
	fileStore.EXPECT().Find(context.TODO(), TargetName, testFilePath).Return(nil, filestore.ErrFileEntryNotFound)
	minioClient := mocks.NewMockMinioWrapper(controller)
	minioClient.EXPECT().FPutObjectWithContext(context.TODO(), testBucket, testFilePath, testFilePath, gomock.Any()).Return(int64(0), fmt.Errorf("error"))

	store := &minioStore{
		bucketName: testBucket,
		fs:         fileStore,
		client:     minioClient,
	}

	err := store.Save(context.TODO(), file)
	assert.Error(t, err)
}

func TestPing(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	fileStore := mocks.NewMockFileStore(controller)
	minioClient := mocks.NewMockMinioWrapper(controller)
	minioClient.EXPECT().BucketExists(testBucket).Return(true, nil)
	store := &minioStore{
		bucketName: testBucket,
		fs:         fileStore,
		client:     minioClient,
	}

	err := store.Ping(context.TODO())
	assert.NoError(t, err)
}
