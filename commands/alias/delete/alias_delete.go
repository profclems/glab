package delete

import (
	"fmt"

	"github.com/profclems/glab/pkg/iostreams"

	"github.com/profclems/glab/commands/cmdutils"

	"github.com/profclems/glab/internal/config"
	"github.com/spf13/cobra"
)

type DeleteOptions struct {
	Config func() (config.Config, error)
	Name   string
	IO     *iostreams.IOStreams
}

func NewCmdDelete(f *cmdutils.Factory, runF func(*DeleteOptions) error) *cobra.Command {
	opts := &DeleteOptions{
		Config: f.Config,
		IO:     f.IO,
	}

	var aliasDeleteCmd = &cobra.Command{
		Use:   "delete <alias name> [flags]",
		Short: `Delete an alias.`,
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Name = args[0]
			if runF != nil {
				return runF(opts)
			}
			return deleteRun(cmd, opts)
		},
	}
	return aliasDeleteCmd
}

func deleteRun(cmd *cobra.Command, opts *DeleteOptions) error {
	cfg, err := opts.Config()
	if err != nil {
		return err
	}

	aliasCfg, err := cfg.Aliases()
	if err != nil {
		return fmt.Errorf("couldn't read aliases config: %w", err)
	}

	expansion, ok := aliasCfg.Get(opts.Name)
	if !ok {
		return fmt.Errorf("no such alias %s", opts.Name)

	}
	err = aliasCfg.Delete(opts.Name)
	if err != nil {
		return fmt.Errorf("failed to delete alias %s: %w", opts.Name, err)
	}
	redCheck := iostreams.Red("✓")
	fmt.Fprintf(opts.IO.StdErr, "%s Deleted alias %s; was %s\n", redCheck, opts.Name, expansion)
	return nil
}
