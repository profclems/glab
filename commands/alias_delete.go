package commands

import (
	"fmt"
	"github.com/profclems/glab/internal/config"
	"github.com/profclems/glab/internal/utils"

	"github.com/spf13/cobra"
)

var aliasDeleteCmd = &cobra.Command{
	Use:   "delete <alias name> [flags]",
	Short: `Delete an alias.`,
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		expansion := config.GetAlias(args[0])
		err := config.DeleteAlias(args[0])

		if err != nil {
			return fmt.Errorf("failed to delete alias %s: %w", args[0], err)
		}
		out := colorableOut(cmd)
		redCheck := utils.Red("✓")
		fmt.Fprintf(out, "%s Deleted alias %s; was %s\n", redCheck, args[0], expansion)
		return nil
	},
}

func init() {
	aliasCmd.AddCommand(aliasDeleteCmd)
}
