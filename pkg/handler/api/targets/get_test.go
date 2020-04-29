package targets_test

import (
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sevigo/hokan/mocks"
	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/handler/api"
	"github.com/sevigo/hokan/pkg/testing/tools"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

const targetOK = "test"
const targetBad = "bad"

func TestGetTargetByNameOK(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	targets := mocks.NewMockTargetRegister(controller)
	targets.EXPECT().GetConfig(gomock.Any(), targetOK).Return(&core.TargetConfig{
		Active:      true,
		Name:        targetOK,
		Description: "this is test",
	}, nil)

	url := fmt.Sprintf("/targets/%s", targetOK)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", url, nil)

	s := api.Server{
		Targets: targets,
		Logger:  logrus.StandardLogger(),
	}
	s.Handler().ServeHTTP(w, r)
	assert.Equal(t, 200, w.Code)

	assert.NotEmpty(t, w.Body.String())
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
}

func TestGetTargetByNameNotFound(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	targets := mocks.NewMockTargetRegister(controller)
	targets.EXPECT().GetConfig(gomock.Any(), targetBad).Return(nil, core.ErrTargetConfigNotFound)

	url := fmt.Sprintf("/targets/%s", targetBad)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", url, nil)

	s := api.Server{
		Targets: targets,
		Logger:  logrus.StandardLogger(),
	}
	s.Handler().ServeHTTP(w, r)

	assert.Equal(t, 404, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

	body := strings.TrimSpace(w.Body.String())
	tools.TestJSONPath(t, "404", "code", body)
	tools.TestJSONPath(t, "default config for target not found", "message", body)
	// tools.TestJSONPath(t, "error", "status", body)
}

func TestGetTargetByNameError(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	targets := mocks.NewMockTargetRegister(controller)
	targets.EXPECT().GetConfig(gomock.Any(), targetBad).Return(nil, fmt.Errorf("Some error"))

	url := fmt.Sprintf("/targets/%s", targetBad)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", url, nil)

	s := api.Server{
		Targets: targets,
		Logger:  logrus.StandardLogger(),
	}
	s.Handler().ServeHTTP(w, r)

	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

	body := strings.TrimSpace(w.Body.String())
	tools.TestJSONPath(t, "400", "code", body)
	tools.TestJSONPath(t, "Some error", "message", body)
	// tools.TestJSONPath(t, "error", "status", body)
}
