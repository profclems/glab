package snippet

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/snippet/create"
	"github.com/spf13/cobra"
)

func NewCmdSnippet(f *cmdutils.Factory) *cobra.Command {
	var snippetCmd = &cobra.Command{
		Use:   "snippet <command> [flags]",
		Short: `Create, view and manage snippets`,
		Long:  ``,
		Example: heredoc.Doc(`
			$ glab snippet create --title "Title of the snippet" --filename "main.go"
		`),
		Annotations: map[string]string{
			"help:arguments": heredoc.Doc(`
			A snippet can be supplied as argument in the following format:
			- by number, e.g. "123"
			`),
		},
	}

	cmdutils.EnableRepoOverride(snippetCmd, f)

	snippetCmd.AddCommand(create.NewCmdCreate(f))
	return snippetCmd
}
