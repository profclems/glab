package glrepo

import (
	"errors"
	"net/url"
	"testing"

	"github.com/profclems/glab/internal/config"

	"github.com/stretchr/testify/assert"
	"github.com/xanzy/go-gitlab"
)

func Test_RemoteURL(t *testing.T) {
	type args struct {
		project *gitlab.Project
		args    *RemoteArgs
	}

	for _, tt := range []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "is_https",
			args: args{
				project: &gitlab.Project{
					SSHURLToRepo:      "git@gitlab.com:profclems/glab.git",
					HTTPURLToRepo:     "https://gitlab.com/profclems/glab.git",
					PathWithNamespace: "profclems/glab",
				},
				args: &RemoteArgs{
					Protocol: "https",
					Token:    "token",
					Url:      "https://gitlab.com",
					Username: "user",
				},
			},
			want: "https://user:token@gitlab.com/profclems/glab.git",
		},
		{
			name: "host_is_http",
			args: args{
				project: &gitlab.Project{
					SSHURLToRepo:      "git@gitlab.com:profclems/glab.git",
					HTTPURLToRepo:     "http://gitlab.example.com/profclems/glab.git",
					PathWithNamespace: "profclems/glab",
				},
				args: &RemoteArgs{
					Protocol: "https",
					Token:    "token",
					Url:      "http://gitlab.example.com",
					Username: "user",
				},
			},
			want: "http://user:token@gitlab.example.com/profclems/glab.git",
		},
		{
			name: "is_ssh",
			args: args{
				project: &gitlab.Project{
					SSHURLToRepo: "git@gitlab.com:profclems/glab.git",
				},
				args: &RemoteArgs{
					Protocol: "ssh",
				},
			},
			want: "git@gitlab.com:profclems/glab.git",
		},
		{
			name: "no username means oauth2",
			args: args{
				project: &gitlab.Project{
					SSHURLToRepo:      "git@gitlab.com:profclems/glab.git",
					HTTPURLToRepo:     "https://gitlab.com/profclems/glab.git",
					PathWithNamespace: "profclems/glab",
				},
				args: &RemoteArgs{
					Protocol: "https",
					Token:    "token",
					Url:      "https://gitlab.com",
				},
			},
			want: "https://oauth2:token@gitlab.com/profclems/glab.git",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RemoteURL(tt.args.project, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoteURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RemoteURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_repoFromURL(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		result string
		host   string
		err    error
	}{
		{
			name:   "gitlab.com URL",
			input:  "https://gitlab.com/monalisa/octo-cat.git",
			result: "monalisa/octo-cat",
			host:   "gitlab.com",
			err:    nil,
		},
		{
			name:   "gitlab.com URL with trailing slash",
			input:  "https://gitlab.com/monalisa/octo-cat/",
			result: "monalisa/octo-cat",
			host:   "gitlab.com",
			err:    nil,
		},
		{
			name:   "www.gitlab.com URL",
			input:  "http://www.GITLAB.com/monalisa/octo-cat.git",
			result: "monalisa/octo-cat",
			host:   "gitlab.com",
			err:    nil,
		},
		{
			name:   "group namespacing",
			input:  "https://gitlab.com/monalisa/octo-cat/minor",
			result: "monalisa/octo-cat/minor",
			host:   "gitlab.com",
			err:    nil,
		},
		{
			name:   "non-GitLab hostname",
			input:  "https://example.com/one/two",
			result: "one/two",
			host:   "example.com",
			err:    nil,
		},
		{
			name:   "filesystem path",
			input:  "/path/to/file",
			result: "",
			host:   "",
			err:    errors.New("no hostname detected"),
		},
		{
			name:   "filesystem path with scheme",
			input:  "file:///path/to/file",
			result: "",
			host:   "",
			err:    errors.New("no hostname detected"),
		},
		{
			name:   "gitlab.com SSH URL",
			input:  "ssh://gitlab.com/monalisa/octo-cat.git",
			result: "monalisa/octo-cat",
			host:   "gitlab.com",
			err:    nil,
		},
		{
			name:   "gitlab.com HTTPS+SSH URL",
			input:  "https+ssh://gitlab.com/monalisa/octo-cat.git",
			result: "monalisa/octo-cat",
			host:   "gitlab.com",
			err:    nil,
		},
		{
			name:   "gitlab.com git URL",
			input:  "git://gitlab.com/monalisa/octo-cat.git",
			result: "monalisa/octo-cat",
			host:   "gitlab.com",
			err:    nil,
		},
		{
			name:   "gitlab.com deep nested",
			input:  "git://gitlab.com/owner/subgroup/subgroup1/subgroup2/subgroup3/namespace/repo.git",
			result: "owner/subgroup/subgroup1/subgroup2/subgroup3/namespace/repo",
			host:   "gitlab.com",
			err:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := url.Parse(tt.input)
			if err != nil {
				t.Fatalf("got error %q", err)
			}

			repo, err := FromURL(u)
			if err != nil {
				if tt.err == nil {
					t.Fatalf("got error %q", err)
				} else if tt.err.Error() == err.Error() {
					return
				}
				t.Fatalf("got error %q", err)
			}

			got := repo.FullName()
			if tt.result != got {
				t.Errorf("expected %q, got %q", tt.result, got)
			}
			if tt.host != repo.RepoHost() {
				t.Errorf("expected %q, got %q", tt.host, repo.RepoHost())
			}
		})
	}
}

