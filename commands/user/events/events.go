package events

import (
	"fmt"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/utils"
	"github.com/profclems/glab/pkg/api"
	"github.com/profclems/glab/pkg/tableprinter"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func NewCmdEvents(f *cmdutils.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "events",
		Short: "View user events",
		Args:  cobra.MaximumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			apiClient, err := f.HttpClient()
			if err != nil {
				return err
			}

			events, err := api.CurrentUserEvents(apiClient)
			if err != nil {
				return err
			}
			projects := make(map[int]*gitlab.Project)
			for _, e := range events {
				project, err := api.GetProject(apiClient, e.ProjectID)
				if err != nil {
					return err
				}
				projects[e.ProjectID] = project
			}

			title := utils.NewListTitle("User events")
			title.CurrentPageTotal = len(events)

			if err = f.IO.StartPager(); err != nil {
				return err
			}
			defer f.IO.StopPager()
			fmt.Fprintf(f.IO.StdOut, "%s\n%s\n", title.Describe(), DisplayAllEvents(events, projects))

			return nil
		},
	}

	return cmd
}

func DisplayAllEvents(events []*gitlab.ContributionEvent, projects map[int]*gitlab.Project) string {
	table := tableprinter.NewTablePrinter()
	for _, e := range events {
		table.AddCell(e.ActionName)
		table.AddCell(projects[e.ProjectID].Name)
		if e.ActionName == "pushed to" {
			table.AddCell(e.PushData.Ref)
			table.AddCell(e.PushData.CommitTitle)
		} else if e.ActionName == "commented on" {
			table.AddCell(e.Note.NoteableType)
			table.AddCell(e.Note.Body)
		} else {
			table.AddCell(e.TargetType)
			table.AddCell(e.Title)
		}
		table.EndRow()
	}

	return table.Render()
}
