package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_defaultMinIO(t *testing.T) {
	c := &Config{
		Backup: Backup{
			configName: "config.exampl.yaml",
		},
	}
	defaultMinIO(c)
	assert.Equal(t, "minio", c.Backup.Name)
	assert.Equal(t, "localhost", c.Backup.MinIO.Endpoint)
	assert.Equal(t, "secret", c.Backup.MinIO.SecretAccessKey)
	assert.Equal(t, "keyid", c.Backup.MinIO.AccessKeyID)
	assert.Equal(t, false, c.Backup.MinIO.UseSSL)
}
