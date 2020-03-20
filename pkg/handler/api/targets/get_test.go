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

	body := strings.TrimSpace(w.Body.String())
	assert.Equal(t, `{"Active":true,"Name":"test","Description":"this is test","Settings":null}`, body)
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

	body := strings.TrimSpace(w.Body.String())
	assert.Equal(t, `{"code":404,"message":"default config for target not found"}`, body)
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
	assert.Equal(t, 500, w.Code)

	body := strings.TrimSpace(w.Body.String())
	assert.Equal(t, `{"code":500,"message":"Some error"}`, body)
}
