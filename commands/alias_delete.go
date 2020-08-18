package commands

import (
	"glab/internal/config"
	"log"

	"github.com/spf13/cobra"
)

var aliasDeleteCmd = &cobra.Command{
	Use:   "delete <alias name> [flags]",
	Short: `Delete an alias.`,
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := config.DeleteAlias(args[0])
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	aliasCmd.AddCommand(aliasDeleteCmd)
}
