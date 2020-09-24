package list

import (
	"fmt"
	"github.com/gosuri/uitable"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/utils"
	"sort"

	"github.com/profclems/glab/internal/config"

	"github.com/spf13/cobra"
)

type ListOptions struct {
	Config func() (config.Config, error)
}

func NewCmdList(f *cmdutils.Factory, runF func(*ListOptions) error) *cobra.Command {
	opts := &ListOptions{
		Config: f.Config,
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

		fmt.Fprintf(utils.ColorableErr(cmd), "no aliases configured\n")
		return nil
	}

	table := uitable.New()
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
	fmt.Fprintf(utils.ColorableOut(cmd), "%v", table)

	return nil
}
