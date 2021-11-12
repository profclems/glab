package update

import (
	"fmt"
	"testing"

	"github.com/profclems/glab/pkg/iostreams"

	"github.com/alecthomas/assert"
	"github.com/jarcoal/httpmock"
)

func TestNewCheckUpdateCmd(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", `https://api.github.com/repos/profclems/glab/releases/latest`,
		httpmock.NewStringResponder(200, `{
    "url": "https://api.github.com/repos/profclems/glab/releases/33385584",
  "html_url": "https://github.com/profclems/glab/releases/tag/v1.11.1",
  "tag_name": "v1.11.1",
  "name": "v1.11.1",
  "draft": false,
  "prerelease": false,
  "created_at": "2020-11-03T05:33:29Z",
  "published_at": "2020-11-03T05:39:04Z"}`))

	ioStream, _, stdout, stderr := iostreams.Test()
	type args struct {
		s       *iostreams.IOStreams
		version string
	}
	tests := []struct {
		name    string
		args    args
		stdOut  string
		stdErr  string
		wantErr bool
	}{
		{
			name: "same version",
			args: args{
				s:       ioStream,
				version: "v1.11.1",
			},
			stdOut: "✓ You are already using the latest version of glab\n",
			stdErr: "",
		},
		{
			name: "older version",
			args: args{
				s:       ioStream,
				version: "v1.11.0",
			},
			stdOut: "A new version of glab has been released: v1.11.0 → v1.11.1\nhttps://github.com/profclems/glab/releases/tag/v1.11.1\n",
			stdErr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewCheckUpdateCmd(tt.args.s, tt.args.version).Execute()
			if tt.wantErr {
				assert.Nil(t, err)
			}

			assert.Equal(t, tt.stdOut, stdout.String())
			assert.Equal(t, tt.stdErr, stderr.String())

			// clean up
			stdout.Reset()
			stderr.Reset()
		})
	}
}

func TestNewCheckUpdateCmd_error(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", `https://api.github.com/repos/profclems/glab/releases/latest`,
		httpmock.NewErrorResponder(fmt.Errorf("an error expected")))

	ioStream, _, stdout, stderr := iostreams.Test()

	err := NewCheckUpdateCmd(ioStream, "1.11.0").Execute()
	assert.NotNil(t, err)
	assert.Equal(t, "could not check for update! Make sure you have a stable internet connection", err.Error())
	assert.Equal(t, "", stdout.String())
	assert.Equal(t, "", stderr.String())
}

func TestNewCheckUpdateCmd_json_error(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", `https://api.github.com/repos/profclems/glab/releases/latest`,
		httpmock.NewStringResponder(200, ``))

	ioStream, _, stdout, stderr := iostreams.Test()

	err := NewCheckUpdateCmd(ioStream, "1.11.0").Execute()
	assert.NotNil(t, err)
	assert.Equal(t, "could not check for update! Make sure you have a stable internet connection", err.Error())
	assert.Equal(t, "", stdout.String())
	assert.Equal(t, "", stderr.String())
}
