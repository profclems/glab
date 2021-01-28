package list

import (
	"fmt"
	"sort"

	"github.com/profclems/glab/pkg/iostreams"

	"github.com/profclems/glab/pkg/tableprinter"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/config"

	"github.com/spf13/cobra"
)

type ListOptions struct {
	Config func() (config.Config, error)
	IO     *iostreams.IOStreams
}

func NewCmdList(f *cmdutils.Factory, runF func(*ListOptions) error) *cobra.Command {
	opts := &ListOptions{
		Config: f.Config,
		IO:     f.IO,
	}

	var aliasListCmd = &cobra.Command{
		Use:   "list [flags]",
		Short: `List the available aliases.`,
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}
			return listRun(cmd, opts)
		},
	}
	return aliasListCmd
}

func listRun(cmd *cobra.Command, opts *ListOptions) error {
	cfg, err := opts.Config()
	if err != nil {
		return err
	}

	aliasCfg, err := cfg.Aliases()
	if err != nil {
		return fmt.Errorf("couldn't read aliases config: %w", err)
	}

	if aliasCfg.Empty() {

		fmt.Fprintf(opts.IO.StdErr, "no aliases configured\n")
		return nil
	}

	table := tableprinter.NewTablePrinter()
	table.MaxColWidth = 70

	aliasMap := aliasCfg.All()
	var keys []string
	for alias := range aliasMap {
		keys = append(keys, alias)
	}
	sort.Strings(keys)

	for _, alias := range keys {
		table.AddRow(alias, aliasMap[alias])
	}
	fmt.Fprintf(opts.IO.StdOut, "%s", table.Render())

	return nil
}
