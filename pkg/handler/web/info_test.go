package web_test

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sevigo/hokan/pkg/handler/web"
	"github.com/sevigo/hokan/pkg/testing/tools"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestHandleInfo(t *testing.T) {
	s := web.Server{
		Logger: logrus.StandardLogger(),
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/info", nil)
	s.Handler().ServeHTTP(w, r)
	assert.Equal(t, 200, w.Code)

	body := strings.TrimSpace(w.Body.String())
	tools.TestJSONPathNotEmpty(t, "machine", body)
	tools.TestJSONPathNotEmpty(t, "os", body)
	tools.TestJSONPathNotEmpty(t, "user", body)
	tools.TestJSONPathNotEmpty(t, "version", body)
}