func TestFromFullName(t *testing.T) {
	defer config.StubConfig(`---
hosts:
  gitlab.com:
    token: xxxxxxxxxxxxxxxxxxxx
    git_protocol: ssh
    api_protocol: https
  example.org:
    token: xxxxxxxxxxxxxxxxxxxxx
`, "")()
	tests := []struct {
		name          string
		input         string
		wantOwner     string
		wantName      string
		wantHost      string
		wantFullname  string
		wantGroup     string
		wantNamespace string
		wantErr       error
	}{
		{
			name:          "OWNER/REPO combo",
			input:         "OWNER/REPO",
			wantHost:      "gitlab.com",
			wantOwner:     "OWNER",
			wantName:      "REPO",
			wantFullname:  "OWNER/REPO",
			wantNamespace: "OWNER",
			wantErr:       nil,
		},
		{
			name:    "too few elements",
			input:   "OWNER",
			wantErr: errors.New(`expected the "[HOST/]OWNER/[NAMESPACE/]REPO" format, got "OWNER"`),
		},
		{
			name:          "group namespace",
			input:         "example.org/b/c/d",
			wantHost:      "example.org",
			wantOwner:     "b/c",
			wantName:      "d",
			wantFullname:  "b/c/d",
			wantNamespace: "c",
			wantGroup:     "b",
			wantErr:       nil,
		},
		{
			name:          "with group namespace",
			input:         "gitlab.com/owner/namespace/repo",
			wantHost:      "gitlab.com",
			wantOwner:     "owner/namespace",
			wantName:      "repo",
			wantFullname:  "owner/namespace/repo",
			wantNamespace: "namespace",
			wantGroup:     "owner",
			wantErr:       nil,
		},
		{
			name:    "blank value",
			input:   "a/",
			wantErr: errors.New(`expected the "[HOST/]OWNER/[NAMESPACE/]REPO" format, got "a/"`),
		},
		{
			name:    "blank value inner",
			input:   "a//c",
			wantErr: errors.New(`expected the "[HOST/]OWNER/[NAMESPACE/]REPO" format, got "a//c"`),
		},
		{
			name:          "with hostname",
			input:         "example.org/OWNER/REPO",
			wantHost:      "example.org",
			wantOwner:     "OWNER",
			wantName:      "REPO",
			wantFullname:  "OWNER/REPO",
			wantNamespace: "OWNER",
			wantGroup:     "",
			wantErr:       nil,
		},
		{
			name:          "group name has dot",
			input:         "my.group/sub.group/repo",
			wantHost:      "gitlab.com",
			wantOwner:     "my.group/sub.group",
			wantName:      "repo",
			wantFullname:  "my.group/sub.group/repo",
			wantNamespace: "sub.group",
			wantGroup:     "my.group",
			wantErr:       nil,
		},
		{
			name:          "full URL",
			input:         "https://example.org/OWNER/REPO.git",
			wantHost:      "example.org",
			wantOwner:     "OWNER",
			wantName:      "REPO",
			wantFullname:  "OWNER/REPO",
			wantNamespace: "OWNER",
			wantGroup:     "",
			wantErr:       nil,
		},
		{
			name:          "SSH URL",
			input:         "git@example.org:OWNER/REPO.git",
			wantHost:      "example.org",
			wantOwner:     "OWNER",
			wantName:      "REPO",
			wantFullname:  "OWNER/REPO",
			wantNamespace: "OWNER",
			wantGroup:     "",
			wantErr:       nil,
		},
		{
			name:          "Deep Nested Groups",
			input:         "git@example.org:GROUP/SUBGROUP1/SUBGROUP2/SUBGROUP3/SUBGROUP4/REPO.git",
			wantHost:      "example.org",
			wantOwner:     "GROUP/SUBGROUP1/SUBGROUP2/SUBGROUP3/SUBGROUP4",
			wantName:      "REPO",
			wantFullname:  "GROUP/SUBGROUP1/SUBGROUP2/SUBGROUP3/SUBGROUP4/REPO",
			wantNamespace: "SUBGROUP1/SUBGROUP2/SUBGROUP3/SUBGROUP4",
			wantGroup:     "GROUP",
			wantErr:       nil,
		},
		{
			name:    "invalid URL",
			input:   "git@example.com/%/url",
			wantErr: errors.New(`parse "git@example.com/%/url": invalid URL escape "%/u"`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := FromFullName(tt.input)
			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("no error in result, expected %v", tt.wantErr)
				} else if err.Error() != tt.wantErr.Error() {
					t.Fatalf("expected error %q, got %q", tt.wantErr.Error(), err.Error())
				}
				return
			}
			if err != nil {
				t.Fatalf("got error %v", err)
			}
			if r.RepoHost() != tt.wantHost {
				t.Errorf("expected host %q, got %q", tt.wantHost, r.RepoHost())
			}
			if r.RepoOwner() != tt.wantOwner {
				t.Errorf("expected owner %q, got %q", tt.wantOwner, r.RepoOwner())
			}
			if r.RepoName() != tt.wantName {
				t.Errorf("expected name %q, got %q", tt.wantName, r.RepoName())
			}
			if r.RepoGroup() != tt.wantGroup {
				t.Errorf("expected group %q, got %q", tt.wantGroup, r.RepoGroup())
			}
			if r.FullName() != tt.wantFullname {
				t.Errorf("expected fullname %q, got %q", tt.wantFullname, r.FullName())
			}
			if r.RepoNamespace() != tt.wantNamespace {
				t.Errorf("expected namespace %q, got %q", tt.wantNamespace, r.RepoNamespace())
			}
		})
	}
}

