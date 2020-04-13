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
	"github.com/sevigo/hokan/pkg/watcher/utils"
)

var testFilePath = "testdata/test.txt"
var testBucket = "test"
var expectedValue = `{"ID":"Ecsz84HlZ1rfn2L4DBUNZn2","Path":"testdata/test.txt","Checksum":"5e2bf57d3f40c4b6df69daf1936cb766f832374b4fc0259a7cbff06e2f70f269","Info":{"ModTime":"2020-04-13T23:36:02.6754743+02:00","Mode":438,"Name":"test.txt","Size":11},"Targets":["test"]}
`

func getTestingFile(t *testing.T) string {
	pwd, err := os.Getwd()
	assert.NoError(t, err)
	return filepath.Join(pwd, testFilePath)
}

func Test_fileStore_Save(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	db := mocks.NewMockDB(controller)
	db.EXPECT().Write(testBucket, "testdata/test.txt", expectedValue).Return(nil)

	checksum, info, err := utils.FileChecksumInfo(getTestingFile(t))
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
