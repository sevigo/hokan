package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	assert.Equal(t, "0.5.0", Version.String())
}
