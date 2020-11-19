package lint

import (
	"testing"

	"github.com/profclems/glab/internal/config"
	"github.com/profclems/glab/internal/utils"
	"github.com/profclems/glab/pkg/api"
	"github.com/stretchr/testify/assert"
	"github.com/xanzy/go-gitlab"

	"github.com/profclems/glab/commands/cmdtest"
)

func TestMain(m *testing.M) {
	cmdtest.InitTest(m, "pipeline_ci_lint_test")
}

func Test_pipelineCILint(t *testing.T) {
	defer config.StubConfig(`---
hosts:
  gitlab.com:
    username: monalisa
    token: OTOKEN
`, "")()

	api.PipelineCILint = func(client *gitlab.Client, content string) (*gitlab.LintResult, error) {
		return &gitlab.LintResult{
			Status: "200",
			Errors: nil,
		}, nil
	}
	io, _, stdout, stderr := utils.IOTest()

	stubFactory, err := cmdtest.StubFactoryWithConfig("")
	assert.Nil(t, err)
	stubFactory.IO = io
	stubFactory.IO.IsaTTY = true
	stubFactory.IO.IsErrTTY = true

	cmd := NewCmdLint(stubFactory)

	_, err = cmd.ExecuteC()
	if err != nil {
		t.Fatal(err)
	}
	assert.Contains(t, stdout.String(), "CI yml is Valid!")
	assert.Equal(t, "", stderr.String())
}
