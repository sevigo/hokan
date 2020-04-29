package settings_test

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sevigo/hokan/mocks"
	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/handler/api"
	"github.com/sirupsen/logrus"
	"gotest.tools/assert"
)

func TestHandleListOK(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	target := mocks.NewMockTargetStorage(controller)
	target.EXPECT().ValidateSettings(core.TargetSettings{"FOO": "bar"}).Return(true, nil)

	targets := mocks.NewMockTargetRegister(controller)
	targets.EXPECT().GetTarget("test").Return(target)
	targets.EXPECT().GetConfig(gomock.Any(), "test").Return(&core.TargetConfig{
		Name:        "test",
		Description: "test target",
		Active:      false,
		Settings:    core.TargetSettings{},
	}, nil)
	targets.EXPECT().SetConfig(gomock.Any(), &core.TargetConfig{
		Name:        "test",
		Description: "test target",
		Active:      false,
		Settings: core.TargetSettings{
			"FOO": "bar",
		},
	}).Return(nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/config/targets/test/settings", strings.NewReader(`{"FOO":"bar"}`))

	s := api.Server{
		Targets: targets,
		Logger:  logrus.StandardLogger(),
	}
	s.Handler().ServeHTTP(w, r)
	assert.Equal(t, 201, w.Code)
}

func TestHandleList404(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	targets := mocks.NewMockTargetRegister(controller)
	targets.EXPECT().GetTarget("test").Return(nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/config/targets/test/settings", strings.NewReader(`{"FOO":"bar"}`))

	s := api.Server{
		Targets: targets,
		Logger:  logrus.StandardLogger(),
	}
	s.Handler().ServeHTTP(w, r)
	assert.Equal(t, 404, w.Code)
}
