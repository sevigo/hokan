package directories

import (
	"encoding/base64"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sevigo/hokan/mocks"
	"github.com/sevigo/hokan/pkg/core"
	"github.com/stretchr/testify/assert"
)

const testPath = "C:\\Documents\\Fotos"

func TestFindByPath(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	dirStore := mocks.NewMockDirectoryStore(controller)
	dirStore.EXPECT().FindName(gomock.Any(), testPath).Return(&core.Directory{
		Path:      testPath,
		Recursive: true,
		Machine:   "test",
	}, nil)

	pathEnc := base64.StdEncoding.EncodeToString([]byte(testPath))
	url := fmt.Sprintf("/?path=%s", pathEnc)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", url, nil)

	HandleFind(dirStore)(w, r)

	assert.Equal(t, 200, w.Code)

	body := strings.TrimSpace(w.Body.String())
	assert.Equal(t, `{"Path":"C:\\Documents\\Fotos","Recursive":true,"Machine":"test","IgnoreFiles":null,"Target":null}`, body)
}

func TestErrorPath(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	dirStore := mocks.NewMockDirectoryStore(controller)
	url := fmt.Sprintf("/?path=%s", testPath)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", url, nil)

	HandleFind(dirStore)(w, r)

	assert.Equal(t, 400, w.Code)

	body := strings.TrimSpace(w.Body.String())
	assert.Equal(t, `{"code":400,"message":"invalid path"}`, body)
}
