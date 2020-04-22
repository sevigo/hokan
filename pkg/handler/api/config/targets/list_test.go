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
	"github.com/sevigo/hokan/pkg/testing/tools"
)

func TestHandleList(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	targets := mocks.NewMockTargetRegister(controller)
	targets.EXPECT().AllConfigs().Return(map[string]core.TargetConfig{
		"test": {
			Active:      true,
			Name:        "test",
			Description: "this is test",
			Settings: map[string]string{
				"FOO": "bar",
			},
		},
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/config/targets", nil)

	s := api.Server{
		Targets: targets,
		Logger:  logrus.StandardLogger(),
	}
	s.Handler().ServeHTTP(w, r)
	assert.Equal(t, 200, w.Code)
	body := w.Body.String()
	assert.NotEmpty(t, body)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

	tools.TestJSONPath(t, "true", "targets.test.active", body)
	tools.TestJSONPath(t, "test", "targets.test.name", body)
	tools.TestJSONPath(t, "this is test", "targets.test.description", body)
	tools.TestJSONPath(t, "bar", "targets.test.settings.FOO", body)
}
