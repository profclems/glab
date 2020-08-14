package commands

import (
	"github.com/spf13/cobra"
)

var mrBrowseCmd = &cobra.Command{
	Use:     "browse <id> [help]",
	Aliases: []string{"b"},
	Short:   "View merge request in a browser",
	Long:    ``,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	mrCmd.AddCommand(mrBrowseCmd)
}
