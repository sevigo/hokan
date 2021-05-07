package directories_test

import (
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nicksnyder/basen"
	"github.com/stretchr/testify/assert"

	"github.com/sevigo/hokan/mocks"
	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/handler/api"
	"github.com/sevigo/hokan/pkg/testing/tools"
)

func TestFindByPath(t *testing.T) {
	testPath, err := ioutil.TempDir(os.TempDir(), "")
	assert.NoError(t, err)
	defer os.RemoveAll(testPath)

	controller := gomock.NewController(t)
	defer controller.Finish()

	dirStore := mocks.NewMockDirectoryStore(controller)
	pathID := basen.Base62Encoding.EncodeToString([]byte(testPath))
	dirStore.EXPECT().FindName(gomock.Any(), pathID).Return(&core.Directory{
		ID:        pathID,
		Path:      testPath,
		Recursive: true,
		Machine:   "test",
	}, nil)

	url := fmt.Sprintf("/directories/%s", pathID)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", url, nil)

	s := api.Server{
		Dirs: dirStore,
	}
	s.Handler().ServeHTTP(w, r)
	assert.Equal(t, 200, w.Code)

	body := strings.TrimSpace(w.Body.String())

	tools.TestJSONPathNotEmpty(t, "directory.id", body)
	tools.TestJSONPath(t, testPath, "directory.path", body)
	tools.TestJSONPath(t, "0", "stats.total-files", body)
	tools.TestJSONPath(t, "1", "stats.total-dirs", body)
}

func TestFindByPathNotFound(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	dirStore := mocks.NewMockDirectoryStore(controller)
	pathID := basen.Base62Encoding.EncodeToString([]byte("/test/path"))
	dirStore.EXPECT().FindName(gomock.Any(), pathID).Return(nil, core.ErrDirectoryNotFound)

	url := fmt.Sprintf("/directories/%s", pathID)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", url, nil)

	s := api.Server{
		Dirs: dirStore,
	}
	s.Handler().ServeHTTP(w, r)
	assert.Equal(t, 404, w.Code)
}
