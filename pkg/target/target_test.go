package target

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sevigo/hokan/mocks"
	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/target/void"
	"github.com/stretchr/testify/assert"
)

func TestRegister_initTarget(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	configStore := mocks.NewMockConfigStore(controller)
	configStore.EXPECT().Find(gomock.Any(), "void").Return(&core.TargetConfig{
		Name:   "void",
		Active: true,
	}, nil)

	targetConfigurator := void.NewConfigurator()
	r := &Register{
		ctx: context.TODO(),
		// fileStore:      fileStore,
		configStore: configStore,
		// event:          event,
		storages:               make(map[string]core.TargetStorage),
		storagesStatus:         make(map[string]core.TargetStorageStatus),
		storagesDefaultConfigs: make(map[string]*core.TargetConfig),

		Results: make(chan core.TargetOperationResult),
	}
	tests := []struct {
		name       string
		targetName string
		wantErr    bool
	}{
		{
			name:       "case 1",
			targetName: "void",
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := r.initTarget(context.TODO(), targetConfigurator)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
