package cmdutils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/profclems/glab/internal/glrepo"
	"github.com/profclems/glab/internal/utils"
	"github.com/profclems/glab/pkg/api"
	"github.com/xanzy/go-gitlab"

	"github.com/AlecAivazis/survey/v2"
	"github.com/profclems/glab/pkg/prompt"
	"github.com/profclems/glab/pkg/surveyext"

	"github.com/profclems/glab/internal/config"

	"github.com/profclems/glab/internal/git"
)

const (
	IssueTemplate        = "issue_templates"
	MergeRequestTemplate = "merge_request_templates"
)

// LoadGitLabTemplate finds and loads the GitLab template from the working git directory
// Follows the format officially supported by GitLab
// https://docs.gitlab.com/ee/user/project/description_templates.html#setting-a-default-template-for-issues-and-merge-requests.
//
// TODO: load from remote repository if repo is overriden by -R flag
func LoadGitLabTemplate(tmplType, tmplName string) (string, error) {
	wdir, err := git.ToplevelDir()
	if err != nil {
		return "", err
	}

	if !strings.HasSuffix(tmplName, ".md") {
		tmplName = tmplName + ".md"
	}

	tmplFile := filepath.Join(wdir, ".gitlab", tmplType, tmplName)
	f, err := os.Open(tmplFile)
	if os.IsNotExist(err) {
		return "", nil
	} else if err != nil {
		return "", err
	}

	tmpl, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(tmpl)), nil
}

// TODO: properly handle errors in this function.
//       For now, it returns nil and empty slice if there's an error
func ListGitLabTemplates(tmplType string) ([]string, error) {
	wdir, err := git.ToplevelDir()
	tmplFolder := filepath.Join(wdir, ".gitlab", tmplType)
	var files []string
	f, err := os.Open(tmplFolder)
	// if error return an empty slice since it only returns PathError
	if err != nil {
		return files, nil
	}
	fileNames, err := f.Readdirnames(-1)
	defer f.Close()
	if err != nil {
		// return empty slice if error
		return files, nil
	}

	for _, file := range fileNames {
		files = append(files, strings.TrimSuffix(file, ".md"))
	}
	return files, nil
}

func GetEditor(cf func() (config.Config, error)) (string, error) {
	cfg, err := cf()
	if err != nil {
		return "", fmt.Errorf("could not read config: %w", err)
	}
	// will search in the order glab_editor, visual, editor first from the env before the config file
	editorCommand, _ := cfg.Get("", "editor")

	return editorCommand, nil
}

func DescriptionPrompt(response *string, templateContent, editorCommand string) error {
	defaultBody := *response
	if templateContent != "" {
		if defaultBody != "" {
			// prevent excessive newlines between default body and template
			defaultBody = strings.TrimRight(defaultBody, "\n")
			defaultBody += "\n\n"
		}
		defaultBody += templateContent
	}

	qs := []*survey.Question{
		{
			Name: "Description",
			Prompt: &surveyext.GLabEditor{
				BlankAllowed:  true,
				EditorCommand: editorCommand,
				Editor: &survey.Editor{
					Message:       "Description",
					FileName:      "*.md",
					Default:       defaultBody,
					HideDefault:   true,
					AppendDefault: true,
				},
			},
		},
	}

	err := prompt.Ask(qs, response)
	if err != nil {
		return err
	}
	if *response == "" {
		*response = defaultBody
	}
	return nil
}

func LabelsPrompt(response *[]string, apiClient *gitlab.Client, repoRemote *glrepo.Remote) (err error) {
	var addLabels bool
	err = prompt.Confirm(&addLabels, "Do you want to add labels?", true)
	if err != nil {
		return
	}
	if addLabels {
		var labelOptions string
		if repoRemote.Name != "" {
			labelOptions, _ = git.Config("remote." + repoRemote.Name + ".glab-cached-labels")
		}
		if labelOptions == "" {
			lOpts := &gitlab.ListLabelsOptions{}
			lOpts.PerPage = 100
			labels, err := api.ListLabels(apiClient, repoRemote.FullName(), lOpts)
			if err == nil && labels != nil {
				for i, label := range labels {
					if i > 0 {
						labelOptions += ","
					}
					labelOptions += label.Name
				}
				if labelOptions != "" && repoRemote.Name != "" {
					// silently fails if not a git repo
					_ = git.SetConfig(repoRemote.Name, "glab-cached-labels", labelOptions)
				}
			}
		}
		if labelOptions != "" {
			var selectedLabels []string
			err = prompt.MultiSelect(&selectedLabels, "Select Labels", strings.Split(labelOptions, ","))
			if err != nil {
				return err
			}
			*response = selectedLabels

		} else {
			var responseString string
			err = prompt.AskQuestionWithInput(&responseString, "Label(s) [Comma Separated]", "", false)
			if err != nil {
				return err
			}
			*response = strings.Split(responseString, ",")
		}
	}
	return nil
}

