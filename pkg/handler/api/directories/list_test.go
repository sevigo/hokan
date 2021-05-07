package directories_test

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sevigo/hokan/mocks"
	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/handler/api"
	"github.com/sevigo/hokan/pkg/testing/tools"
	"github.com/stretchr/testify/assert"
)

func TestHandleListEmpty(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	dirStore := mocks.NewMockDirectoryStore(controller)
	dirStore.EXPECT().List(gomock.Any()).Return([]*core.Directory{}, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/directories", nil)

	s := api.Server{
		Dirs: dirStore,
	}
	s.Handler().ServeHTTP(w, r)
	assert.Equal(t, 200, w.Code)

	body := strings.TrimSpace(w.Body.String())
	assert.NotEmpty(t, body)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

	tools.TestJSONPath(t, "self", "links.0.rel", body)
	tools.TestJSONPath(t, "0", "meta.total_items", body)
}

func TestHandleLisOK(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	dirStore := mocks.NewMockDirectoryStore(controller)
	dirStore.EXPECT().List(gomock.Any()).Return([]*core.Directory{
		{
			ID:        "abc",
			Path:      "/test/dir",
			Recursive: true,
			Machine:   "test",
		},
	}, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/directories", nil)

	s := api.Server{
		Dirs: dirStore,
	}
	s.Handler().ServeHTTP(w, r)
	assert.Equal(t, 200, w.Code)

	body := strings.TrimSpace(w.Body.String())
	assert.NotEmpty(t, body)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

	tools.TestJSONPath(t, "self", "links.0.rel", body)
	tools.TestJSONPath(t, "1", "meta.total_items", body)

	tools.TestJSONPath(t, "abc", "directories.0.id", body)
	tools.TestJSONPath(t, "/test/dir", "directories.0.path", body)
	tools.TestJSONPath(t, "true", "directories.0.recursive", body)
}
