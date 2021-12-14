package cmdutils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/profclems/glab/pkg/iostreams"

	"github.com/profclems/glab/api"
	"github.com/profclems/glab/internal/glrepo"
	"github.com/profclems/glab/pkg/utils"
	"github.com/xanzy/go-gitlab"

	"github.com/AlecAivazis/survey/v2"
	"github.com/profclems/glab/pkg/prompt"
	"github.com/profclems/glab/pkg/surveyext"

	"github.com/profclems/glab/internal/config"

	"github.com/profclems/glab/pkg/git"
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
	if err != nil {
		return []string{}, nil
	}
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
		if strings.HasPrefix(file, ".") || !strings.HasSuffix(file, ".md") {
			continue
		}
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

func EditorPrompt(response *string, question, templateContent, editorCommand string) error {
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
			Name: question,
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
	lOpts := &gitlab.ListLabelsOptions{}
	lOpts.PerPage = 100
	labels, err := api.ListLabels(apiClient, repoRemote.FullName(), lOpts)
	if err != nil {
		return err
	}

	if len(labels) != 0 {
		var labelOptions []string

		for i := range labels {
			labelOptions = append(labelOptions, labels[i].Name)
		}

		var selectedLabels []string
		err = prompt.MultiSelect(&selectedLabels, "labels", "Select Labels", labelOptions)
		if err != nil {
			return err
		}
		*response = append(*response, selectedLabels...)
		return nil
	}

	var responseString string
	err = prompt.AskQuestionWithInput(&responseString, "labels", "Label(s) [Comma Separated]", "", false)
	if err != nil {
		return err
	}
	if responseString != "" {
		*response = append(*response, strings.Split(responseString, ",")...)
	}
	return nil
}

func MilestonesPrompt(response *int, apiClient *gitlab.Client, repoRemote *glrepo.Remote, io *iostreams.IOStreams) (err error) {
	var milestoneOptions []string
	milestoneMap := map[string]int{}

	lOpts := &api.ListMilestonesOptions{
		IncludeParentMilestones: gitlab.Bool(true),
		State:                   gitlab.String("active"),
		PerPage:                 100,
	}
	milestones, err := api.ListAllMilestones(apiClient, repoRemote.FullName(), lOpts)
	if err != nil {
		return err
	}
	if len(milestones) == 0 {
		fmt.Fprintln(io.StdErr, "There are no active milestones in this project")
		return nil
	}

	for i := range milestones {
		milestoneOptions = append(milestoneOptions, milestones[i].Title)
		milestoneMap[milestones[i].Title] = milestones[i].ID
	}

	var selectedMilestone string
	err = prompt.Select(&selectedMilestone, "milestone", "Select Milestone", milestoneOptions)
	if err != nil {
		return err
	}
	*response = milestoneMap[selectedMilestone]

	return nil
}

// GroupMemberLevel maps a number representing the access level to a string shown to the
// user.
// API docs:
// https://docs.gitlab.com/ce/api/members.html#valid-access-levels
var GroupMemberLevel = map[int]string{
	0:  "no access",
	5:  "minimal access",
	10: "guest",
	20: "reporter",
	30: "developer",
	40: "maintainer",
	50: "owner",
}

