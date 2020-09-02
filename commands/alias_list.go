package commands

import (
	"fmt"

	"github.com/profclems/glab/internal/config"

	"github.com/spf13/cobra"
)

var aliasListCmd = &cobra.Command{
	Use:   "list [flags]",
	Short: `List the available aliases.`,
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		aliasMap := config.GetAllAliases()

		if len(aliasMap) == 0 {
			fmt.Println("There are currently no aliases")
			fmt.Println("See 'glab alias set --help' for more info")
			return
		}

		for name, command := range aliasMap {
			fmt.Println(name + ": " + command)
		}
	},
}

func init() {
	aliasCmd.AddCommand(aliasListCmd)
}
