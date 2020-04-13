package minio

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/minio/minio-go"
	"github.com/sevigo/hokan/mocks"
	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/watcher/utils"

	"github.com/stretchr/testify/assert"
)

var testFilePath = "testdata/test.txt"
var testBucket = "test"
var expectedOpts = minio.PutObjectOptions{
	UserMetadata: map[string]string{
		"path":      "testdata/test.txt",
		"size":      "11",
		"name":      "test.txt",
		"mode-time": "2020-04-13 22:07:56.013644 +0200 CEST",
		"mode":      "-rw-rw-rw-",
		"checksum":  "5e2bf57d3f40c4b6df69daf1936cb766f832374b4fc0259a7cbff06e2f70f269",
	},
}

func getTestingFile(t *testing.T) string {
	pwd, err := os.Getwd()
	assert.NoError(t, err)
	return filepath.Join(pwd, testFilePath)
}

func TestConfig(t *testing.T) {
	conf := DefaultConfig()
	assert.Equal(t, "minio", conf.Name)
	assert.Equal(t, false, conf.Active)
}

func TestNewNotActive(t *testing.T) {
	conf := DefaultConfig()
	_, err := New(context.Background(), nil, *conf)
	assert.EqualError(t, err, "target is not active")
}

func TestNewActiveErr(t *testing.T) {
	conf := DefaultConfig()
	conf.Active = true
	_, err := New(context.Background(), nil, *conf)
	assert.Error(t, err)
}

func Test_minioStore_SaveNewFile(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	checksum, info, err := utils.FileChecksumInfo(getTestingFile(t))
	assert.NoError(t, err)

	file := &core.File{
		Path:     testFilePath,
		Checksum: checksum,
		Targets:  []string{"minio"},
		Info:     info,
	}

	fileStore := mocks.NewMockFileStore(controller)
	fileStore.EXPECT().Find(context.TODO(), &core.FileSearchOptions{
		FilePath:   testFilePath,
		TargetName: TargetName,
	}).Return(nil, core.ErrFileNotFound)
	fileStore.EXPECT().Save(context.TODO(), TargetName, file).Return(nil)

	minioClient := mocks.NewMockMinioWrapper(controller)
	minioClient.EXPECT().FPutObjectWithContext(context.TODO(), testBucket, testFilePath, testFilePath, expectedOpts).Return(int64(64), nil)

	store := &minioStore{
		bucketName: testBucket,
		fileStore:  fileStore,
		client:     minioClient,
	}

	err = store.Save(context.TODO(), file)
	assert.NoError(t, err)
}

func Test_minioStore_SaveFileChange(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	checksum, info, err := utils.FileChecksumInfo(getTestingFile(t))
	assert.NoError(t, err)

	fileA := &core.File{
		Path:     testFilePath,
		Checksum: checksum,
		Targets:  []string{"minio"},
		Info:     info,
	}
	fileB := &core.File{
		Path:     testFilePath,
		Checksum: "abX",
		Targets:  []string{"minio"},
		Info:     info,
	}

	fileStore := mocks.NewMockFileStore(controller)
	fileStore.EXPECT().Find(context.TODO(), &core.FileSearchOptions{
		TargetName: TargetName,
		FilePath:   testFilePath,
	}).Return(fileA, nil)
	fileStore.EXPECT().Save(context.TODO(), TargetName, fileB).Return(nil)

	minioClient := mocks.NewMockMinioWrapper(controller)
	minioClient.EXPECT().FPutObjectWithContext(context.TODO(), testBucket, testFilePath, testFilePath, gomock.Any()).Return(int64(64), nil)

	store := &minioStore{
		bucketName: testBucket,
		fileStore:  fileStore,
		client:     minioClient,
	}

	err = store.Save(context.TODO(), fileB)
	assert.NoError(t, err)
}

func Test_minioStore_NoSave(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	checksum, info, err := utils.FileChecksumInfo("minio_test.go")
	assert.NoError(t, err)

	fileA := &core.File{
		Path:     testFilePath,
		Checksum: checksum,
		Targets:  []string{"minio"},
		Info:     info,
	}

	fileStore := mocks.NewMockFileStore(controller)
	fileStore.EXPECT().Find(context.TODO(), &core.FileSearchOptions{
		TargetName: TargetName,
		FilePath:   testFilePath,
	}).Return(fileA, nil)
	minioClient := mocks.NewMockMinioWrapper(controller)

	store := &minioStore{
		bucketName: testBucket,
		fileStore:  fileStore,
		client:     minioClient,
	}

	err = store.Save(context.TODO(), fileA)
	assert.NoError(t, err)
}

func Test_minioStore_ErrorNoSave(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	checksum, info, err := utils.FileChecksumInfo("minio_test.go")
	assert.NoError(t, err)

	file := &core.File{
		Path:     testFilePath,
		Checksum: checksum,
		Info:     info,
		Targets:  []string{"minio"},
	}

	fileStore := mocks.NewMockFileStore(controller)
	fileStore.EXPECT().Find(context.TODO(), &core.FileSearchOptions{
		TargetName: TargetName,
		FilePath:   testFilePath,
	}).Return(nil, core.ErrFileNotFound)

	minioClient := mocks.NewMockMinioWrapper(controller)
	minioClient.EXPECT().FPutObjectWithContext(context.TODO(), testBucket, testFilePath, testFilePath, gomock.Any()).Return(int64(0), fmt.Errorf("error"))

	store := &minioStore{
		bucketName: testBucket,
		fileStore:  fileStore,
		client:     minioClient,
	}

	err = store.Save(context.TODO(), file)
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
		fileStore:  fileStore,
		client:     minioClient,
	}

	err := store.Ping(context.TODO())
	assert.NoError(t, err)
}
