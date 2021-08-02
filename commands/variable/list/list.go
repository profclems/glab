package list

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/profclems/glab/api"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/glrepo"
	"github.com/profclems/glab/pkg/iostreams"
	"github.com/profclems/glab/pkg/tableprinter"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

type ListOpts struct {
	HTTPClient func() (*gitlab.Client, error)
	IO         *iostreams.IOStreams
	BaseRepo   func() (glrepo.Interface, error)

	ValueSet bool
}

func NewCmdSet(f *cmdutils.Factory, runE func(opts *ListOpts) error) *cobra.Command {
	opts := &ListOpts{
		IO: f.IO,
	}

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List project or group variables",
		Aliases: []string{"new", "create"},
		Args:    cobra.ExactArgs(0),
		Example: heredoc.Doc(`
			$ glab variable list
		`),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Supports repo override
			opts.HTTPClient = f.HttpClient
			opts.BaseRepo = f.BaseRepo

			if runE != nil {
				err = runE(opts)
				return
			}
			err = listRun(opts)
			return
		},
	}

	return cmd
}

func listRun(opts *ListOpts) error {
	//c := opts.IO.Color()
	httpClient, err := opts.HTTPClient()
	if err != nil {
		return err
	}

	repo, err := opts.BaseRepo()
	if err != nil {
		return err
	}

	createVarOpts := &gitlab.ListProjectVariablesOptions{}

	variables, err := api.ListProjectVariables(httpClient, repo.FullName(), createVarOpts)
	if err != nil {
		return err
	}

	table := tableprinter.NewTablePrinter()
	table.AddRow("KEY", "PROTECTED", "MASKED", "SCOPE")
	for _, variable := range variables {
		table.AddRow(variable.Key, variable.Protected, variable.Masked, variable.EnvironmentScope)
	}

	opts.IO.Log(table.String())
	return nil
}
