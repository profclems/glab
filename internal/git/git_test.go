package git

import (
	"os/exec"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/profclems/glab/internal/run"
	"github.com/profclems/glab/test"

	"github.com/stretchr/testify/require"
)

func Test_isFilesystemPath(t *testing.T) {
	type args struct {
		p string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Filesystem",
			args: args{"./.git"},
			want: true,
		},
		{
			name: "Filesystem",
			args: args{".git"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isFilesystemPath(tt.args.p); got != tt.want {
				t.Errorf("isFilesystemPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_UncommittedChangeCount(t *testing.T) {
	type c struct {
		Label    string
		Expected int
		Output   string
	}
	cases := []c{
		{Label: "no changes", Expected: 0, Output: ""},
		{Label: "one change", Expected: 1, Output: " M poem.txt"},
		{Label: "untracked file", Expected: 2, Output: " M poem.txt\n?? new.txt"},
	}

	teardown := run.SetPrepareCmd(func(*exec.Cmd) run.Runnable {
		return &test.OutputStub{}
	})
	defer teardown()

	for _, v := range cases {
		_ = run.SetPrepareCmd(func(*exec.Cmd) run.Runnable {
			return &test.OutputStub{Out: []byte(v.Output)}
		})
		ucc, _ := UncommittedChangeCount()

		if ucc != v.Expected {
			t.Errorf("got unexpected ucc value: %d for case %s", ucc, v.Label)
		}
	}
}

func Test_CurrentBranch(t *testing.T) {
	cs, teardown := test.InitCmdStubber()
	defer teardown()

	expected := "branch-name"

	cs.Stub(expected)

	result, err := CurrentBranch()
	if err != nil {
		t.Errorf("got unexpected error: %v", err)
	}
	if len(cs.Calls) != 1 {
		t.Errorf("expected 1 git call, saw %d", len(cs.Calls))
	}
	if result != expected {
		t.Errorf("unexpected branch name: %s instead of %s", result, expected)
	}
}

func Test_CurrentBranch_detached_head(t *testing.T) {
	cs, teardown := test.InitCmdStubber()
	defer teardown()

	cs.StubError("")

	_, err := CurrentBranch()
	if err == nil {
		t.Errorf("expected an error")
	}
	if err != ErrNotOnAnyBranch {
		t.Errorf("got unexpected error: %s instead of %s", err, ErrNotOnAnyBranch)
	}
	if len(cs.Calls) != 1 {
		t.Errorf("expected 1 git call, saw %d", len(cs.Calls))
	}
}

func Test_CurrentBranch_unexpected_error(t *testing.T) {
	cs, teardown := test.InitCmdStubber()
	defer teardown()

	cs.StubError("lol")

	expectedError := "lol\nstub: lol"

	_, err := CurrentBranch()
	if err == nil {
		t.Errorf("expected an error")
	}
	if err.Error() != expectedError {
		t.Errorf("got unexpected error: %s instead of %s", err.Error(), expectedError)
	}
	if len(cs.Calls) != 1 {
		t.Errorf("expected 1 git call, saw %d", len(cs.Calls))
	}
}

func TestParseExtraCloneArgs(t *testing.T) {
	type Wanted struct {
		args []string
		dir  string
	}
	tests := []struct {
		name string
		args []string
		want Wanted
	}{
		{
			name: "args and target",
			args: []string{"target_directory", "-o", "upstream", "--depth", "1"},
			want: Wanted{
				args: []string{"-o", "upstream", "--depth", "1"},
				dir:  "target_directory",
			},
		},
		{
			name: "only args",
			args: []string{"-o", "upstream", "--depth", "1"},
			want: Wanted{
				args: []string{"-o", "upstream", "--depth", "1"},
				dir:  "",
			},
		},
		{
			name: "only target",
			args: []string{"target_directory"},
			want: Wanted{
				args: []string{},
				dir:  "target_directory",
			},
		},
		{
			name: "no args",
			args: []string{},
			want: Wanted{
				args: []string{},
				dir:  "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args, dir := parseCloneArgs(tt.args)
			got := Wanted{
				args: args,
				dir:  dir,
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %#v want %#v", got, tt.want)
			}
		})
	}

}

func TestReadBranchConfig(t *testing.T) {
	cs, teardown := test.InitCmdStubber()
	defer teardown()

	cs.Stub(`branch.branch-name.remote origin
branch.branch.remote git@gitlab.com:glab-test/test.git
branch.branch.merge refs/heads/branch-name`)

	u, err := ParseURL("git@gitlab.com:glab-test/test.git")
	assert.Nil(t, err)
	wantCfg := BranchConfig{
		"origin",
		u,
		"refs/heads/branch-name",
	}

	t.Run("", func(t *testing.T) {
		if gotCfg := ReadBranchConfig("branch-name"); !reflect.DeepEqual(gotCfg, wantCfg) {
			t.Errorf("ReadBranchConfig() = %v, want %v", gotCfg, wantCfg)
		}
	})
}

func Test_parseRemotes(t *testing.T) {
	remoteList := []string{
		"mona\tgit@gitlab.com:monalisa/myfork.git (fetch)",
		"origin\thttps://gitlab.com/monalisa/octo-cat.git (fetch)",
		"origin\thttps://gitlab.com/monalisa/octo-cat-push.git (push)",
		"upstream\thttps://example.com/nowhere.git (fetch)",
		"upstream\thttps://gitlab.com/hubot/tools (push)",
		"zardoz\thttps://example.com/zed.git (push)",
	}
	r := parseRemotes(remoteList)
	eq(t, len(r), 4)

	eq(t, r[0].Name, "mona")
	eq(t, r[0].FetchURL.String(), "ssh://git@gitlab.com/monalisa/myfork.git")
	if r[0].PushURL != nil {
		t.Errorf("expected no PushURL, got %q", r[0].PushURL)
	}
	eq(t, r[1].Name, "origin")
	eq(t, r[1].FetchURL.Path, "/monalisa/octo-cat.git")
	eq(t, r[1].PushURL.Path, "/monalisa/octo-cat-push.git")

	eq(t, r[2].Name, "upstream")
	eq(t, r[2].FetchURL.Host, "example.com")
	eq(t, r[2].PushURL.Host, "gitlab.com")

	eq(t, r[3].Name, "zardoz")
}

func TestGetDefaultBranch(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    string
		wantErr bool
	}{
		{
			name: "No Params",
			want: "master",
		},
		{
			name: "Different Remote",
			want: "master",
			args: "profclems/test",
		},
		{
			name:    "Invalid repo",
			want:    "master",
			args:    "testssz",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetDefaultBranch(tt.args)
			if (err != nil) != tt.wantErr {
				t.Logf("GetDefaultBranch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetDefaultBranch() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetRemoteURL(t *testing.T) {
	tests := []struct {
		name        string
		remoteAlias string
		want        string
		wantErr     bool
	}{
		{
			name:        "isInvalid",
			remoteAlias: "someorigin",
			wantErr:     true,
		},
		{
			name:        "isInvalid",
			remoteAlias: "origin",
			want:        "profclems/glab",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetRemoteURL(tt.remoteAlias)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRemoteURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.Contains(t, got, tt.want)
		})
	}
}
