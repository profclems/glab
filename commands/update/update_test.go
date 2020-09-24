package update

import "testing"

func Test_isLatestVersion(t *testing.T) {
	type args struct {
		latestVersion  string
		currentVersion string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "latest is newer",
			args: args{"v1.10.0", "v1.9.1"},
			want: true,
		},
		{
			name: "latest is current",
			args: args{"v1.9.2", "v1.9.2"},
			want: false,
		},
		{
			name: "latest is older",
			args: args{"v1.9.0", "v1.9.2-pre.1"},
			want: false,
		},
		{
			name: "current is prerelease",
			args: args{"v1.9.0", "v1.9.0-pre.1"},
			want: true,
		},
		{
			name: "latest is older (against prerelease)",
			args: args{"v1.9.0", "v1.10.0-pre.1"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isLatestVersion(tt.args.latestVersion, tt.args.currentVersion); got != tt.want {
				t.Errorf("isLatestVersion(%s, %s) = %v, want %v",
					tt.args.latestVersion, tt.args.currentVersion, got, tt.want)
			}
		})
	}
}