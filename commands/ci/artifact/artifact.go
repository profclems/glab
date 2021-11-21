package ci

import (
	"archive/zip"
	"io"
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/profclems/glab/api"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/glrepo"
	"github.com/profclems/glab/pkg/iostreams"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

type TraceOpts struct {
	Branch string
	JobID  int

	BaseRepo   func() (glrepo.Interface, error)
	HTTPClient func() (*gitlab.Client, error)
	IO         *iostreams.IOStreams
}

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

			artifact, _ := api.DownloadArtifactJob(apiClient, repo.FullName(), args[0], &gitlab.DownloadArtifactsFileOptions{Job: &args[1]})

			zr, _ := zip.NewReader(artifact, artifact.Size())
			os.Mkdir(p, 0777)
			for _, v := range zr.File {
				if v.FileInfo().IsDir() {
					os.Mkdir(p+v.Name, 0777)
				} else {
					srcFile, _ := zr.Open(v.Name)
					defer srcFile.Close()
					dstFile, err := os.OpenFile(p+v.Name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, v.Mode())
					if err != nil {
						panic(err)
					}
					io.Copy(dstFile, srcFile)
				}
			}

			return nil
		},
	}
	jobArtifactCmd.Flags().StringP("path", "p", "", "Path to download the Artifact files (default ./)")

	return jobArtifactCmd
}
