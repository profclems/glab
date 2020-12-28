package view

import (
	"fmt"
	"log"
	"runtime/debug"

	"github.com/AlecAivazis/survey/v2"
	"github.com/gdamore/tcell"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/glrepo"
	"github.com/profclems/glab/pkg/api"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

var (
	apiClient     *gitlab.Client
	project       *gitlab.Project
	repo          glrepo.Interface
	selectedBoard string
)

func NewCmdView(f *cmdutils.Factory) *cobra.Command {
	var viewCmd = &cobra.Command{
		Use:   "view [flags]",
		Short: `View project issue board.`,
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			a := tview.NewApplication()
			defer recoverPanic(a)

			//out := utils.ColorableOut(cmd)

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

			issueBoards, err := api.ListIssueBoards(apiClient, repo.FullName(), &gitlab.ListIssueBoardsOptions{})
			if err != nil {
				return fmt.Errorf("error retrieving issue board: %w", err)
			}
			boardStr := make([]string, len(issueBoards))
			boardInfo := map[string]int{}
			for i, board := range issueBoards {
				boardStr[i] = board.Name
				boardInfo[boardStr[i]] = board.ID
			}

			prompt := &survey.Select{
				Message: "Select Board:",
				Options: boardStr,
			}
			err = survey.AskOne(prompt, &selectedBoard)
			if err != nil {
				return err
			}

			boadLists, err := api.GetIssueBoardLists(apiClient, repo.FullName(), boardInfo[selectedBoard], &gitlab.GetIssueBoardListsOptions{})
			if err != nil {
				return err
			}

			root := tview.NewFlex()
			/*
				AddItem(tview.NewBox().SetBorder(true).SetTitle("Left (1/2 x width of Top)"), 0, 1, false).
				AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
					AddItem(tview.NewBox().SetBorder(true).SetTitle("Top"), 0, 1, false).
					AddItem(tview.NewBox().SetBorder(true).SetTitle("Middle (3 x height of Top)"), 0, 3, false).
					AddItem(tview.NewBox().SetBorder(true).SetTitle("Bottom (5 rows)"), 5, 1, false), 0, 2, false).
				AddItem(tview.NewBox().SetBorder(true).SetTitle("Right (20 cols)"), 20, 1, false)
			*/
			var issues []*gitlab.Issue
			// TODO: add `open` and `closed` board list. Both are not returned in the List API response payload
			for _, list := range boadLists {
				var boardIssues string
				issues, err = api.ListIssues(apiClient, repo.FullName(), &gitlab.ListProjectIssuesOptions{
					Labels: gitlab.Labels{list.Label.Name},
				})
				if err != nil {
					return fmt.Errorf("error retrieving list issues: %w", err)
				}
				bx := tview.NewTextView().SetDynamicColors(true)
				for _, issue := range issues {
					//label, _ := issue.Labels.MarshalJSON()
					//labelPrint := strings.Split(string(label), ", ")
					var labelPrint string
					var assignee string
					totalLables := len(issue.Labels)
					if totalLables > 0 {
						for i, l := range issue.Labels {
							if i == (totalLables - 1) {
								labelPrint += l
							} else {
								labelPrint += l + ", "
							}
						}
					}
					if issue.Assignee != nil {
						assignee = issue.Assignee.Username
					}
					boardIssues += fmt.Sprintf("[white]%s\n[blue](%s)\n[green]#%d[white] - %s\n\n", issue.Title, labelPrint, issue.IID, assignee)
				}
				bx.SetText(boardIssues).SetWrap(true)
				bx.SetBorder(true).SetTitle(list.Label.Name).SetTitleColor(tcell.GetColor(list.Label.Color))
				root.AddItem(bx, 0, 1, false)
			}

			root.SetBorderPadding(1, 1, 2, 2).SetBorder(true).SetTitle(fmt.Sprintf(" %s • Boards • %s ", selectedBoard, project.NameWithNamespace))
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

	return viewCmd
}

func recoverPanic(app *tview.Application) {
	if r := recover(); r != nil {
		app.Stop()
		log.Fatalf("%s\n%s\n", r, string(debug.Stack()))
	}
}

func drawBIssueoards(app *tview.Application) error {
	//flex := tview.NewFlex().
	return nil
}
