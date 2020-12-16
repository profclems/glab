// This package contains the old `glab pipeline ci` command which has been deprecated
// in favour of the `glab ci` command.
// This package is kept for backward compatibility but issues a deprecation warning
package legacyci

import (
	"testing"

	"github.com/profclems/glab/internal/utils"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/stretchr/testify/assert"
)

func TestNewCmdCI(t *testing.T) {
	ioStrm, stdin, stdout, stderr := utils.IOTest()

	cmd := NewCmdCI(&cmdutils.Factory{
		IO: ioStrm,
	})

	cmd.SetIn(stdin)
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)

	_, err := cmd.ExecuteC()

	assert.Nil(t, err)

	assert.Contains(t, stdout.String(), "Work with GitLab CI pipelines and jobs\n")
	assert.Contains(t, stderr.String(), "")
	assert.Contains(t, stdout.String(), "This command is deprecated. All the commands under it has been moved to `ci` or `pipeline` command\n")

}
