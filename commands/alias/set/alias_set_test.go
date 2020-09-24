package set

import (
	"bytes"
	"github.com/MakeNowJust/heredoc"
	"github.com/google/shlex"
	"github.com/profclems/glab/commands/cmdutils"
	"io/ioutil"
	"testing"

	"github.com/profclems/glab/internal/config"
	"github.com/profclems/glab/test"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func runCommand(cfg config.Config, cli string) (*test.CmdOut, error) {

	factoryConf := &cmdutils.Factory{
		Config: func() (config.Config, error) {
			return cfg, nil
		},
	}

	cmd := NewCmdSet(factoryConf, nil)

	// fake command nesting structure needed for validCommand
	rootCmd := &cobra.Command{}
	rootCmd.AddCommand(cmd)
	mrCmd := &cobra.Command{Use: "mr"}
	mrCmd.AddCommand(&cobra.Command{Use: "checkout"})
	mrCmd.AddCommand(&cobra.Command{Use: "rebase"})
	rootCmd.AddCommand(mrCmd)
	issueCmd := &cobra.Command{Use: "issue"}
	issueCmd.AddCommand(&cobra.Command{Use: "list"})
	rootCmd.AddCommand(issueCmd)

	argv, err := shlex.Split("set " + cli)
	if err != nil {
		return nil, err
	}
	rootCmd.SetArgs(argv)
	var stderr bytes.Buffer
	var stdout bytes.Buffer

	rootCmd.SetIn(&bytes.Buffer{})
	rootCmd.SetOut(&stdout)
	rootCmd.SetErr(&stderr)

	_, err = rootCmd.ExecuteC()
	return &test.CmdOut{
		OutBuf: &stdout,
		ErrBuf: &stderr,
	}, err
}

func TestAliasSet_gh_command(t *testing.T) {
	defer config.StubWriteConfig(ioutil.Discard, ioutil.Discard)()

	cfg := config.NewFromString(``)

	_, err := runCommand(cfg,"mr 'mr rebase'")

	if assert.Error(t, err) {
		assert.Equal(t, `could not create alias: "mr" is already a glab command`, err.Error())
	}
}

func TestAliasSet_empty_aliases(t *testing.T) {
	mainBuf := bytes.Buffer{}
	defer config.StubWriteConfig(ioutil.Discard, &mainBuf)()

	cfg := config.NewFromString(heredoc.Doc(`
		aliases:
		editor: vim
	`))

	output, err := runCommand(cfg,"co 'mr checkout'")

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	test.ExpectLines(t, output.String(), "Added alias")
	test.ExpectLines(t, output.Stderr(), "")

	expected := `co: mr checkout
`
	assert.Equal(t, expected, mainBuf.String())
}

func TestAliasSet_existing_alias(t *testing.T) {
	mainBuf := bytes.Buffer{}
	defer config.StubWriteConfig(ioutil.Discard, &mainBuf)()

	cfg := config.NewFromString(heredoc.Doc(`
		aliases:
		  co: mr checkout
	`))

	output, err := runCommand(cfg,"co 'mr checkout -Rcool/repo'")
	require.NoError(t, err)

	test.ExpectLines(t, output.String(), "Changed alias.*co.*from.*mr checkout.*to.*mr checkout -Rcool/repo")
}

func TestAliasSet_space_args(t *testing.T) {
	mainBuf := bytes.Buffer{}
	defer config.StubWriteConfig(ioutil.Discard, &mainBuf)()

	cfg := config.NewFromString(``)

	output, err := runCommand(cfg,`il 'issue list -l "cool story"'`)
	require.NoError(t, err)

	test.ExpectLines(t, output.String(), `Adding alias for.*il.*issue list -l "cool story"`)

	test.ExpectLines(t, mainBuf.String(), `il: issue list -l "cool story"`)
}

func TestAliasSet_arg_processing(t *testing.T) {
	cases := []struct {
		Cmd                string
		ExpectedOutputLine string
		ExpectedConfigLine string
	}{
		{`il "issue list"`, "- Adding alias for.*il.*issue list", "il: issue list"},

		{`iz 'issue list'`, "- Adding alias for.*iz.*issue list", "iz: issue list"},

		{`ii 'issue list --author="$1" --label="$2"'`,
			`- Adding alias for.*ii.*issue list --author="\$1" --label="\$2"`,
			`ii: issue list --author="\$1" --label="\$2"`},

		{`ix "issue list --author='\$1' --label='\$2'"`,
			`- Adding alias for.*ix.*issue list --author='\$1' --label='\$2'`,
			`ix: issue list --author='\$1' --label='\$2'`},
	}

	for _, c := range cases {
		t.Run(c.Cmd, func(t *testing.T) {
			mainBuf := bytes.Buffer{}
			defer config.StubWriteConfig(ioutil.Discard, &mainBuf)()

			cfg := config.NewFromString(``)

			output, err := runCommand(cfg,c.Cmd)
			if err != nil {
				t.Fatalf("got unexpected error running %s: %s", c.Cmd, err)
			}

			test.ExpectLines(t, output.String(), c.ExpectedOutputLine)
			test.ExpectLines(t, mainBuf.String(), c.ExpectedConfigLine)
		})
	}
}

func TestAliasSet_init_alias_cfg(t *testing.T) {
	mainBuf := bytes.Buffer{}
	defer config.StubWriteConfig(ioutil.Discard, &mainBuf)()

	cfg := config.NewFromString(heredoc.Doc(`
		editor: vim
	`))

	output, err := runCommand(cfg,"diff 'mr diff'")
	require.NoError(t, err)

	expected := `diff: mr diff
`

	test.ExpectLines(t, output.String(), "Adding alias for.*diff.*mr diff", "Added alias.")
	assert.Equal(t, expected, mainBuf.String())
}

func TestAliasSet_existing_aliases(t *testing.T) {
	mainBuf := bytes.Buffer{}
	defer config.StubWriteConfig(ioutil.Discard, &mainBuf)()

	cfg := config.NewFromString(heredoc.Doc(`
		aliases:
		  foo: bar
	`))

	output, err := runCommand(cfg,"view 'mr view'")
	require.NoError(t, err)

	expected := `foo: bar
view: mr view
`

	test.ExpectLines(t, output.String(), "Adding alias for.*view.*mr view", "Added alias.")
	assert.Equal(t, expected, mainBuf.String())

}

func TestAliasSet_invalid_command(t *testing.T) {
	defer config.StubWriteConfig(ioutil.Discard, ioutil.Discard)()

	cfg := config.NewFromString(``)

	_, err := runCommand(cfg,"co 'pe checkout'")
	if assert.Error(t, err) {
		assert.Equal(t, "could not create alias: pe checkout does not correspond to a glab command", err.Error())
	}
}

func TestShellAlias_flag(t *testing.T) {
	mainBuf := bytes.Buffer{}
	defer config.StubWriteConfig(ioutil.Discard, &mainBuf)()

	cfg := config.NewFromString(``)

	output, err := runCommand(cfg,"--shell igrep 'glab issue list | grep'")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	test.ExpectLines(t, output.String(), "Adding alias for.*igrep")

	expected := `igrep: '!glab issue list | grep'
`
	assert.Equal(t, expected, mainBuf.String())
}

func TestShellAlias_bang(t *testing.T) {
	mainBuf := bytes.Buffer{}
	defer config.StubWriteConfig(ioutil.Discard, &mainBuf)()

	cfg := config.NewFromString(``)

	output, err := runCommand(cfg,"igrep '!glab issue list | grep'")
	require.NoError(t, err)

	test.ExpectLines(t, output.String(), "Adding alias for.*igrep")

	expected := `igrep: '!glab issue list | grep'
`
	assert.Equal(t, expected, mainBuf.String())
}