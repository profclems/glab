package git

import "testing"

func TestIsValidURL(t *testing.T) {
	type args struct {
		toTest string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Personal repo",args{"profclems/glab"}, false},
		{"Group namespace",args{"group/namespace/repo"}, false},
		{"HTTPS Protocol",args{"https://gitlab.com/profclems/glab.git"}, true},
		{"With SSH",args{"git@gitlab.com:profclems/glab.git"}, true},
		{"SSH Protocol",args{"ssh:user@example.com:my-project"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidURL(tt.args.toTest); got != tt.want {
				t.Errorf("IsValidURL() = %v, want %v", got, tt.want)
			}
		})
	}
}