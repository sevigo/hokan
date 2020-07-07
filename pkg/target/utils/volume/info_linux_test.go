// +build linux

package volume

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetVolumeInformation(t *testing.T) {
	free, total := GetVolumeInformation(context.TODO(), `/`)
	assert.NotEmpty(t, free)
	assert.NotEmpty(t, total)
}
