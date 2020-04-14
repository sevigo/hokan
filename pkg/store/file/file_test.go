package file

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/sevigo/hokan/mocks"
	"github.com/sevigo/hokan/pkg/core"
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
		tools.TestJSONPath(t, "test", "targets.0", value)
		return nil
	})

	checksum, info, err := utils.FileChecksumInfo(testFilePath)
	assert.NoError(t, err)

	file := &core.File{
		Path:     testFilePath,
		Checksum: checksum,
		Targets:  []string{"test"},
		Info:     info,
	}

	store := New(db)
	err = store.Save(context.TODO(), testBucket, file)
	assert.NoError(t, err)
}
