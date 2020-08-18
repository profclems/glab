package commands

import "github.com/spf13/cobra"

var aliasCmd = &cobra.Command{
	Use:   "alias [command] [flags]",
	Short: `Create, list and delete aliases`,
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	RootCmd.AddCommand(aliasCmd)
}
