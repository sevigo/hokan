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
	"github.com/stretchr/testify/assert"

	"github.com/sevigo/hokan/mocks"
	"github.com/sevigo/hokan/pkg/core"
)

func TestCreate(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	dirStore := mocks.NewMockDirectoryStore(controller)
	dirStore.EXPECT().Create(gomock.Any(), gomock.Any()).Do(func(_ context.Context, dir *core.Directory) error {
		assert.Equal(t, "C:\\Documents\\Fotos", dir.Path)
		return nil
	})
	eventCreator := mocks.NewMockEventCreator(controller)

	in := new(bytes.Buffer)
	err := json.NewEncoder(in).Encode(&core.Directory{Path: "C:\\Documents\\Fotos"})
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/", in)

	HandleCreate(dirStore, eventCreator)(w, r)
	body := strings.TrimSpace(w.Body.String())

	assert.Equal(t, 201, w.Code)
	assert.Equal(t, `{"Path":"C:\\Documents\\Fotos","Recursive":false,"Machine":"","IgnoreFiles":null,"Target":null}`, body)
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
