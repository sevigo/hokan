package directories

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
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
	assert.Equal(t, w.Code, 201)
}
