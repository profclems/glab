package view

import (
	"fmt"
	"log"
	"runtime/debug"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/gdamore/tcell/v2"
	"github.com/profclems/glab/api"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/glrepo"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

var (
	apiClient *gitlab.Client
	project   *gitlab.Project
	repo      glrepo.Interface
)

type issueBoardViewOptions struct {
	Assignee  string
	Labels    string
	Milestone string
	State     string
}

func NewCmdView(f *cmdutils.Factory) *cobra.Command {
	var opts = &issueBoardViewOptions{}
	var viewCmd = &cobra.Command{
		Use:   "view [flags]",
		Short: `View project issue board.`,
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			a := tview.NewApplication()
			defer recoverPanic(a)

			apiClient, err = f.HttpClient()
			if err != nil {
				return err
			}

			repo, err = f.BaseRepo()
			if err != nil {
				return err
			}

			project, err = api.GetProject(apiClient, repo.FullName())
			if err != nil {
				return fmt.Errorf("failed to get project: %w", err)
			}

			// list the groups that are ancestors to project:
			// https://docs.gitlab.com/ee/api/projects.html#list-a-projects-groups
			projectGroups, err := api.ListProjectGroups(apiClient, project.ID, &gitlab.ListProjectGroupOptions{})
			if err != nil {
				return err
			}

			// get group issue boards for each ancestor group returned:
			// https://docs.gitlab.com/ee/api/group_boards.html#list-all-group-issue-boards-in-a-group
			projectGroupIssueBoards := []*gitlab.GroupIssueBoard{}
			for _, projectGroup := range projectGroups {
				groupIssueBoards, err := api.ListGroupIssueBoards(apiClient, projectGroup.ID, &gitlab.ListGroupIssueBoardsOptions{})
				if err != nil {
					return err
				}
				projectGroupIssueBoards = append(groupIssueBoards, projectGroupIssueBoards...)
			}

			projectIssueBoards, err := api.ListProjectIssueBoards(apiClient, repo.FullName(), &gitlab.ListIssueBoardsOptions{})
			if err != nil {
				return fmt.Errorf("error retrieving issue board: %w", err)
			}

			type info struct {
				id    int
				group *gitlab.Group
			}

			boardSelectionStr := []string{}
			boardInfo := map[string]info{}
			for _, board := range projectGroupIssueBoards {
				boardSelectionStr = append(boardSelectionStr, fmt.Sprintf("%s%*s", board.Name, 50-len(board.Name), "(Group)"))
				boardInfo[board.Name] = info{id: board.ID, group: board.Group}
			}
			for _, board := range projectIssueBoards {
				boardSelectionStr = append(boardSelectionStr, fmt.Sprintf("%s%*s", board.Name, 50-len(board.Name), "(Project)"))
				boardInfo[board.Name] = info{id: board.ID}
			}

			var selectedOption string
			prompt := &survey.Select{
				Message: "Select Board:",
				Options: boardSelectionStr,
			}
			err = survey.AskOne(prompt, &selectedOption)
			if err != nil {
				return err
			}
			selectedBoard := strings.Split(selectedOption, " ")[0]

			var boardLists []*gitlab.BoardList
			if boardInfo[selectedBoard].group != nil {
				boardLists, err = api.GetGroupIssueBoardLists(apiClient, boardInfo[selectedBoard].group.ID,
					boardInfo[selectedBoard].id, &gitlab.ListGroupIssueBoardListsOptions{})
				if err != nil {
					return err
				}
			}

			if boardInfo[selectedBoard].group == nil {
				boardLists, err = api.GetPojectIssueBoardLists(apiClient, repo.FullName(),
					boardInfo[selectedBoard].id, &gitlab.GetIssueBoardListsOptions{})
				if err != nil {
					return err
				}
			}

			// manually add 'opened' and 'closed' lists before and after fetched lists
			opened := &gitlab.BoardList{
				Label: &gitlab.Label{
					Name:      "Open",
					Color:     "#fabd2f",
					TextColor: "#000000",
				},
				Position: 0,
			}
			boardLists = append([]*gitlab.BoardList{opened}, boardLists...)

			closed := &gitlab.BoardList{
				Label: &gitlab.Label{
					Name:      "Closed",
					Color:     "#8ec07c",
					TextColor: "#000000",
				},
				Position: len(boardLists),
			}
			boardLists = append(boardLists, closed)

			root := tview.NewFlex()
			issues := []*gitlab.Issue{}
			for _, list := range boardLists {
				if list.Label == nil {
					continue
				}

				var boardIssues, listTitle, listColor string

				if list.Label == nil {
					listTitle = "Unnamed"
					listColor = "#FFA500" // orange
				}

				if list.Label != nil {
					listTitle = list.Label.Name
					listColor = list.Label.Color
				}

				// automatically request using state for default "open" and "closed" lists
				// this is required because these lists aren't returned with the board lists api call
				if list.Label.Name == "Closed" {
					opts.State = "closed"
				}

				if list.Label.Name == "Open" {
					opts.State = "opened"
				}

				// "closed" and "open" are considered "state lists" since they
				// need to be requested using state and not label
				isStateList := list.Label.Name == "Open" || list.Label.Name == "Closed"

				if boardInfo[selectedBoard].group != nil {
					reqOpts, err := parseListGroupIssueOptions(opts)
					if err != nil {
						return err
					}
					issues, err = api.ListGroupIssues(apiClient, boardInfo[selectedBoard].group.ID, reqOpts)
					if err != nil {
						return err
					}
				}

				if boardInfo[selectedBoard].group == nil {
					reqOpts, err := parseListProjectIssueOptions(opts)
					if err != nil {
						return err
					}
					issues, err = api.ListProjectIssues(apiClient, repo.FullName(), reqOpts)
					if err != nil {
						return err
					}
				}

				if err != nil {
					return fmt.Errorf("error retrieving list issues: %w", err)
				}

				bx := tview.NewTextView().SetDynamicColors(true)

			next:
				for _, issue := range issues {
					// skip issues labeled for other board lists when populating the "open" list
					if opts.State == "opened" && isStateList {
						for _, boardList := range boardLists {
							for _, issueLabel := range issue.Labels {
								if issueLabel == boardList.Label.Name {
									continue next
								}
							}
						}
					}

					// skip all issues without the "closed" state for the "closed" list
					if opts.State == "closed" && isStateList {
						if issue.State != "closed" {
							continue next
						}
					}

					// filter labeled issues into matching label board lists
					if !isStateList {
						var hasListLabel bool
						for _, issueLabel := range issue.Labels {
							if issueLabel == list.Label.Name {
								hasListLabel = true
							}
						}

						if !hasListLabel {
							continue next
						}
					}

					var assignee, labelString string
					if len(issue.Labels) > 0 {
						labelString = buildLabelString(issue.LabelDetails)
					}
					if issue.Assignee != nil {
						assignee = issue.Assignee.Username
					}

					boardIssues += fmt.Sprintf("[white::b]%s\n%s[green:-:-]#%d[darkgray] - %s\n\n",
						issue.Title, labelString, issue.IID, assignee)
				}
				bx.SetText(boardIssues).SetWrap(true)
				bx.SetBorder(true).SetTitle(listTitle).SetTitleColor(tcell.GetColor(listColor))
				root.AddItem(bx, 0, 1, false)
			}

			root.SetBorderPadding(1, 1, 2, 2).SetBorder(true).SetTitle(
				fmt.Sprintf(" %s • Boards • %s ", selectedBoard, project.NameWithNamespace))
			screen, err := tcell.NewScreen()
			if err != nil {
				return err
			}
			_ = screen.Init()
			if err := a.SetScreen(screen).SetRoot(root, true).Run(); err != nil {
				return err
			}
			return nil
		},
	}

	viewCmd.Flags().StringVarP(&opts.Assignee, "assignee", "a", "", "Filter board issues by assignee username")
	viewCmd.Flags().StringVarP(&opts.Labels, "labels", "l", "", "Filter board issues by labels (comma separated)")
	viewCmd.Flags().StringVarP(&opts.Milestone, "milestone", "m", "", "Filter board issues by milestone")
	return viewCmd
}

