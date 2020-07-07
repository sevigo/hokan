package settings_test

import (
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

func TestHandleListOK(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	configurator := mocks.NewMockTargetStorageConfigurator(controller)
	configurator.EXPECT().ValidateSettings(core.TargetSettings{"FOO": "bar"}).Return(true, nil)

	targets := mocks.NewMockTargetRegister(controller)
	targets.EXPECT().GetTargetStorageConfigurator("test").Return(configurator)
	targets.EXPECT().GetConfig(gomock.Any(), "test").Return(&core.TargetConfig{
		Name:        "test",
		Description: "test target",
		Active:      false,
		Settings:    core.TargetSettings{},
	}, nil)
	targets.EXPECT().SetConfig(gomock.Any(), &core.TargetConfig{
		Name:        "test",
		Description: "test target",
		Active:      true,
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

	body := w.Body.String()
	assert.NotEmpty(t, body)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

	tools.TestJSONPath(t, "success", "status", body)
	tools.TestJSONPath(t, "201", "code", body)
	tools.TestJSONPath(t, "target storage config saved successfully", "message", body)
}

func TestHandleList404(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	targets := mocks.NewMockTargetRegister(controller)
	targets.EXPECT().GetTargetStorageConfigurator("test").Return(nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/config/targets/test/settings", strings.NewReader(`{"FOO":"bar"}`))

	s := api.Server{
		Targets: targets,
		Logger:  logrus.StandardLogger(),
	}
	s.Handler().ServeHTTP(w, r)
	assert.Equal(t, 404, w.Code)

	body := w.Body.String()
	assert.NotEmpty(t, body)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

	tools.TestJSONPath(t, "error", "status", body)
	tools.TestJSONPath(t, "404", "code", body)
	tools.TestJSONPath(t, "target not found", "message", body)
}
