package target

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sevigo/hokan/mocks"
	"github.com/sevigo/hokan/pkg/core"
	"github.com/stretchr/testify/assert"
)

var allTargets = []string{"void", "local", "minio"}

func TestRegister_GetConfigOK(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	configStore := mocks.NewMockConfigStore(controller)

	configStore.EXPECT().Find(gomock.Any(), "void").Return(&core.TargetConfig{
		Name:   "void",
		Active: true,
	}, nil)

	r := &Register{
		configStore: configStore,
	}

	tConf, err := r.GetConfig(context.TODO(), "void")
	assert.NoError(t, err)
	assert.Equal(t, "void", tConf.Name)
	assert.True(t, tConf.Active)
}

// TODO: test more!
