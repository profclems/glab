package note

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/profclems/glab/pkg/iostreams"

	"github.com/alecthomas/assert"
	"github.com/google/shlex"
	"github.com/profclems/glab/api"
	"github.com/profclems/glab/commands/cmdtest"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/config"
	"github.com/profclems/glab/internal/git"
	"github.com/profclems/glab/internal/glrepo"
	"github.com/profclems/glab/pkg/httpmock"
	"github.com/profclems/glab/pkg/prompt"
	"github.com/profclems/glab/test"
	"github.com/xanzy/go-gitlab"
)

func TestMain(m *testing.M) {
	cmdtest.InitTest(m, "mr_note_create_test")
}

func runCommand(rt http.RoundTripper, isTTY bool, cli string) (*test.CmdOut, error) {
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
		Branch: git.CurrentBranch,
	}

	// TODO: shouldn't be there but the stub doesn't work without it
	_, _ = factory.HttpClient()

	cmd := NewCmdNote(factory)

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

func Test_NewCmdNote(t *testing.T) {
	fakeHTTP := httpmock.New()
	defer fakeHTTP.Verify(t)

	t.Run("--message flag specified", func(t *testing.T) {

		fakeHTTP.RegisterResponder("POST", "/projects/OWNER/REPO/merge_requests/1/notes",
			httpmock.NewStringResponse(201, `
		{
			"id": 301,
  			"created_at": "2013-10-02T08:57:14Z",
  			"updated_at": "2013-10-02T08:57:14Z",
  			"system": false,
  			"noteable_id": 1,
  			"noteable_type": "MergeRequest",
  			"noteable_iid": 1
		}
	`))

		fakeHTTP.RegisterResponder("GET", "/projects/OWNER/REPO/merge_requests/1",
			httpmock.NewStringResponse(200, `
		{
  			"id": 1,
  			"iid": 1,
			"web_url": "https://gitlab.com/OWNER/REPO/merge_requests/1"
		}
	`))

		// glab mr note 1 --message "Here is my note"
		output, err := runCommand(fakeHTTP, true, `1 --message "Here is my note"`)
		if err != nil {
			t.Error(err)
			return
		}
		assert.Equal(t, output.Stderr(), "")
		assert.Equal(t, output.String(), "https://gitlab.com/OWNER/REPO/merge_requests/1#note_301\n")
	})

	t.Run("merge request not found", func(t *testing.T) {
		fakeHTTP.RegisterResponder("GET", "/projects/OWNER/REPO/merge_requests/122",
			httpmock.NewStringResponse(404, `
		{
  			"message" : "merge request not found"
		}
	`))

		// glab mr note 1 --message "Here is my note"
		_, err := runCommand(fakeHTTP, true, `122`)
		assert.NotNil(t, err)
		assert.Equal(t, "failed to get merge request 122: GET https://gitlab.com/api/v4/projects/OWNER/REPO/merge_requests/122: 404 {message: merge request not found}", err.Error())
	})
}

func Test_NewCmdNote_error(t *testing.T) {
	fakeHTTP := httpmock.New()
	defer fakeHTTP.Verify(t)

	t.Run("note could not be created", func(t *testing.T) {
		fakeHTTP.RegisterResponder("POST", "/projects/OWNER/REPO/merge_requests/1/notes",
			httpmock.NewStringResponse(401, `
		{
			"message" : "Unauthorized"
		}
	`))

		fakeHTTP.RegisterResponder("GET", "/projects/OWNER/REPO/merge_requests/1",
			httpmock.NewStringResponse(200, `
		{
  			"id": 1,
  			"iid": 1,
			"web_url": "https://gitlab.com/OWNER/REPO/merge_requests/1"
		}
	`))

		// glab mr note 1 --message "Here is my note"
		_, err := runCommand(fakeHTTP, true, `1 -m "Some message"`)
		assert.NotNil(t, err)
		assert.Equal(t, "POST https://gitlab.com/api/v4/projects/OWNER/REPO/merge_requests/1/notes: 401 {message: Unauthorized}", err.Error())
	})
}

func Test_mrNoteCreate_prompt(t *testing.T) {
	fakeHTTP := httpmock.New()
	defer fakeHTTP.Verify(t)

	t.Run("message provided", func(t *testing.T) {

		fakeHTTP.RegisterResponder("POST", "/projects/OWNER/REPO/merge_requests/1/notes",
			httpmock.NewStringResponse(201, `
		{
			"id": 301,
  			"created_at": "2013-10-02T08:57:14Z",
  			"updated_at": "2013-10-02T08:57:14Z",
  			"system": false,
  			"noteable_id": 1,
  			"noteable_type": "MergeRequest",
  			"noteable_iid": 1
		}
	`))

		fakeHTTP.RegisterResponder("GET", "/projects/OWNER/REPO/merge_requests/1",
			httpmock.NewStringResponse(200, `
		{
  			"id": 1,
  			"iid": 1,
			"web_url": "https://gitlab.com/OWNER/REPO/merge_requests/1"
		}
	`))
		as, teardown := prompt.InitAskStubber()
		defer teardown()
		as.StubOne("some note message")

		// glab mr note 1
		output, err := runCommand(fakeHTTP, true, `1`)
		if err != nil {
			t.Error(err)
			return
		}
		assert.Equal(t, output.Stderr(), "")
		assert.Equal(t, output.String(), "https://gitlab.com/OWNER/REPO/merge_requests/1#note_301\n")
	})

	t.Run("message is empty", func(t *testing.T) {

		fakeHTTP.RegisterResponder("GET", "/projects/OWNER/REPO/merge_requests/1",
			httpmock.NewStringResponse(200, `
		{
  			"id": 1,
  			"iid": 1,
			"web_url": "https://gitlab.com/OWNER/REPO/merge_requests/1"
		}
	`))

		as, teardown := prompt.InitAskStubber()
		defer teardown()
		as.StubOne("")

		// glab mr note 1
		_, err := runCommand(fakeHTTP, true, `1`)
		if err == nil {
			t.Error("expected error")
			return
		}
		assert.Equal(t, err.Error(), "aborted... Note has an empty message")
	})
}
