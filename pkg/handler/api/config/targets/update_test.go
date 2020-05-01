package targets_test

import (
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sevigo/hokan/mocks"
	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/handler/api"
	"github.com/sevigo/hokan/pkg/testing/tools"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestHandleActivateOK(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	targets := mocks.NewMockTargetRegister(controller)
	targets.EXPECT().GetConfig(gomock.Any(), "test").Return(&core.TargetConfig{
		Name:        "test",
		Description: "test target",
		Active:      false,
	}, nil)
	targets.EXPECT().SetConfig(gomock.Any(), &core.TargetConfig{
		Name:        "test",
		Description: "test target",
		Active:      true,
	}).Return(nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/config/targets/test/activate", nil)

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
	tools.TestJSONPath(t, "target status changed successfully", "message", body)
}

func TestHandleDeActivateOK(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	targets := mocks.NewMockTargetRegister(controller)
	targets.EXPECT().GetConfig(gomock.Any(), "test").Return(&core.TargetConfig{
		Name:        "test",
		Description: "test target",
		Active:      true,
	}, nil)
	targets.EXPECT().SetConfig(gomock.Any(), &core.TargetConfig{
		Name:        "test",
		Description: "test target",
		Active:      false,
	}).Return(nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/config/targets/test/deactivate", nil)

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
	tools.TestJSONPath(t, "target status changed successfully", "message", body)
}
