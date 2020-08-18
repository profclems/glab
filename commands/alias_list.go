package commands

import (
	"github.com/spf13/cobra"
)

var aliasListCmd = &cobra.Command{
	Use:   "list [flags]",
	Short: `List the available aliases.`,
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	aliasCmd.AddCommand(aliasListCmd)
}
