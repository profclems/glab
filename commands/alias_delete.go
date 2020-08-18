package commands

import (
	"github.com/spf13/cobra"
)

var aliasDeleteCmd = &cobra.Command{
	Use:   "delete <alias name> [flags]",
	Short: `Delete an alias.`,
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	aliasCmd.AddCommand(aliasDeleteCmd)
}
