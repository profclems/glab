package ci

import (
	"archive/zip"
	"io"
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/profclems/glab/api"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func NewCmdRun(f *cmdutils.Factory) *cobra.Command {
	var jobArtifactCmd = &cobra.Command{
		Use:     "artifact <refName> <jobName> [flags]",
		Short:   `Download all Artifacts from the last pipeline`,
		Aliases: []string{"push"},
		Example: heredoc.Doc(`
	$ glab ci artifact main build
	$ glab ci artifact main deploy --path="artifacts/"
	`),
		Long: ``,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			repo, err := f.BaseRepo()
			if err != nil {
				return err
			}
			apiClient, err := f.HttpClient()
			if err != nil {
				return err
			}
			p, err := cmd.Flags().GetString("path")
			if err != nil {
				return err
			}

			artifact, err := api.DownloadArtifactJob(apiClient, repo.FullName(), args[0], &gitlab.DownloadArtifactsFileOptions{Job: &args[1]})
			if err != nil {
				return err
			}

			zr, err := zip.NewReader(artifact, artifact.Size())
			if err != nil {
				return err
			}

			if err := os.Mkdir(p, 0755); err != nil {
				return err
			}

			for _, v := range zr.File {
				if v.FileInfo().IsDir() {
					if err := os.Mkdir(p+v.Name, v.Mode()); err != nil {
						return err
					}
				} else {
					srcFile, err := zr.Open(v.Name)
					if err != nil {
						return err
					}
					defer srcFile.Close()
					dstFile, err := os.OpenFile(p+v.Name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, v.Mode())
					if err != nil {
						return err
					}
					if _, err := io.Copy(dstFile, srcFile); err != nil {
						return err
					}
				}
			}
			return nil
		},
	}
	jobArtifactCmd.Flags().StringP("path", "p", "", "Path to download the Artifact files (default ./)")

	return jobArtifactCmd
}
