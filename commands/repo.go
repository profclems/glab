package commands

import (
	"fmt"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
)

//CloneWriter w
type CloneWriter struct {
	Total uint64
}

func (wc *CloneWriter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.progress()
	return n, nil
}

func (wc CloneWriter) progress() {
	fmt.Printf("\r%s", strings.Repeat(" ", 35))
	fmt.Printf("\rCloning... %s complete", humanize.Bytes(wc.Total))
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func fixRepoNamespace(repo string) (string, error) {
	if !strings.Contains(repo, "/") {
		u, err := currentUser()
		if err != nil {
			return "", err
		}
		repo = u + "/" + repo
	}
	return repo, nil
}

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
