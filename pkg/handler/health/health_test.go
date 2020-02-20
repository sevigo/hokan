package health

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleHealthz(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/healthz", nil)

	Handler().ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)
}