func MilestonesPrompt(response *int, apiClient *gitlab.Client, repoRemote *glrepo.Remote, io *utils.IOStreams) (err error) {
	var addMilestones bool
	err = prompt.Confirm(&addMilestones, "Do you wish to attach a milestone?", true)
	if err != nil {
		return err
	}
	if addMilestones {
		var milestoneOptions []string
		milestoneMap := map[string]*gitlab.Milestone{}

		lOpts := &gitlab.ListMilestonesOptions{
			State: gitlab.String("active"),
		}
		lOpts.PerPage = 100
		milestones, err := api.ListMilestones(apiClient, repoRemote.FullName(), lOpts)
		if err != nil {
			return err
		}
		if len(milestones) == 0 {
			fmt.Fprintln(io.StdErr, "There are no active milestones in this project")
			return nil
		}

		for i := range milestones {
			milestoneOptions = append(milestoneOptions, milestones[i].Title)
			milestoneMap[milestones[i].Title] = milestones[i]
		}

		var selectedMilestone string
		err = prompt.AskOne(&survey.Select{
			Message: "Select Milestone",
			Options: milestoneOptions,
		}, &selectedMilestone)
		if err != nil {
			return err
		}
		*response = milestoneMap[selectedMilestone].ID
	}
	return nil
}

func AssigneesPrompt(response *[]string) (err error) {
	var addAssignees bool
	err = prompt.Confirm(&addAssignees, "Do you wish to assign users?", true)
	if err != nil {
		return err
	}
	if addAssignees {
		var responseString string
		err = prompt.AskQuestionWithInput(&responseString, "Username(s) [Comma Separated]", "", false)
		if err != nil {
			return err
		}
		*response = strings.Split(responseString, ",")
	}
	return nil
}

type Action int

const (
	NoAction Action = iota
	SubmitAction
	PreviewAction
	AddMetadataAction
	CancelAction
)

func ConfirmSubmission(allowPreview bool, allowAddMetadata bool) (Action, error) {
	const (
		submitLabel      = "Submit"
		previewLabel     = "Continue in browser"
		addMetadataLabel = "Add metadata"
		cancelLabel      = "Cancel"
	)

	options := []string{submitLabel}
	if allowPreview {
		options = append(options, previewLabel)
	}
	if allowAddMetadata {
		options = append(options, addMetadataLabel)
	}
	options = append(options, cancelLabel)

	confirmAnswers := struct {
		Confirmation int
	}{}
	confirmQs := []*survey.Question{
		{
			Name: "confirmation",
			Prompt: &survey.Select{
				Message: "What's next?",
				Options: options,
			},
		},
	}

	err := prompt.Ask(confirmQs, &confirmAnswers)
	if err != nil {
		return -1, fmt.Errorf("could not prompt: %w", err)
	}

	switch options[confirmAnswers.Confirmation] {
	case submitLabel:
		return SubmitAction, nil
	case previewLabel:
		return PreviewAction, nil
	case addMetadataLabel:
		return AddMetadataAction, nil
	case cancelLabel:
		return CancelAction, nil
	default:
		return -1, fmt.Errorf("invalid index: %d", confirmAnswers.Confirmation)
	}
}

//IDsFromUsers collects all user IDs from a slice of users
func IDsFromUsers(users []*gitlab.User) []int {
	ids := make([]int, len(users))
	for i, user := range users {
		ids[i] = user.ID
	}
	return ids
}

func ParseMilestone(apiClient *gitlab.Client, repo glrepo.Interface, milestoneTitle string) (int, error) {
	if milestoneID, err := strconv.Atoi(milestoneTitle); err == nil {
		return milestoneID, nil
	}

	milestone, err := api.MilestoneByTitle(apiClient, repo.FullName(), milestoneTitle)
	if err != nil {
		return 0, err
	}

	return milestone.ID, nil
}
