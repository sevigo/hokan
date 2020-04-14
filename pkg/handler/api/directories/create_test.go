package directories

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nicksnyder/basen"
	"github.com/stretchr/testify/assert"

	"github.com/sevigo/hokan/mocks"
	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/testing/tools"
)

var testFotosPath string = "C:\\Documents\\Fotos"

func TestCreate(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	dirStore := mocks.NewMockDirectoryStore(controller)
	dirStore.EXPECT().Create(gomock.Any(), gomock.Any()).Do(func(_ context.Context, dir *core.Directory) error {
		assert.Equal(t, testFotosPath, dir.Path)
		return nil
	})
	eventCreator := mocks.NewMockEventCreator(controller)
	eventCreator.EXPECT().Publish(gomock.Any(), &core.EventData{
		Type: core.WatchDirStart,
		Data: &core.Directory{
			ID:        basen.Base62Encoding.EncodeToString([]byte(testFotosPath)),
			Path:      testFotosPath,
			Active:    true,
			Machine:   "test",
			Recursive: true,
		},
	})

	in := new(bytes.Buffer)
	err := json.NewEncoder(in).Encode(&core.Directory{Path: testFotosPath, Active: true, Machine: "test", Recursive: true})
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/", in)

	HandleCreate(dirStore, eventCreator)(w, r)
	body := strings.TrimSpace(w.Body.String())

	assert.Equal(t, 201, w.Code)
	tools.TestJSONPath(t, "YsmKL73TlYdFBq4g6vBYaZKl", "id", body)
	tools.TestJSONPath(t, "true", "active", body)
	tools.TestJSONPath(t, "C:\\Documents\\Fotos", "path", body)
	tools.TestJSONPath(t, "test", "machine", body)
	tools.TestJSONPath(t, "true", "recursive", body)
}

func TestCreateBadRequest(t *testing.T) {
	in := new(bytes.Buffer)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/", in)

	HandleCreate(nil, nil)(w, r)
	body := strings.TrimSpace(w.Body.String())

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, `{"code":400,"message":"invalid request body"}`, body)
}

func TestCreateError(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	dirStore := mocks.NewMockDirectoryStore(controller)
	dirStore.EXPECT().Create(gomock.Any(), gomock.Any()).Return(fmt.Errorf("error"))
	eventCreator := mocks.NewMockEventCreator(controller)

	in := new(bytes.Buffer)
	err := json.NewEncoder(in).Encode(&core.Directory{Path: "C:\\Documents\\Fotos"})
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/", in)

	HandleCreate(dirStore, eventCreator)(w, r)
	body := strings.TrimSpace(w.Body.String())

	assert.Equal(t, 500, w.Code)
	assert.Equal(t, `{"code":500,"message":"cannot store a new directory"}`, body)
}
