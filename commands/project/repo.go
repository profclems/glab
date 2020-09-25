package project

import (
	"fmt"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
)

var repoCmd = &cobra.Command{
	Use:     "repo <command> [flags]",
	Short:   `Work with GitLab repositories and projects`,
	Long:    ``,
	Aliases: []string{"project"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 || len(args) > 2 {
			_ = cmd.Help()
			return
		}
	},
}

func init() {
	RootCmd.AddCommand(repoCmd)
}
