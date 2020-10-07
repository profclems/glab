package glrepo

import (
	"errors"
	"net/url"
	"testing"

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
	tests := []struct {
		name      string
		input     string
		wantOwner string
		wantName  string
		wantHost  string
		wantErr   error
	}{
		{
			name:      "OWNER/REPO combo",
			input:     "OWNER/REPO",
			wantHost:  "gitlab.com",
			wantOwner: "OWNER",
			wantName:  "REPO",
			wantErr:   nil,
		},
		{
			name:    "too few elements",
			input:   "OWNER",
			wantErr: errors.New(`expected the "[HOST/]OWNER/[NAMESPACE/]REPO" format, got "OWNER"`),
		},
		{
			name:      "group namespace",
			input:     "a/b/c/d",
			wantHost:  "a",
			wantOwner: "b",
			wantName:  "c/d",
			wantErr:   nil,
		},
		{
			name:      "with group namespace",
			input:     "gitlab.com/owner/namespace/repo",
			wantHost:  "gitlab.com",
			wantOwner: "owner",
			wantName:  "namespace/repo",
			wantErr:   nil,
		},
		{
			name:    "blank value",
			input:   "a/",
			wantErr: errors.New(`expected the "[HOST/]OWNER/[NAMESPACE/]REPO" format, got "a/"`),
		},
		{
			name:      "with hostname",
			input:     "example.org/OWNER/REPO",
			wantHost:  "example.org",
			wantOwner: "OWNER",
			wantName:  "REPO",
			wantErr:   nil,
		},
		{
			name:      "full URL",
			input:     "https://example.org/OWNER/REPO.git",
			wantHost:  "example.org",
			wantOwner: "OWNER",
			wantName:  "REPO",
			wantErr:   nil,
		},
		{
			name:      "SSH URL",
			input:     "git@example.org:OWNER/REPO.git",
			wantHost:  "example.org",
			wantOwner: "OWNER",
			wantName:  "REPO",
			wantErr:   nil,
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
		})
	}
}

func TestFullNameFromURL(t *testing.T) {

	tests := []struct {
		name      string
		remoteURL string
		want      string
		wantErr   error
	}{
		{
			remoteURL: "gitlab.com/profclems/glab.git",
			wantErr:   errors.New("cannot parse remote: gitlab.com/profclems/glab.git"),
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
