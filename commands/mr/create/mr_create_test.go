package create

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/profclems/glab/pkg/iostreams"

	"github.com/MakeNowJust/heredoc"
	"github.com/profclems/glab/commands/cmdtest"
	"github.com/profclems/glab/pkg/prompt"

	"github.com/google/shlex"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/config"
	"github.com/profclems/glab/internal/git"
	"github.com/profclems/glab/internal/glrepo"
	"github.com/profclems/glab/pkg/api"
	"github.com/profclems/glab/pkg/httpmock"
	"github.com/profclems/glab/test"
	"github.com/stretchr/testify/assert"
	"github.com/xanzy/go-gitlab"
)

func runCommand(rt http.RoundTripper, remotes glrepo.Remotes, runE func(opts *CreateOpts) error, branch string, isTTY bool, cli string) (*test.CmdOut, error) {
	io, _, stdout, stderr := iostreams.IOTest()
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
		Remotes: func() (glrepo.Remotes, error) {
			if remotes != nil {
				return remotes, nil
			}
			return glrepo.Remotes{
				{
					Remote: &git.Remote{
						Name:     "origin",
						Resolved: "base",
					},
					Repo: glrepo.New("OWNER", "REPO"),
				},
			}, nil
		},
		Branch: func() (string, error) {
			return branch, nil
		},
		BaseRepo: func() (glrepo.Interface, error) {
			return glrepo.New("OWNER", "REPO"), nil
		},
	}

	// TODO: shouldn't be there but the stub doesn't work without it
	_, _ = factory.HttpClient()

	cmd := NewCmdCreate(factory, runE)
	cmd.PersistentFlags().StringP("repo", "R", "", "")

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

func TestNewCmdCreate_tty(t *testing.T) {
	fakeHTTP := httpmock.New()
	defer fakeHTTP.Verify(t)

	fakeHTTP.RegisterResponder("POST", "/projects/OWNER/REPO/merge_requests",
		httpmock.NewStringResponse(200, `
			{
 				"id": 1,
 				"iid": 12,
 				"project_id": 3,
 				"title": "myMRtitle",
 				"description": "myMRbody",
 				"state": "open",
 				"target_branch": "master",
 				"source_branch": "feat-new-mr",
				"web_url": "https://gitlab.com/OWNER/REPO/-/merge_requests/12"
			}
		`),
	)
	fakeHTTP.RegisterResponder("GET", "/projects/OWNER/REPO",
		httpmock.NewStringResponse(200, `
			{
 				"id": 1,
				"description": null,
				"default_branch": "master",
				"web_url": "http://gitlab.com/OWNER/REPO",
				"name": "OWNER",
				"path": "REPO",
				"merge_requests_enabled": true,
				"path_with_namespace": "OWNER/REPO"
			}
		`),
	)
	fakeHTTP.RegisterResponder("GET", "/users",
		httpmock.NewStringResponse(200, `
			[{
 				"username": "testuser"
			}]
		`),
	)

	ask, teardown := prompt.InitAskStubber()
	defer teardown()

	ask.Stub([]*prompt.QuestionStub{
		{
			Name:  "confirmation",
			Value: 0,
		},
	})

	cs, csTeardown := test.InitCmdStubber()
	defer csTeardown()
	cs.Stub("")
	cs.Stub("HEAD branch: master\n")
	cs.Stub(heredoc.Doc(`
		deadbeef HEAD
		deadb00f refs/remotes/upstream/feat-new-mr
		deadbeef refs/remotes/origin/feat-new-mr
	`))

	cliStr := []string{"-t", "myMRtitle",
		"-d", "myMRbody",
		"-l", "test,bug",
		"--milestone", "1",
		"--assignee", "testuser",
	}

	cli := strings.Join(cliStr, " ")

	runE := func(opts *CreateOpts) error {
		opts.HeadRepo = func() (glrepo.Interface, error) {
			return glrepo.New("OWNER", "REPO"), nil
		}
		return createRun(opts)
	}
	t.Log(cli)

	pu, _ := url.Parse("https://github.com/OWNER/REPO.git")
	t.Log(pu)
	remotes := glrepo.Remotes{
		{
			Remote: &git.Remote{
				Name:     "upstream",
				Resolved: "base",
				PushURL:  pu,
			},
			Repo: glrepo.New("OWNER", "REPO"),
		},
		{
			Remote: &git.Remote{
				Name:     "origin",
				Resolved: "base",
				PushURL:  pu,
			},
			Repo: glrepo.New("monalisa", "REPO"),
		},
	}
	output, err := runCommand(fakeHTTP, remotes, runE, "feat-new-mr", true, cli)
	if err != nil {
		if errors.Is(err, cmdutils.SilentError) {
			t.Errorf("Unexpected error: %q", output.Stderr())
		}
		t.Error(err)
		return
	}

	assert.Contains(t, cmdtest.FirstLine([]byte(output.String())), `!12 myMRtitle (feat-new-mr)`)
	assert.Contains(t, output.Stderr(), "\nCreating merge request for feat-new-mr into master in OWNER/REPO\n\n")
	assert.Contains(t, output.String(), "https://gitlab.com/OWNER/REPO/-/merge_requests/12")
}

func TestMRCreate_nontty_insufficient_flags(t *testing.T) {
	fakeHTTP := httpmock.New()
	defer fakeHTTP.Verify(t)

	_, err := runCommand(fakeHTTP, nil, nil, "test-br", false, "")
	if err == nil {
		t.Fatal("expected error")
	}

	assert.Equal(t, "--title or --fill required for non-interactive mode", err.Error())
}

func TestMrBodyAndTitle(t *testing.T) {
	opts := &CreateOpts{
		SourceBranch:         "mr-autofill-test-br",
		TargetBranch:         "master",
		TargetTrackingBranch: "origin/master",
	}
	t.Run("", func(t *testing.T) {
		cs, csTeardown := test.InitCmdStubber()
		defer csTeardown()
		cs.Stub("d1sd2e,docs: add some changes to txt file")                           // git log
		cs.Stub("Here, I am adding some commit body.\nLittle longer\n\nResolves #1\n") // git log

		if err := mrBodyAndTitle(opts); err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		assert.Equal(t, "docs: add some changes to txt file", opts.Title)
		assert.Equal(t, "Here, I am adding some commit body.\nLittle longer\n\nResolves #1\n", opts.Description)
	})
	t.Run("given-title", func(t *testing.T) {
		cs, csTeardown := test.InitCmdStubber()
		defer csTeardown()

		cs.Stub("d1sd2e,docs: add some changes to txt file")
		cs.Stub("Here, I am adding some commit body.\nLittle longer\n\nResolves #1\n") // git log

		opts := *opts
		opts.Title = "docs: make some other stuff"
		if err := mrBodyAndTitle(&opts); err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		assert.Equal(t, "docs: make some other stuff", opts.Title)
		assert.Equal(t, `Here, I am adding some commit body.
Little longer

Resolves #1
`, opts.Description)
	})
	t.Run("given-description", func(t *testing.T) {
		cs, csTeardown := test.InitCmdStubber()
		defer csTeardown()

		cs.Stub("d1sd2e,docs: add some changes to txt file")

		opts := *opts
		opts.Description = `Make it multiple lines
like this

resolves #1
`
		if err := mrBodyAndTitle(&opts); err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		assert.Equal(t, "docs: add some changes to txt file", opts.Title)
		assert.Equal(t, `Make it multiple lines
like this

resolves #1
`, opts.Description)
	})
}