func parseListProjectIssueOptions(opts *issueBoardViewOptions) (*gitlab.ListProjectIssuesOptions, error) {
	withLabelDetails := true
	reqOpts := &gitlab.ListProjectIssuesOptions{
		WithLabelDetails: &withLabelDetails,
	}

	if opts.Assignee != "" {
		reqOpts.AssigneeUsername = &opts.Assignee
	}

	if opts.Labels != "" {
		reqOpts.Labels = gitlab.Labels(strings.Split(opts.Labels, ","))
	}

	if opts.State != "" {
		reqOpts.State = &opts.State
	}

	if opts.Milestone != "" {
		reqOpts.Milestone = &opts.Milestone
	}
	return reqOpts, nil
}

func parseListGroupIssueOptions(opts *issueBoardViewOptions) (*gitlab.ListGroupIssuesOptions, error) {
	withLabelDetails := true
	reqOpts := &gitlab.ListGroupIssuesOptions{
		WithLabelDetails: &withLabelDetails,
	}

	if opts.Assignee != "" {
		reqOpts.AssigneeUsername = &opts.Assignee
	}

	if opts.Labels != "" {
		reqOpts.Labels = gitlab.Labels(strings.Split(opts.Labels, ","))
	}

	if opts.State != "" {
		reqOpts.State = &opts.State
	}

	if opts.Milestone != "" {
		reqOpts.Milestone = &opts.Milestone
	}
	return reqOpts, nil
}

func recoverPanic(app *tview.Application) {
	if r := recover(); r != nil {
		app.Stop()
		log.Fatalf("%s\n%s\n", r, string(debug.Stack()))
	}
}

func buildLabelString(labelDetails []*gitlab.LabelDetails) string {
	var labels string
	for _, ld := range labelDetails {
		labels += fmt.Sprintf("[white:%s:-]%s[white:-:-] ", ld.Color, ld.Name)
	}
	labels += fmt.Sprintf("\n")
	return labels
}
