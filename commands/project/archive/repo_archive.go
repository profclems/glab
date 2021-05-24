package archive

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/dustin/go-humanize"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/cmdutils/action"
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func NewCmdArchive(f *cmdutils.Factory) *cobra.Command {
	var repoArchiveCmd = &cobra.Command{
		Use:   "archive <command> [flags]",
		Short: `Get an archive of the repository.`,
		Example: heredoc.Doc(`
	$ glab repo archive profclems/glab
	$ glab repo archive  # Downloads zip file of current repository
	$ glab repo archive profclems/glab mydirectory  # Downloads repo zip file into mydirectory
	$ glab repo archive profclems/glab --format=zip   # Finds repo for current user and download in zip format
	`),
		Long: heredoc.Doc(`
	Clone supports these shorthands
	- repo
	- namespace/repo
	- namespace/group/repo
	`),
		Args: cobra.MaximumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			var name string
			var err error

			if len(args) != 0 {
				err = f.RepoOverride(args[0])
				if err != nil {
					return err
				}
				if len(args) > 1 {
					name = args[1]
				}
			}

			apiClient, err := f.HttpClient()
			if err != nil {
				return err
			}

			repo, err := f.BaseRepo()
			if err != nil {
				return err
			}

			format, _ := cmd.Flags().GetString("format")

			// tar.gz, tar.bz2, tbz, tbz2, tb2, bz2, tar, and zip
			extensions := []string{"tar.gz", "tar.bz2", "tbz", "tbz2", "tb2", "bz2", "tar", "zip"}
			if b := contains(extensions, format); !b {
				return errors.New("format must be one of " + strings.Join(extensions, ","))
			}

			l := &gitlab.ArchiveOptions{}
			l.Format = gitlab.String(format)
			if sha, _ := cmd.Flags().GetString("sha"); sha != "" {
				l.SHA = gitlab.String(sha)
			}
			ext := *l.Format
			archiveName := strings.ReplaceAll(repo.FullName(), "/", "-") + "." + ext
			if strings.TrimSpace(name) != "" {
				archiveName = name + "." + ext
			}

			bt, _, err := apiClient.Repositories.Archive(repo.FullName(), l)
			if err != nil {
				return err
			}

			r := bytes.NewReader(bt)
			out, err := os.Create(archiveName + ".tmp")
			if err != nil {
				return fmt.Errorf("failed to create repo archive: %v", err)
			}

			counter := &CloneWriter{}
			if _, err = io.Copy(out, io.TeeReader(r, counter)); err != nil {
				_ = out.Close()
				return fmt.Errorf("failed to write repos: %v", err)
			}

			fmt.Fprint(f.IO.StdOut, "\n")
			_ = out.Close()
			if err = os.Rename(archiveName+".tmp", archiveName); err != nil {
				return fmt.Errorf("failed to rename tmp repos: %v", err)
			}
			fmt.Fprintln(f.IO.StdOut, "Complete...", archiveName)
			return nil
		},
	}

	repoArchiveCmd.Flags().StringP("format", "f", "zip", "Optionally Specify format if you want a downloaded archive: {tar.gz|tar.bz2|tbz|tbz2|tb2|bz2|tar|zip} (Default: zip)")
	repoArchiveCmd.Flags().StringP("sha", "s", "", "The commit SHA to download. A tag, branch reference, or SHA can be used. This defaults to the tip of the default branch if not specified")

	carapace.Gen(repoArchiveCmd).FlagCompletion(carapace.ActionMap{
		"format": carapace.ActionValues("tar.gz", "tar.bz2", "tbz", "tbz2", "tb2", "bz2", "tar", "zip"),
		"sha": carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			branches := action.ActionBranches(repoArchiveCmd, f, &gitlab.ListBranchesOptions{}).Invoke(c)
			tags := action.ActionTags(repoArchiveCmd, f, &gitlab.ListTagsOptions{}).Invoke(c)
			return branches.Merge(tags).ToA() // TODO sha
		}),
	})

	return repoArchiveCmd
}

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
