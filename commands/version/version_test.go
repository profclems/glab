package version

import (
	"testing"

	"github.com/profclems/glab/pkg/iostreams"

	"github.com/stretchr/testify/assert"
)

func Test_Version(t *testing.T) {
	ios, _, stdout, stderr := iostreams.Test()
	assert.Nil(t, NewCmdVersion(ios, "v1.0.0", "2020-01-01").Execute())

	assert.Equal(t, "glab version 1.0.0 (2020-01-01)\n", stdout.String())
	assert.Equal(t, "", stderr.String())
}
