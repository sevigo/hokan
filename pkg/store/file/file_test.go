package file

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/sevigo/hokan/mocks"
	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/store/db"
	"github.com/sevigo/hokan/pkg/testing/tools"
	"github.com/sevigo/hokan/pkg/watcher/utils"
)

var testFilePath = "testdata/test.txt"
var testBucket = "test"

func getTestingFile(t *testing.T) string {
	pwd, err := os.Getwd()
	assert.NoError(t, err)
	return filepath.Join(pwd, testFilePath)
}

func Test_fileStore_Save(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	db := mocks.NewMockDB(controller)
	testFilePath := getTestingFile(t)

	db.EXPECT().Write(gomock.Any(), gomock.Any(), gomock.Any()).Do(func(bucketName, key, value string) error {
		assert.Equal(t, testBucket, bucketName)
		tools.TestJSONPath(t, "5e2bf57d3f40c4b6df69daf1936cb766f832374b4fc0259a7cbff06e2f70f269", "checksum", value)
		tools.TestJSONPath(t, "test.txt", "info.name", value)
		tools.TestJSONPath(t, "11", "info.size", value)
		return nil
	})

	checksum, info, err := utils.FileChecksumInfo(testFilePath)
	assert.NoError(t, err)

	file := &core.File{
		Path:     testFilePath,
		Checksum: checksum,
		Info:     info,
	}

	store := New(db)
	err = store.Save(context.TODO(), testBucket, file)
	assert.NoError(t, err)
}

func TestListOffsetLimit(t *testing.T) {
	tmpDir, err := ioutil.TempDir(os.TempDir(), "")
	assert.NoError(t, err)
	dbFile := path.Join(tmpDir, "test.db")
	defer os.RemoveAll(tmpDir)

	storage, err := db.Connect(dbFile)
	assert.NoError(t, err)
	assert.NotNil(t, storage)

	_, info, err := utils.FileChecksumInfo(dbFile)
	assert.NoError(t, err)

	fileStore := New(storage)
	for i := 0; i < 1003; i++ {
		name := fmt.Sprintf("/test/foo/file_%04d.png", i)
		file := &core.File{
			Path:     name,
			Checksum: "abc",
			Info:     info,
		}
		err := fileStore.Save(context.TODO(), testBucket, file)
		assert.NoError(t, err)
	}

	tests := []struct {
		name         string
		opt          *core.FileListOptions
		expectedData []string
	}{
		{
			name: "case 1",
			opt: &core.FileListOptions{
				Offset: 0,
				Limit:  5,
			},
			expectedData: []string{
				"/test/foo/file_0000.png",
				"/test/foo/file_0001.png",
				"/test/foo/file_0002.png",
				"/test/foo/file_0003.png",
				"/test/foo/file_0004.png",
			},
		},
		{
			name: "case 2",
			opt: &core.FileListOptions{
				Offset: 5,
				Limit:  5,
			},
			expectedData: []string{
				"/test/foo/file_0005.png",
				"/test/foo/file_0006.png",
				"/test/foo/file_0007.png",
				"/test/foo/file_0008.png",
				"/test/foo/file_0009.png",
			},
		},
		{
			name: "case 3",
			opt: &core.FileListOptions{
				Offset: 100,
				Limit:  3,
			},
			expectedData: []string{
				"/test/foo/file_0100.png",
				"/test/foo/file_0101.png",
				"/test/foo/file_0102.png",
			},
		},
		{
			name: "case 4",
			opt: &core.FileListOptions{
				Offset: 1000,
				Limit:  5,
			},
			expectedData: []string{
				"/test/foo/file_1000.png",
				"/test/foo/file_1001.png",
				"/test/foo/file_1002.png",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := fileStore.List(context.TODO(), testBucket, tt.opt)
			assert.NoError(t, err)
			files := []string{}
			for _, f := range data {
				files = append(files, f.Path)
			}
			assert.Equal(t, tt.expectedData, files)
		})
	}
}
