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

func TestConfig(t *testing.T) {
	store := &minioStore{}
	conf := store.DefaultConfig()
	assert.Equal(t, "minio", conf.Name)
	assert.Equal(t, false, conf.Active)
}

func TestNewNotActive(t *testing.T) {
	store := &minioStore{}
	conf := store.DefaultConfig()
	_, err := New(context.Background(), nil, *conf)
	assert.EqualError(t, err, "target is not active")
}

func TestNewActiveErr(t *testing.T) {
	store := &minioStore{}
	conf := store.DefaultConfig()
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
	fileStore.EXPECT().Save(context.TODO(), TargetName, file).Return(nil)

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

	result := make(chan core.TargetOperationResult)
	go store.Save(context.TODO(), result, file, nil)
	data := <-result

	assert.Equal(t, core.TargetOperationResult{
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
		Targets:  []string{"minio"},
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

	result := make(chan core.TargetOperationResult)
	go store.Save(context.TODO(), result, file, nil)
	data := <-result

	assert.Equal(t, core.TargetOperationResult{
		Success: false,
		Error:   errSave,
		Message: "error on save",
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
		Targets:  []string{"minio"},
		Info:     info,
	}

	fileStore := mocks.NewMockFileStore(controller)
	fileStore.EXPECT().Save(context.TODO(), TargetName, fileB).Return(nil)

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

	result := make(chan core.TargetOperationResult)
	go store.Save(context.TODO(), result, fileB, nil)
	data := <-result

	assert.Equal(t, core.TargetOperationResult{
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

func Test_localStorage_ValidateSettings(t *testing.T) {
	store := &minioStore{}

	tests := []struct {
		name     string
		settings core.TargetSettings
		wantErr  bool
	}{
		{
			name: "case 1",
			settings: core.TargetSettings{
				"MINIO_HOST":        "http://localhost:8081",
				"MINIO_ACCESS_KEY":  "abc",
				"MINIO_SECRET_KEY":  "xyz",
				"MINIO_USE_SSL":     "false",
				"MINIO_BUCKET_NAME": "test",
			},
			wantErr: false,
		},
		{
			name: "case 2",
			settings: core.TargetSettings{
				"MINIO_ACCESS_KEY":  "",
				"MINIO_SECRET_KEY":  "",
				"MINIO_USE_SSL":     "",
				"MINIO_BUCKET_NAME": "",
			},
			wantErr: true,
		},
		{
			name: "case 3",
			settings: core.TargetSettings{
				"MINIO_HOST":        "",
				"MINIO_ACCESS_KEY":  "",
				"MINIO_SECRET_KEY":  "",
				"MINIO_USE_SSL":     "",
				"MINIO_BUCKET_NAME": "",
			},
			wantErr: true,
		},
		{
			name: "case 4",
			settings: core.TargetSettings{
				"MINIO_HOST":        "http://localhost:8081",
				"MINIO_ACCESS_KEY":  "abc",
				"MINIO_SECRET_KEY":  "xyz",
				"MINIO_USE_SSL":     "no",
				"MINIO_BUCKET_NAME": "test",
			},
			wantErr: true,
		},
		{
			name: "case 5",
			settings: core.TargetSettings{
				"MINIO_HOST":        "http://localhost:8081",
				"MINIO_ACCESS_KEY":  "abc",
				"MINIO_SECRET_KEY":  "xyz",
				"MINIO_USE_SSL":     "true",
				"MINIO_BUCKET_NAME": "!<.test",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := store.ValidateSettings(tt.settings)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_minioStore_Info(t *testing.T) {
	store := &minioStore{}
	info := store.Info(context.TODO())
	assert.NotNil(t, info)
}

func Test_minioStore_Delete(t *testing.T) {
	store := &minioStore{}
	result := make(chan core.TargetOperationResult)
	go store.Delete(context.TODO(), result, []*core.File{}, &core.TargetStorageDeleteOpt{})
	data := <-result
	assert.Equal(t, core.TargetOperationResult{
		Success: false,
		Error:   core.ErrNotImplemented,
		Message: "not implemented",
	}, data)
}

func Test_minioStore_Restore(t *testing.T) {
	store := &minioStore{}
	result := make(chan core.TargetOperationResult)
	go store.Restore(context.TODO(), result, []*core.File{}, &core.TargetStorageRestoreOpt{})
	data := <-result
	assert.Equal(t, core.TargetOperationResult{
		Success: false,
		Error:   core.ErrNotImplemented,
		Message: "not implemented",
	}, data)
}
