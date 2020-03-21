package targets_test

import (
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/sevigo/hokan/mocks"
	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/handler/api"
)

func TestHandleList(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	targets := mocks.NewMockTargetRegister(controller)
	targets.EXPECT().AllConfigs().Return(map[string]*core.TargetConfig{
		"test": &core.TargetConfig{
			Active: true,
			Name:   "test",
		},
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/targets", nil)

	s := api.Server{
		Targets: targets,
		Logger:  logrus.StandardLogger(),
	}
	s.Handler().ServeHTTP(w, r)
	assert.Equal(t, 200, w.Code)
	assert.NotEmpty(t, w.Body.String())
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
}
