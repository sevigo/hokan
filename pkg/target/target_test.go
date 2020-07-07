package target

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sevigo/hokan/mocks"
	"github.com/sevigo/hokan/pkg/core"
	configstore "github.com/sevigo/hokan/pkg/store/config"
	"github.com/stretchr/testify/assert"
)

func TestRegisterInitTargetOK(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	configStore := mocks.NewMockConfigStore(controller)
	configStore.EXPECT().Find(gomock.Any(), "local").Return(&core.TargetConfig{
		Name: "local",
		Settings: core.TargetSettings{
			"LOCAL_STORAGE_PATH": "/",
			"LOCAL_BUCKET_NAME":  "test",
		},
		Active: true,
	}, nil)
	configStore.EXPECT().Find(gomock.Any(), "minio").Return(&core.TargetConfig{
		Name:     "minio",
		Active:   false,
		Settings: core.TargetSettings{},
	}, nil)
	configStore.EXPECT().Find(gomock.Any(), "void").Return(nil, configstore.ErrConfigNotFound)
	configStore.EXPECT().Save(gomock.Any(), gomock.Any()).Do(func(_ context.Context, conf *core.TargetConfig) {
		assert.Equal(t, "void", conf.Name)
		assert.Equal(t, true, conf.Active)
		assert.NotEmpty(t, conf.Settings)
	}).Return(nil)

	r := &Register{
		ctx:                    context.TODO(),
		configStore:            configStore,
		storages:               make(map[string]core.TargetStorage),
		storagesStatus:         make(map[string]core.TargetStorageStatus),
		storagesDefaultConfigs: make(map[string]*core.TargetConfig),
		configurators:          make(map[string]core.TargetStorageConfigurator),
	}
	r.initTargets(context.TODO())
	assert.Equal(t, r.GetTarget("local").Name(), "local")
	assert.Equal(t, r.GetTarget("void").Name(), "void")
	assert.Empty(t, r.GetTarget("minio"))

	// test allTargets are not empty
	all := r.AllTargets()
	assert.Equal(t, all["local"].Status, core.TargetStorageOK)
	assert.Equal(t, all["void"].Status, core.TargetStorageOK)
	// TODO: should be not empty, but a different status!
	assert.Empty(t, all["minio"])

	assert.Equal(t, "void", r.GetTargetStorageConfigurator("void").Name())
	assert.Equal(t, "minio", r.GetTargetStorageConfigurator("minio").Name())
}
