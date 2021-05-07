package directory

import (
	"context"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sevigo/hokan/mocks"
	"github.com/sevigo/hokan/pkg/core"
	"github.com/stretchr/testify/assert"
)

func Test_directoryStore_List(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	db := mocks.NewMockDB(controller)
	db.EXPECT().ReadBucket("watch:directories", &core.ReadBucketOptions{}).Return([]core.BucketData{
		{
			Key:   "foo",
			Value: `{"active":true,"path":"/foo"}`,
		},
	}, nil)

	s := directoryStore{db}
	dirs, err := s.List(context.TODO())
	assert.Equal(t, "/foo", dirs[0].Path)
	assert.NoError(t, err)
}

func Test_directoryStore_FindName(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	db := mocks.NewMockDB(controller)
	db.EXPECT().Read("watch:directories", "/foo").Return([]byte(`{"active":true,"path":"/foo"}`), nil)

	s := directoryStore{db}
	dir, err := s.FindName(context.TODO(), "/foo")
	assert.Equal(t, "/foo", dir.Path)
	assert.NoError(t, err)
}

func Test_directoryStore_Create(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	db := mocks.NewMockDB(controller)
	db.EXPECT().Write(gomock.Any(), gomock.Any(), gomock.Any()).
		Do(func(bucket, key, value string) {
			assert.Equal(t, "watch:directories", bucket)
			assert.Equal(t, "rokrX", key)
			assert.Equal(t, `{"id":"rokrX","path":"/foo","recursive":false,"machine":"","ignore":null}`, strings.TrimSpace(value))
		}).Return(nil)

	s := directoryStore{db}
	err := s.Create(context.TODO(), &core.Directory{Path: "/foo"})
	assert.NoError(t, err)
}
