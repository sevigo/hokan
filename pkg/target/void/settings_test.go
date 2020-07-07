package void

import (
	"context"
	"testing"

	"github.com/sevigo/hokan/pkg/core"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	configurator := NewConfigurator()

	conf := configurator.DefaultConfig()
	assert.Equal(t, "void", conf.Name)
	assert.Equal(t, true, conf.Active)
}

func TestNewNotActive(t *testing.T) {
	configurator := NewConfigurator()

	conf := configurator.DefaultConfig()
	conf.Active = false
	_, err := New(context.Background(), nil, *conf)
	assert.EqualError(t, err, "target is not active")
}

func TestNewActive(t *testing.T) {
	configurator := NewConfigurator()

	conf := configurator.DefaultConfig()
	_, err := New(context.Background(), nil, *conf)
	assert.NoError(t, err)
}

func Test_voidStore_ValidateSettings(t *testing.T) {
	configurator := NewConfigurator()

	ok, err := configurator.ValidateSettings(core.TargetSettings{})
	assert.True(t, ok)
	assert.NoError(t, err)
}
