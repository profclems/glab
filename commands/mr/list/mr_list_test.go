package list

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"regexp"
	"testing"

	"github.com/profclems/glab/pkg/iostreams"

	"github.com/profclems/glab/commands/cmdtest"

	"github.com/MakeNowJust/heredoc"

	"github.com/alecthomas/assert"
	"github.com/google/shlex"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/config"
	"github.com/profclems/glab/internal/glrepo"
	"github.com/profclems/glab/pkg/api"
	"github.com/profclems/glab/pkg/httpmock"
	"github.com/profclems/glab/test"
	"github.com/xanzy/go-gitlab"
)

func runCommand(rt http.RoundTripper, isTTY bool, cli string, runE func(opts *ListOptions) error) (*test.CmdOut, error) {
	io, _, stdout, stderr := iostreams.Test()
	io.IsaTTY = isTTY
	io.IsInTTY = isTTY
	io.IsErrTTY = isTTY

	factory := &cmdutils.Factory{
		IO: io,
		HttpClient: func() (*gitlab.Client, error) {
			a, err := api.TestClient(&http.Client{Transport: rt}, "", "", false)
			if err != nil {
				return nil, err
			}
			return a.Lab(), err
		},
		Config: func() (config.Config, error) {
			return config.NewBlankConfig(), nil
		},
		BaseRepo: func() (glrepo.Interface, error) {
			return glrepo.New("OWNER", "REPO"), nil
		},
	}

	// TODO: shouldn't be there but the stub doesn't work without it
	_, _ = factory.HttpClient()

	cmd := NewCmdList(factory, runE)

	argv, err := shlex.Split(cli)
	if err != nil {
		return nil, err
	}
	cmd.SetArgs(argv)

	cmd.SetIn(&bytes.Buffer{})
	cmd.SetOut(ioutil.Discard)
	cmd.SetErr(ioutil.Discard)

	_, err = cmd.ExecuteC()
	return &test.CmdOut{
		OutBuf: stdout,
		ErrBuf: stderr,
	}, err
}

func TestNewCmdList(t *testing.T) {
	io, _, _, _ := iostreams.Test()
	io.IsaTTY = true
	io.IsInTTY = true
	io.IsErrTTY = true

	fakeHTTP := httpmock.New()
	defer fakeHTTP.Verify(t)

	factory := &cmdutils.Factory{
		IO: io,
		HttpClient: func() (*gitlab.Client, error) {
			a, err := api.TestClient(&http.Client{Transport: fakeHTTP}, "", "", false)
			if err != nil {
				return nil, err
			}
			return a.Lab(), err
		},
		Config: func() (config.Config, error) {
			return config.NewBlankConfig(), nil
		},
		BaseRepo: func() (glrepo.Interface, error) {
			return glrepo.New("OWNER", "REPO"), nil
		},
	}
	t.Run("MergeRequest_NewCmdList", func(t *testing.T) {
		gotOpts := &ListOptions{}
		err := NewCmdList(factory, func(opts *ListOptions) error {
			gotOpts = opts
			return nil
		}).Execute()

		assert.Nil(t, err)
		assert.Equal(t, factory.IO, gotOpts.IO)

		gotBaseRepo, _ := gotOpts.BaseRepo()
		expectedBaseRepo, _ := factory.BaseRepo()
		assert.Equal(t, gotBaseRepo, expectedBaseRepo)
	})
}

func TestMergeRequestList_tty(t *testing.T) {
	fakeHTTP := httpmock.New()
	defer fakeHTTP.Verify(t)

	fakeHTTP.RegisterResponder("GET", "/projects/OWNER/REPO/merge_requests",
		httpmock.NewStringResponse(200, `
[
  {
    "state" : "opened",
    "description" : "a description here",
    "project_id" : 1,
    "updated_at" : "2016-01-04T15:31:51.081Z",
    "id" : 76,
    "title" : "MergeRequest one",
    "created_at" : "2016-01-04T15:31:51.081Z",
    "iid" : 6,
    "labels" : ["foo", "bar"],
	"target_branch": "master",
    "source_branch": "test1",
    "web_url": "http://gitlab.com/OWNER/REPO/merge_requests/6"
  },
  {
    "state" : "opened",
    "description" : "description two here",
    "project_id" : 1,
    "updated_at" : "2016-01-04T15:31:51.081Z",
    "id" : 77,
    "title" : "MergeRequest two",
    "created_at" : "2016-01-04T15:31:51.081Z",
    "iid" : 7,
	"target_branch": "master",
    "source_branch": "test2",
    "labels" : ["fooz", "baz"],
    "web_url": "http://gitlab.com/OWNER/REPO/merge_requests/7"
  }
]
`))

	output, err := runCommand(fakeHTTP, true, "", nil)
	if err != nil {
		t.Errorf("error running command `issue list`: %v", err)
	}

	out := output.String()
	timeRE := regexp.MustCompile(`\d+ years`)
	out = timeRE.ReplaceAllString(out, "X years")

	assert.Equal(t, heredoc.Doc(`
		Showing 2 open merge requests on OWNER/REPO (Page 1)

		!6	MergeRequest one	(master) ← (test1)
		!7	MergeRequest two	(master) ← (test2)

	`), out)
	assert.Equal(t, ``, output.Stderr())
}

func TestMergeRequestList_tty_withFlags(t *testing.T) {
	fakeHTTP := httpmock.New()
	defer fakeHTTP.Verify(t)

	fakeHTTP.RegisterResponder("GET", "/projects/OWNER/REPO/merge_requests",
		httpmock.NewStringResponse(200, `[]`))

	fakeHTTP.RegisterResponder("GET", "/users",
		httpmock.NewStringResponse(200, `[{"id" : 1, "iid" : 1, "username": "john_smith"}]`))

	output, err := runCommand(fakeHTTP, true, "--opened -P1 -p100 -a someuser -l bug -m1", nil)
	if err != nil {
		t.Errorf("error running command `issue list`: %v", err)
	}

	cmdtest.Eq(t, output.Stderr(), "")
	cmdtest.Eq(t, output.String(), `No open merge requests match your search in OWNER/REPO


`)
}
