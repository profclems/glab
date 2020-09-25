package list

import (
	"fmt"
	"github.com/profclems/glab/pkg/api"
	"strings"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/utils"

	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func NewCmdList(f *cmdutils.Factory) *cobra.Command {
	var labelListCmd = &cobra.Command{
		Use:     "list [flags]",
		Short:   `List labels in repository`,
		Long:    ``,
		Aliases: []string{"ls"},
		Args:    cobra.ExactArgs(0),
		RunE:    func(cmd *cobra.Command, args []string) error {
			var err error
			out := utils.ColorableOut(cmd)
			if r, _ := cmd.Flags().GetString("repo"); r != "" {
				f, err = f.NewClient(r)
				if err != nil {
					return err
				}
			}
			apiClient, err := f.HttpClient()
			cfg, _ := f.Config()
			if err != nil {
				return err
			}
			repo, err := f.BaseRepo()
			if err != nil {
				return err
			}

			l := &gitlab.ListLabelsOptions{}
			if p, _ := cmd.Flags().GetInt("page"); p != 0 {
				l.Page = p
			}
			if p, _ := cmd.Flags().GetInt("per-page"); p != 0 {
				l.PerPage = p
			}

			// List all labels
			labels, err := api.ListLabels(apiClient, repo.FullName(), l)
			if err != nil {
				return err
			}
			fmt.Fprintf(out, "Showing label %d of %d on %s\n\n", len(labels), len(labels), repo.FullName())
			var labelPrintInfo string
			for _, label := range labels {
				labelPrintInfo += label.Name
				if label.Description != "" {
					labelPrintInfo += " -> " + label.Description
				}
				labelPrintInfo += "\n"
			}
			fmt.Fprintf(out, utils.Indent(labelPrintInfo, " "))

			// Cache labels for host
			labelNames := make([]string, 0, len(labels))
			for _, label := range labels {
				labelNames = append(labelNames, label.Name)
			}
			labelsEntry := strings.Join(labelNames, ",")
			if err := cfg.Set(repo.RepoHost(), "labels", labelsEntry); err != nil {
				_ = cfg.Write()
			}

			return nil

		},
	}

	labelListCmd.Flags().IntP("page", "p", 1, "Page number")
	labelListCmd.Flags().IntP("per-page", "P", 20, "Number of items to list per page")

	return labelListCmd
}
