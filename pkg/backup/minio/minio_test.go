package minio

import (
	"context"
	"errors"
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

func getTestingFile(t *testing.T) string {
	pwd, err := os.Getwd()
	assert.NoError(t, err)
	return filepath.Join(pwd, testFilePath)
}

func Test_minioStore_SaveNewFile(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	checksum, info, err := utils.FileChecksumInfo(getTestingFile(t))
	assert.NoError(t, err)

	file := &core.File{
		Path:     testFilePath,
		Checksum: checksum,
		Info:     info,
	}

	fileStore := mocks.NewMockFileStore(controller)
	fileStore.EXPECT().Save(context.TODO(), name, file).Return(nil)

	minioClient := mocks.NewMockMinioWrapper(controller)
	minioClient.EXPECT().FPutObjectWithContext(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Do(func(_ context.Context, bucketName, objectName, filePath string, opts minio.PutObjectOptions) (int64, error) {
			assert.Equal(t, testBucket, bucketName)
			assert.Equal(t, "5e2bf57d3f40c4b6df69daf1936cb766f832374b4fc0259a7cbff06e2f70f269", opts.UserMetadata["checksum"])

			return int64(11), nil
		})

	store := &minioStore{
		bucketName: testBucket,
		fileStore:  fileStore,
		client:     minioClient,
	}

	result := make(chan core.BackupResult)
	go store.Save(context.TODO(), result, file, nil)
	data := <-result

	assert.Equal(t, core.BackupResult{
		Success: true,
		Error:   nil,
		Message: "requested operation was successful",
	}, data)
}

func Test_minioStore_SaveError(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	checksum, info, err := utils.FileChecksumInfo(getTestingFile(t))
	assert.NoError(t, err)

	file := &core.File{
		Path:     testFilePath,
		Checksum: checksum,
		Info:     info,
	}

	fileStore := mocks.NewMockFileStore(controller)
	minioClient := mocks.NewMockMinioWrapper(controller)
	errSave := errors.New("error on save")
	minioClient.EXPECT().FPutObjectWithContext(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(int64(0), errSave)

	store := &minioStore{
		bucketName: testBucket,
		fileStore:  fileStore,
		client:     minioClient,
	}

	result := make(chan core.BackupResult)
	go store.Save(context.TODO(), result, file, nil)
	data := <-result

	assert.Equal(t, core.BackupResult{
		Success: false,
		Error:   errSave,
		Message: "can't save file \"testdata/test.txt\" to \"test\"",
	}, data)
}

func Test_minioStore_SaveFileChange(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	_, info, err := utils.FileChecksumInfo(getTestingFile(t))
	assert.NoError(t, err)

	fileB := &core.File{
		Path:     testFilePath,
		Checksum: "abX",
		Info:     info,
	}

	fileStore := mocks.NewMockFileStore(controller)
	fileStore.EXPECT().Save(context.TODO(), name, fileB).Return(nil)

	minioClient := mocks.NewMockMinioWrapper(controller)
	minioClient.EXPECT().FPutObjectWithContext(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Do(func(_ context.Context, bucketName, objectName, filePath string, opts minio.PutObjectOptions) (int64, error) {
			assert.Equal(t, testBucket, bucketName)
			assert.Equal(t, "abX", opts.UserMetadata["checksum"])

			return int64(11), nil
		})

	store := &minioStore{
		bucketName: testBucket,
		fileStore:  fileStore,
		client:     minioClient,
	}

	result := make(chan core.BackupResult)
	go store.Save(context.TODO(), result, fileB, nil)
	data := <-result

	assert.Equal(t, core.BackupResult{
		Success: true,
		Error:   nil,
		Message: "requested operation was successful",
	}, data)
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

func Test_minioStore_Delete(t *testing.T) {
	store := &minioStore{}
	result := make(chan core.BackupResult)
	go store.Delete(context.TODO(), result, []*core.File{}, &core.BackupOperationOptions{})
	data := <-result
	assert.Equal(t, core.BackupResult{
		Success: false,
		Error:   core.ErrNotImplemented,
	}, data)
}

func Test_minioStore_Restore(t *testing.T) {
	store := &minioStore{}
	result := make(chan core.BackupResult)
	go store.Restore(context.TODO(), result, []*core.File{}, &core.BackupOperationOptions{})
	data := <-result
	assert.Equal(t, core.BackupResult{
		Success: false,
		Error:   core.ErrNotImplemented,
	}, data)
}
