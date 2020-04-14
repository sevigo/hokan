package directories_test

import (
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nicksnyder/basen"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/sevigo/hokan/mocks"
	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/handler/api"
	"github.com/sevigo/hokan/pkg/testing/tools"
)

const testPath = "C:\\Documents\\Fotos"

func TestFindByPath(t *testing.T) {
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
		Dirs:   dirStore,
		Logger: logrus.StandardLogger(),
	}
	s.Handler().ServeHTTP(w, r)
	assert.Equal(t, 200, w.Code)

	body := strings.TrimSpace(w.Body.String())
	tools.TestJSONPath(t, "YsmKL73TlYdFBq4g6vBYaZKl", "id", body)
	tools.TestJSONPath(t, "false", "active", body)
	tools.TestJSONPath(t, "C:\\Documents\\Fotos", "path", body)
}

func TestFindByPathNotFound(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	dirStore := mocks.NewMockDirectoryStore(controller)
	pathID := basen.Base62Encoding.EncodeToString([]byte(testPath))
	dirStore.EXPECT().FindName(gomock.Any(), pathID).Return(nil, core.ErrDirectoryNotFound)

	url := fmt.Sprintf("/directories/%s", pathID)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", url, nil)

	s := api.Server{
		Dirs:   dirStore,
		Logger: logrus.StandardLogger(),
	}
	s.Handler().ServeHTTP(w, r)
	assert.Equal(t, 404, w.Code)
}
