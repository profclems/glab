package git

import (
	"reflect"
	"testing"
)

func TestCommits(t *testing.T) {
	type args struct {
		baseRef string
		headRef string
	}
	tests := []struct {
		name    string
		args    args
		want    []*Commit
		wantErr bool
	}{
		{
			name: "Commit",
			args: args{"trunk","HEAD"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Commits(tt.args.baseRef, tt.args.headRef)
			if (err != nil) != tt.wantErr {
				t.Errorf("Commits() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Commits() got = %v, want %v", got, tt.want)
			}
		})
	}
}