func TestFullNameFromURL(t *testing.T) {

	tests := []struct {
		remoteURL string
		want      string
		wantErr   error
	}{
		{
			remoteURL: "gitlab.com/profclems/glab.git",
			wantErr:   errors.New("cannot parse remote: gitlab.com/profclems/glab.git"),
		},
		{
			remoteURL: "ssh://https://gitlab.com/owner/repo",
			wantErr:   errors.New(`cannot parse remote: ssh://https://gitlab.com/owner/repo`),
		},
		{
			remoteURL: "https://gitlab.com/profclems/glab.git",
			want:      "profclems/glab",
			wantErr:   nil,
		},
		{
			remoteURL: "https://gitlab.com/owner/namespace/repo.git",
			want:      "owner/namespace/repo",
			wantErr:   nil,
		},
		{
			remoteURL: "git@gitlab.com:owner/namespace/repo.git",
			want:      "owner/namespace/repo",
			wantErr:   nil,
		},
		{
			remoteURL: "git@gitlab.com:owner/subgroup/subgroup1/subgroup2/subgroup3/namespace/repo.git",
			want:      "owner/subgroup/subgroup1/subgroup2/subgroup3/namespace/repo",
			wantErr:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.remoteURL, func(t *testing.T) {
			got, err := FullNameFromURL(tt.remoteURL)
			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("no error in result, expected %v", tt.wantErr)
				} else if err.Error() != tt.wantErr.Error() {
					t.Fatalf("expected error %q, got %q", tt.wantErr.Error(), err.Error())
				}
				return
			}
			if got != tt.want {
				t.Errorf("FullNameFromURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_NewWitHost(t *testing.T) {
	testCases := []struct {
		name          string
		input         []string
		wantHostname  string
		wantOwner     string
		wantGroup     string
		wantNamespace string
		wantName      string
		wantFullname  string
	}{
		{
			name:          "Simple",
			input:         []string{"profclems", "glab", "gitlab.com"},
			wantHostname:  "gitlab.com",
			wantNamespace: "profclems",
			wantOwner:     "profclems",
			wantName:      "glab",
			wantFullname:  "profclems/glab",
		},
		{
			name:          "group",
			input:         []string{"company/profclems", "glab", "gitlab.com"},
			wantHostname:  "gitlab.com",
			wantNamespace: "profclems",
			wantOwner:     "company/profclems",
			wantGroup:     "company",
			wantName:      "glab",
			wantFullname:  "company/profclems/glab",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			got := NewWithHost(tC.input[0], tC.input[1], tC.input[2])
			if tC.wantHostname != "" {
				assert.Equal(t, tC.wantHostname, got.RepoHost())
			}
			if tC.wantOwner != "" {
				assert.Equal(t, tC.wantOwner, got.RepoOwner())
			}
			if tC.wantGroup != "" {
				assert.Equal(t, tC.wantGroup, got.RepoGroup())
			}
			if tC.wantNamespace != "" {
				assert.Equal(t, tC.wantNamespace, got.RepoNamespace())
			}
			if tC.wantName != "" {
				assert.Equal(t, tC.wantName, got.RepoName())
			}
			if tC.wantFullname != "" {
				assert.Equal(t, tC.wantFullname, got.FullName())
			}
		})
	}
}
