package files_test

import (
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/sevigo/hokan/mocks"
	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/handler/api"
	"github.com/sevigo/hokan/pkg/testing/tools"
)

var targetName = "test"

func TestHandleList(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	fileStore := mocks.NewMockFileStore(controller)
	fileStore.EXPECT().List(gomock.Any(), &core.FileListOptions{
		TargetName: targetName,
	}).Return([]*core.File{
		{
			ID:       "1",
			Path:     "/foo.txt",
			Checksum: "abc",
		},
		{
			ID:       "2",
			Path:     "/bar.txt",
			Checksum: "xzy",
		},
	}, nil)

	url := fmt.Sprintf("/targets/%s/files", targetName)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", url, nil)

	s := api.Server{
		Files:  fileStore,
		Logger: logrus.StandardLogger(),
	}
	s.Handler().ServeHTTP(w, r)
	assert.Equal(t, 200, w.Code)

	jsonStr := strings.TrimSpace(w.Body.String())
	assert.NotEmpty(t, jsonStr)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

	tools.TestJSONPath(t, "/foo.txt", "files.0.path", jsonStr)
	tools.TestJSONPath(t, "/bar.txt", "files.1.path", jsonStr)
	tools.TestJSONPath(t, "/bar.txt", "files.1.path", jsonStr)
	tools.TestJSONPath(t, "self", "links.0.rel", jsonStr)
	tools.TestJSONPath(t, "/targets/test/files", "links.0.href", jsonStr)
	tools.TestJSONPath(t, "GET", "links.0.method", jsonStr)
	tools.TestJSONPath(t, "0", "meta.total_items", jsonStr)
}
