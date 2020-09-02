package commands

import (
	"fmt"
	"os"

	"github.com/profclems/glab/internal/config"

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
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	aliasCmd.AddCommand(aliasDeleteCmd)
}