// AssigneesPrompt creates a multi-selection prompt of all the users below the given access level
// for the remote referenced by the `*glrepo.Remote`
func AssigneesPrompt(response *[]string, apiClient *gitlab.Client, repoRemote *glrepo.Remote, io *iostreams.IOStreams, minimumAccessLevel int) (err error) {
	var assigneeOptions []string
	assigneeMap := map[string]string{}

	lOpts := &gitlab.ListProjectMembersOptions{}
	lOpts.PerPage = 100
	members, err := api.ListProjectMembers(apiClient, repoRemote.FullName(), lOpts)
	if err != nil {
		return err
	}

	for i := range members {
		if members[i].AccessLevel >= gitlab.AccessLevelValue(minimumAccessLevel) {
			assigneeOptions = append(assigneeOptions, fmt.Sprintf("%s (%s)",
				members[i].Username,
				GroupMemberLevel[int(members[i].AccessLevel)],
			))
			assigneeMap[fmt.Sprintf("%s (%s)", members[i].Username, GroupMemberLevel[int(members[i].AccessLevel)])] = members[i].Username
		}
	}
	if len(assigneeOptions) == 0 {
		fmt.Fprintf(io.StdErr, "Couldn't fetch any members with minimum permission level %d\n", minimumAccessLevel)
		return nil
	}

	var selectedAssignees []string
	err = prompt.MultiSelect(&selectedAssignees, "assignees", "Select assignees", assigneeOptions)
	if err != nil {
		return err
	}
	for _, x := range selectedAssignees {
		*response = append(*response, assigneeMap[x])
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
	EditCommitMessageAction
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

	var confirmAnswer string
	err := prompt.Select(&confirmAnswer, "confirmation", "What's next?", options)
	if err != nil {
		return -1, fmt.Errorf("could not prompt: %w", err)
	}

	switch confirmAnswer {
	case submitLabel:
		return SubmitAction, nil
	case previewLabel:
		return PreviewAction, nil
	case addMetadataLabel:
		return AddMetadataAction, nil
	case cancelLabel:
		return CancelAction, nil
	default:
		return -1, fmt.Errorf("invalid value: %s", confirmAnswer)
	}
}

const (
	AddLabelAction Action = iota
	AddAssigneeAction
	AddMilestoneAction
)

func PickMetadata() ([]Action, error) {
	const (
		labelsLabel    = "labels"
		assigneeLabel  = "assignees"
		milestoneLabel = "milestones"
	)

	options := []string{
		labelsLabel,
		assigneeLabel,
		milestoneLabel,
	}

	var confirmAnswers []string
	err := prompt.MultiSelect(&confirmAnswers, "metadata", "Which metadata types to add?", options)
	if err != nil {
		return nil, fmt.Errorf("could not prompt: %w", err)
	}

	var pickedActions []Action

	for _, x := range confirmAnswers {
		switch x {
		case labelsLabel:
			pickedActions = append(pickedActions, AddLabelAction)
		case assigneeLabel:
			pickedActions = append(pickedActions, AddAssigneeAction)
		case milestoneLabel:
			pickedActions = append(pickedActions, AddMilestoneAction)
		}
	}
	return pickedActions, nil
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

	milestone, err := api.ProjectMilestoneByTitle(apiClient, repo.FullName(), milestoneTitle)
	if err != nil {
		return 0, err
	}

	return milestone.ID, nil
}

// UserAssignments holds 3 slice strings that represent which assignees should be added, removed, and replaced
// helper functions are also provided
type UserAssignments struct {
	ToAdd          []string
	ToRemove       []string
	ToReplace      []string
	AssignmentType UserAssignmentType
}

type UserAssignmentType int

const (
	AssigneeAssignment UserAssignmentType = iota
	ReviewerAssignment
)

// ParseAssignees takes a String Slice and splits them into 3 Slice Strings based on
// the first character of a string.
//
// '+' is put in the first slice, '!' and '-' in the second slice and all other cases
// in the third slice.
//
// The 3 String slices are returned regardless if anything was put it in or not the user
// is responsible for checking the length to see if anything is in it
func ParseAssignees(assignees []string) *UserAssignments {
	ua := UserAssignments{
		AssignmentType: AssigneeAssignment,
	}

	for _, assignee := range assignees {
		switch string([]rune(assignee)[0]) {
		case "+":
			ua.ToAdd = append(ua.ToAdd, string([]rune(assignee)[1:]))
		case "!", "-":
			ua.ToRemove = append(ua.ToRemove, string([]rune(assignee)[1:]))
		default:
			ua.ToReplace = append(ua.ToReplace, assignee)
		}
	}
	return &ua
}

// VerifyAssignees is a method for UserAssignments that checks them for validity
func (ua *UserAssignments) VerifyAssignees() error {
	// Fail if relative and absolute assignees were given, there is no reason to mix them.
	if len(ua.ToReplace) != 0 && (len(ua.ToAdd) != 0 || len(ua.ToRemove) != 0) {
		return errors.New("mixing relative (+,!,-) and absolute assignments is forbidden")
	}

	if m := utils.CommonElementsInStringSlice(ua.ToAdd, ua.ToRemove); len(m) != 0 {
		return fmt.Errorf("%s %q present in both add and remove which is forbidden",
			utils.Pluralize(len(m), "element"),
			strings.Join(m, " "))
	}
	return nil
}

// UsersFromReplaces converts all users from the `ToReplace` member of the struct into
// an Slice of String representing the Users' IDs, it also takes a Slice of Strings and
// writes a proper action message to it
func (ua *UserAssignments) UsersFromReplaces(apiClient *gitlab.Client, actions []string) ([]int, []string, error) {
	users, err := api.UsersByNames(apiClient, ua.ToReplace)
	if err != nil {
		return nil, actions, err
	}
	var usernames []string
	for i := range users {
		usernames = append(usernames, fmt.Sprintf("@%s", users[i].Username))
	}
	if len(usernames) != 0 {
		if ua.AssignmentType == ReviewerAssignment {
			actions = append(actions, fmt.Sprintf("requested review from %q", strings.Join(usernames, " ")))
		} else {
			actions = append(actions, fmt.Sprintf("assigned to %q", strings.Join(usernames, " ")))
		}
	}
	return IDsFromUsers(users), actions, nil
}

// UsersFromAddRemove works with both `ToAdd` and `ToRemove` members to produce a Slice of Ints that
// represents the final collection of IDs to assigned.
//
// It starts by getting all IDs already assigned, but ignoring ones present in `ToRemove`, it then
// converts all `usernames` in `ToAdd` into IDs by using the `api` package and adds them to the
// IDs to be assigned
func (ua *UserAssignments) UsersFromAddRemove(
	issueAssignees []*gitlab.IssueAssignee,
	mergeRequestAssignees []*gitlab.BasicUser,
	apiClient *gitlab.Client,
	actions []string,
) ([]int, []string, error) {

	var assignedIDs []int
	var usernames []string

	// Only one of those is required
	if mergeRequestAssignees != nil && issueAssignees != nil {
		return nil, actions, fmt.Errorf("issueAssignes and mergeRequestAssignes can't both not be nil")
	}

	// Path for Issues
	for i := range issueAssignees {
		// Only store them in assigneedIDs if they are not marked for removal
		if !utils.PresentInStringSlice(ua.ToRemove, issueAssignees[i].Username) {
			assignedIDs = append(assignedIDs, issueAssignees[i].ID)
		}
	}

	// Path for Merge Requests
	for i := range mergeRequestAssignees {
		// Only store them in assigneedIDs if they are not marked for removal
		if !utils.PresentInStringSlice(ua.ToRemove, mergeRequestAssignees[i].Username) {
			assignedIDs = append(assignedIDs, mergeRequestAssignees[i].ID)
		}
	}

	// Add action string
	if len(ua.ToRemove) != 0 {
		for _, x := range ua.ToRemove {
			usernames = append(usernames, fmt.Sprintf("@%s", x))
		}
		if ua.AssignmentType == ReviewerAssignment {
			actions = append(actions, fmt.Sprintf("removed review request for %q", strings.Join(usernames, " ")))
		} else {
			actions = append(actions, fmt.Sprintf("unassigned %q", strings.Join(usernames, " ")))
		}
	}

	if len(ua.ToAdd) != 0 {
		users, err := api.UsersByNames(apiClient, ua.ToAdd)
		if err != nil {
			return nil, nil, err
		}
		// Work-around GitLab (the company's own instance, not all instances have this) bug
		// which causes a 500 Internal Error if duplicate `IDs` are used. Filter out any
		// IDs that is already present
		for i := range users {
			if !utils.PresentInIntSlice(assignedIDs, users[i].ID) {
				assignedIDs = append(assignedIDs, users[i].ID)
			}
		}

		// Reset the usernames array because it might have been used by `unassignedUsers`
		usernames = []string{}

		for _, x := range ua.ToAdd {
			usernames = append(usernames, fmt.Sprintf("@%s", x))
		}
		if ua.AssignmentType == ReviewerAssignment {
			actions = append(actions, fmt.Sprintf("requested review from %q", strings.Join(usernames, " ")))
		} else {
			actions = append(actions, fmt.Sprintf("assigned %q", strings.Join(usernames, " ")))
		}
	}

	// That means that all assignees were removed but we can't pass an empty Slice of Ints so
	// pass the documented value of 0
	if len(assignedIDs) == 0 {
		assignedIDs = []int{0}
	}
	return assignedIDs, actions, nil
}
