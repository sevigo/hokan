package minio

import (
	"context"
	"testing"

	"github.com/sevigo/hokan/pkg/core"
	"github.com/stretchr/testify/assert"
)

func TestConfigTargetName(t *testing.T) {
	configurator := NewConfigurator()

	assert.Equal(t, "minio", configurator.Name())
}

func TestConfig(t *testing.T) {
	configurator := NewConfigurator()

	conf := configurator.DefaultConfig()
	assert.Equal(t, "minio", conf.Name)
	assert.Equal(t, false, conf.Active)
}

func TestNewNotActive(t *testing.T) {
	configurator := NewConfigurator()

	conf := configurator.DefaultConfig()
	_, err := New(context.Background(), nil, *conf)
	assert.EqualError(t, err, "target is not active")
}

func TestNewActiveErr(t *testing.T) {
	configurator := NewConfigurator()

	conf := configurator.DefaultConfig()
	conf.Active = true
	_, err := New(context.Background(), nil, *conf)
	assert.Error(t, err)
}

func Test_localStorage_ValidateSettings(t *testing.T) {
	configurator := NewConfigurator()
	tests := []struct {
		name     string
		settings core.TargetSettings
		wantErr  bool
	}{
		{
			name: "case 1",
			settings: core.TargetSettings{
				"MINIO_HOST":        "http://localhost:8081",
				"MINIO_ACCESS_KEY":  "abc",
				"MINIO_SECRET_KEY":  "xyz",
				"MINIO_USE_SSL":     "false",
				"MINIO_BUCKET_NAME": "test",
			},
			wantErr: false,
		},
		{
			name: "case 2",
			settings: core.TargetSettings{
				"MINIO_ACCESS_KEY":  "",
				"MINIO_SECRET_KEY":  "",
				"MINIO_USE_SSL":     "",
				"MINIO_BUCKET_NAME": "",
			},
			wantErr: true,
		},
		{
			name: "case 3",
			settings: core.TargetSettings{
				"MINIO_HOST":        "",
				"MINIO_ACCESS_KEY":  "",
				"MINIO_SECRET_KEY":  "",
				"MINIO_USE_SSL":     "",
				"MINIO_BUCKET_NAME": "",
			},
			wantErr: true,
		},
		{
			name: "case 4",
			settings: core.TargetSettings{
				"MINIO_HOST":        "http://localhost:8081",
				"MINIO_ACCESS_KEY":  "abc",
				"MINIO_SECRET_KEY":  "xyz",
				"MINIO_USE_SSL":     "no",
				"MINIO_BUCKET_NAME": "test",
			},
			wantErr: true,
		},
		{
			name: "case 5",
			settings: core.TargetSettings{
				"MINIO_HOST":        "http://localhost:8081",
				"MINIO_ACCESS_KEY":  "abc",
				"MINIO_SECRET_KEY":  "xyz",
				"MINIO_USE_SSL":     "true",
				"MINIO_BUCKET_NAME": "!<.test",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := configurator.ValidateSettings(tt.settings)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
