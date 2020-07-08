package config

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
	"github.com/sevigo/hokan/mocks"
	"github.com/sevigo/hokan/pkg/core"
)

func Test_configStore_Save(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	db := mocks.NewMockDB(controller)
	db.EXPECT().Write(gomock.Any(), gomock.Any(), gomock.Any()).
		Do(func(bucket, key, value string) {
			assert.Equal(t, "config", bucket)
			assert.Equal(t, "target:test", key)
			assert.Equal(t, `{"active":true,"name":"test","description":"test confing","settings":{"foo":"bar"}}`, strings.TrimSpace(value))
		}).Return(nil)

	s := configStore{db}
	err := s.Save(context.TODO(), &core.TargetConfig{
		Name:        "test",
		Active:      true,
		Description: "test confing",
		Settings: core.TargetSettings{
			"foo": "bar",
		},
	})
	assert.NoError(t, err)
}

func Test_configStore_Find(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	db := mocks.NewMockDB(controller)
	db.EXPECT().Read("config", "target:test").
		Return([]byte(`{"active":true,"name":"test","description":"test confing","settings":{"foo":"bar"}}`), nil)

	s := configStore{db}
	config, err := s.Find(context.TODO(), "test")
	assert.Equal(t, &core.TargetConfig{
		Name:        "test",
		Active:      true,
		Description: "test confing",
		Settings: core.TargetSettings{
			"foo": "bar",
		},
	}, config)
	assert.NoError(t, err)
}
