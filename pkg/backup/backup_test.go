package backup

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getBackupStorage(t *testing.T) {
	tests := []struct {
		name       string
		backupName string
		wantErr    bool
	}{
		{
			name:       "case 1: backup implementation found",
			backupName: "void",
			wantErr:    false,
		},
		{
			name:       "case 2: wrong backup name",
			backupName: "xxx",
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getBackupStorage(tt.backupName)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
			}
		})
	}
}
