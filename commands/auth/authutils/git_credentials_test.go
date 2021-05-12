package authutils

import (
	"testing"
)

func Test_isOurCredentialHelper(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want bool
	}{
		{
			name: "looks like glab but isn't",
			arg:  "glab auth",
			want: false,
		},
		{
			name: "ours",
			arg:  "!/path/to/glab auth",
			want: true,
		},
		{
			name: "blank",
			arg:  "",
			want: false,
		},
		{
			name: "invalid",
			arg:  "!",
			want: false,
		},
		{
			name: "osxkeychain",
			arg:  "osxkeychain",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isOurCredentialHelper(tt.arg); got != tt.want {
				t.Errorf("isOurCredentialHelper() = %v, want %v", got, tt.want)
			}
		})
	}
}
