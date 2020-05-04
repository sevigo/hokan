package directory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_dir_Stats(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "case 1",
			path:    ".",
			wantErr: false,
		},
		{
			name:    "case 2",
			path:    "!",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Stats(tt.path)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
			}
		})
	}
}
