package events

import (
	"fmt"
	"io"

	"github.com/profclems/glab/api"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/pkg/utils"
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

			repo, err := f.BaseRepo()
			if err != nil {
				return err
			}

			events, err := api.CurrentUserEvents(apiClient)
			if err != nil {
				return err
			}

			if err = f.IO.StartPager(); err != nil {
				return err
			}
			defer f.IO.StopPager()

			if lb, _ := cmd.Flags().GetBool("all"); lb {
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

				DisplayAllEvents(f.IO.StdOut, events, projects)
				return nil
			}

			project, err := api.GetProject(apiClient, repo.FullName())
			if err != nil {
				return err
			}

			DisplayProjectEvents(f.IO.StdOut, events, project)
			return nil
		},
	}

	cmd.Flags().BoolP("all", "a", false, "Get events from all projects")

	return cmd
}

func DisplayProjectEvents(w io.Writer, events []*gitlab.ContributionEvent, project *gitlab.Project) {
	for _, e := range events {
		if e.ProjectID != project.ID {
			continue
		}
		printEvent(w, e, project)
	}
}

func DisplayAllEvents(w io.Writer, events []*gitlab.ContributionEvent, projects map[int]*gitlab.Project) {
	for _, e := range events {
		printEvent(w, e, projects[e.ProjectID])
	}
}

func printEvent(w io.Writer, e *gitlab.ContributionEvent, project *gitlab.Project) {
	switch e.ActionName {
	case "pushed to":
		fmt.Fprintf(w, "Pushed to %s %s at %s\n%q\n", e.PushData.RefType, e.PushData.Ref, project.NameWithNamespace, e.PushData.CommitTitle)
	case "deleted":
		fmt.Fprintf(w, "Deleted %s %s at %s\n", e.PushData.RefType, e.PushData.Ref, project.NameWithNamespace)
	case "pushed new":
		fmt.Fprintf(w, "Pushed new %s %s at %s\n", e.PushData.RefType, e.PushData.Ref, project.NameWithNamespace)
	case "commented on":
		fmt.Fprintf(w, "Commented on %s #%s at %s\n%q\n", e.Note.NoteableType, e.Note.Title, project.NameWithNamespace, e.Note.Body)
	case "accepted":
		fmt.Fprintf(w, "Accepted %s %s at %s\n", e.TargetType, e.TargetTitle, project.NameWithNamespace)
	case "opened":
		fmt.Fprintf(w, "Opened %s %s at %s\n", e.TargetType, e.TargetTitle, project.NameWithNamespace)
	case "closed":
		fmt.Fprintf(w, "Closed %s %s at %s\n", e.TargetType, e.TargetTitle, project.NameWithNamespace)
	case "joined":
		fmt.Fprintf(w, "Joined %s\n", project.NameWithNamespace)
	case "left":
		fmt.Fprintf(w, "Left %s\n", project.NameWithNamespace)
	case "created":
		targetType := e.TargetType
		if e.TargetType == "WikiPage::Meta" {
			targetType = "Wiki page"
		}
		fmt.Fprintf(w, "Created %s %s at %s\n", targetType, e.TargetTitle, project.NameWithNamespace)
	default:
		fmt.Fprintf(w, "%s %q", e.TargetType, e.Title)
	}
	fmt.Fprintln(w) // to leave a blank line

}
