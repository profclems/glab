package update

import (
	"github.com/spf13/cobra"
	"testing"
)

func TestNewCheckUpdateCmd(t *testing.T) {
	type args struct {
		version string
		build   string
	}
	tests := []struct {
		name string
		args args
		want *cobra.Command
		wantErr bool
	}{
		{
			name: "older version",
			args: args{
				version: "0.0.1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}
