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

/*
func cloneRepo(cmdOptions map[string]string, cmdArgs map[int]string) {
	repo := cmdArgs[1]
	dir := cmdArgs[2]

	if len(strings.TrimSpace(repo)) == 0 || !strings.Contains(repo, "/") {

		fmt.Println("fatal: You must specify a owner/repository to clone.")

		return
	}
	if manip.CommandArgExists(cmdOptions, "format") {
		archiveRepo(repo, cmdOptions["format"], dir)
		return
	}

	repos := strings.Split(repo, "/")
	u, _ := url.Parse(config.GetEnv("GITLAB_URI"))
	repoUrl := config.GetEnv("GITLAB_URI") + ":" + config.GetEnv("GITLAB_TOKEN") + "@" + u.Host + "/" + repos[0] + "/" + repos[1]
	// git clone https://gitlab.com:<personal_access_token>@gitlab.com/user/repo.git'
	cmd := exec.Command("git", "clone", repoUrl, dir)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Failed to clone repos: %v", err)
	}
	fmt.Printf("%s\n", out)

}
*/
// mrCmd is merge request command
var repoCmd = &cobra.Command{
	Use:   "repo <command> [flags]",
	Short: `Work with GitLab repositories`,
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 || len(args) > 2 {
			cmd.Help()
			return
		}
	},
}

func init() {
	RootCmd.AddCommand(repoCmd)
}
