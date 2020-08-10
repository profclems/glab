package commands

import (
	"bytes"
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
	"glab/internal/config"
	"glab/internal/git"
	"io"
	"log"
	"os"
	"strings"
)

var repoArchiveCmd = &cobra.Command{
	Use:   "archive <command> [flags]",
	Short: `Get an archive of the repository.`,
	Example: heredoc.Doc(`
	$ glab repo archive profclems/glab
	$ glab repo archive  // Downloads zip file of current repository
	$ glab repo clone profclems/glab mydirectory  // Clones repo into mydirectory
	$ glab repo clone profclems/glab --format=zip   // Finds repo for current user and download in zip format
	`),
	Long: heredoc.Doc(`
	Clone supports these shorthands
	- repo
	- namespace/repo
	- namespace/group/repo
	`),
	Run: func (cmd *cobra.Command, args []string) {
		repo := config.GetRepo()
		var name string
		if len(args) != 0 {
			repo = args[0]
			if len(args) > 1 {
				name = args[1]
			}
		}

		format, _ := cmd.Flags().GetString("format")

		// tar.gz, tar.bz2, tbz, tbz2, tb2, bz2, tar, and zip
		extensions := []string{"tar.gz", "tar.bz2", "tbz", "tbz2", "tb2", "bz2", "tar", "zip"}
		if b := contains(extensions, format); !b {

			fmt.Println("fatal: --format must be one of " + strings.Join(extensions[:], ","))

			return
		}

		gitlabClient, _ := git.InitGitlabClient()
		l := &gitlab.ArchiveOptions{}
		l.Format = gitlab.String(format)
		if sha, _ := cmd.Flags().GetString("sha"); sha != "" {
			l.SHA = gitlab.String(sha)
		}
		ext := *l.Format
		archiveName := strings.Replace(repo, "/", "-", -1) + ext
		if len(strings.TrimSpace(name)) != 0 {
			archiveName = name + "." + ext
		}

		bt, _, err := gitlabClient.Repositories.Archive(repo, l)
		if err != nil {
			log.Fatalf("Failed to clone repos: %v", err)
		}

		r := bytes.NewReader(bt)
		out, err := os.Create(archiveName + ".tmp")
		if err != nil {

			log.Fatalf("Failed to create archive repos: %v", err)
		}

		counter := &CloneWriter{}
		if _, err = io.Copy(out, io.TeeReader(r, counter)); err != nil {
			out.Close()
			log.Fatalf("Failed to write repos: %v", err)
		}

		fmt.Print("\n")
		out.Close()
		if err = os.Rename(archiveName+".tmp", archiveName); err != nil {
			log.Fatalf("Failed to rename tmp repos: %v", err)
		}
		fmt.Println("finish ", archiveName)
	},
}

func init() {
	repoArchiveCmd.Flags().StringP("format", "f", "zip", "Optionally Specify format if you want a downloaded archive: {tar.gz|tar.bz2|tbz|tbz2|tb2|bz2|tar|zip} (Default: zip)")
	repoArchiveCmd.Flags().StringP("sha", "s", "", "The commit SHA to download. A tag, branch reference, or SHA can be used. This defaults to the tip of the default branch if not specified")
	repoCmd.AddCommand(repoArchiveCmd)
}
