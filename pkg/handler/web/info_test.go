package web_test

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sevigo/hokan/mocks"
	"github.com/sevigo/hokan/pkg/handler/web"
	"github.com/sevigo/hokan/pkg/testing/tools"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestHandleInfo(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	sse := mocks.NewMockServerSideEventCreator(controller)

	s := web.Server{
		Logger: logrus.StandardLogger(),
		SSE:    sse,
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/info", nil)
	s.Handler().ServeHTTP(w, r)
	assert.Equal(t, 200, w.Code)

	body := strings.TrimSpace(w.Body.String())
	tools.TestJSONPathNotEmpty(t, "machine", body)
	tools.TestJSONPathNotEmpty(t, "os", body)
	tools.TestJSONPathNotEmpty(t, "version", body)
}
