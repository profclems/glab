package cmdutils

import (
	"errors"
	"fmt"
	"testing"

	"github.com/profclems/glab/internal/git"
	"github.com/profclems/glab/internal/glrepo"
	"github.com/profclems/glab/internal/utils"
	"github.com/profclems/glab/pkg/api"
	"github.com/profclems/glab/pkg/prompt"
	"github.com/stretchr/testify/assert"
	"github.com/xanzy/go-gitlab"
)

func Test_ParseAssignees(t *testing.T) {
	testCases := []struct {
		name        string
		input       []string
		wantAdd     []string
		wantRemove  []string
		wantReplace []string
	}{
		{
			name:        "simple replace",
			input:       []string{"foo"},
			wantAdd:     []string{},
			wantRemove:  []string{},
			wantReplace: []string{"foo"},
		},
		{
			name:        "only add",
			input:       []string{"+foo"},
			wantAdd:     []string{"foo"},
			wantRemove:  []string{},
			wantReplace: []string{},
		},
		{
			name:        "only remove",
			input:       []string{"-foo", "!bar"},
			wantAdd:     []string{},
			wantRemove:  []string{"foo", "bar"},
			wantReplace: []string{},
		},
		{
			name:        "only replace",
			input:       []string{"baz"},
			wantAdd:     []string{},
			wantRemove:  []string{},
			wantReplace: []string{"baz"},
		},
		{
			name:        "add and remove",
			input:       []string{"+qux", "-foo", "!bar"},
			wantAdd:     []string{"qux"},
			wantRemove:  []string{"foo", "bar"},
			wantReplace: []string{},
		},
		{
			name:        "add and replace",
			input:       []string{"+foo", "bar"},
			wantAdd:     []string{"foo"},
			wantRemove:  []string{},
			wantReplace: []string{"bar"},
		},
		{
			name:        "remove and replace",
			input:       []string{"-foo", "bar", "!baz"},
			wantAdd:     []string{},
			wantRemove:  []string{"foo", "baz"},
			wantReplace: []string{"bar"},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			uaGot := ParseAssignees(tC.input)
			assert.ElementsMatch(t, uaGot.ToAdd, tC.wantAdd)
			assert.ElementsMatch(t, uaGot.ToRemove, tC.wantRemove)
			assert.ElementsMatch(t, uaGot.ToReplace, tC.wantReplace)
		})
	}

}

func Test_VerifyAssignees(t *testing.T) {
	testCases := []struct {
		name  string
		input UserAssignments
		want  string // expected error message
	}{
		{
			name: "empty, no errors",
			input: UserAssignments{
				ToAdd:     []string{},
				ToRemove:  []string{},
				ToReplace: []string{},
			},
		},
		{
			name: "simple addition, no errors",
			input: UserAssignments{
				ToAdd:     []string{"foo"},
				ToRemove:  []string{},
				ToReplace: []string{},
			},
		},
		{
			name: "simple removal, no errors",
			input: UserAssignments{
				ToAdd:     []string{},
				ToRemove:  []string{"foo"},
				ToReplace: []string{},
			},
		},
		{
			name: "simple replace, no errors",
			input: UserAssignments{
				ToAdd:     []string{},
				ToRemove:  []string{},
				ToReplace: []string{"foo"},
			},
		},
		{
			name: "add and removal with multiple elements, no errors",
			input: UserAssignments{
				ToAdd:     []string{"foo", "bar", "baz"},
				ToRemove:  []string{"qux", "quux", "quz"},
				ToReplace: []string{},
			},
		},
		{
			name: "multi replace, no errors",
			input: UserAssignments{
				ToAdd:     []string{},
				ToRemove:  []string{},
				ToReplace: []string{"foo", "bar"},
			},
		},
		{
			name: "replace with add, error",
			input: UserAssignments{
				ToAdd:     []string{"bar"},
				ToRemove:  []string{},
				ToReplace: []string{"foo"},
			},
			want: "mixing relative (+,!,-) and absolute assignments is forbidden",
		},
		{
			name: "replace with remove, error",
			input: UserAssignments{
				ToAdd:     []string{},
				ToRemove:  []string{"baz"},
				ToReplace: []string{"foo"},
			},
			want: "mixing relative (+,!,-) and absolute assignments is forbidden",
		},
		{
			name: "overlapping add and removal element, error",
			input: UserAssignments{
				ToAdd:     []string{"foo"},
				ToRemove:  []string{"foo"},
				ToReplace: []string{},
			},
			want: "1 element \"foo\" present in both add and remove which is forbidden",
		},
		{
			name: "overlapping add and removal elements, error",
			input: UserAssignments{
				ToAdd:     []string{"foo", "bar", "baz"},
				ToRemove:  []string{"foo", "baz"},
				ToReplace: []string{},
			},
			want: "2 elements \"foo baz\" present in both add and remove which is forbidden",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			err := tC.input.VerifyAssignees()
			if tC.want == "" {
				if err != nil {
					t.Errorf("VerifyAssignees() unexpected error = %s", err)
				}
			} else {
				if tC.want != err.Error() {
					t.Errorf("VerifyAssignees() expected = %s, got = %s", tC.want, err.Error())
				}
			}
		})
	}
}

func Test_UsersFromReplaces(t *testing.T) {
	testCases := []struct {
		name           string
		users          []*gitlab.User
		expectedIDs    []int
		expectedAction []string
	}{
		{
			name:           "nothingness",
			users:          []*gitlab.User{},
			expectedIDs:    []int{},
			expectedAction: []string{},
		},
		{
			name: "single user named foo",
			users: []*gitlab.User{
				{ID: 1, Username: "foo"},
			},
			expectedIDs:    []int{1},
			expectedAction: []string{"assigned to \"@foo\""},
		},
		{
			name: "multiple users named foo, bar and baz",
			users: []*gitlab.User{
				{ID: 1, Username: "foo"},
				{ID: 3, Username: "bar"},
				{ID: 7, Username: "baz"},
			},
			expectedIDs:    []int{1, 3, 7},
			expectedAction: []string{"assigned to \"@foo @bar @baz\""},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			ua := UserAssignments{}
			api.UsersByNames = func(apiClient *gitlab.Client, names []string) ([]*gitlab.User, error) {
				return tC.users, nil
			}
			var gotAction []string
			gotIDs, gotAction, err := ua.UsersFromReplaces(&gitlab.Client{}, gotAction)
			if err != nil {
				t.Errorf("UsersFromReplaces() unexpected error = %s", err)
			}
			assert.ElementsMatch(t, gotIDs, tC.expectedIDs)
			assert.ElementsMatch(t, gotAction, tC.expectedAction)
		})
	}
}

func Test_UserAssignmentsAPIFailure(t *testing.T) {
	want := "failed to get users by their names" // Error message we want
	ua := UserAssignments{
		ToAdd: []string{"foo"},
	} // Fill `ToAdd` so `cmdutils.UsersFromAddRemove()` reaches the api call
	var err error

	api.UsersByNames = func(apiClient *gitlab.Client, names []string) ([]*gitlab.User, error) {
		return nil, fmt.Errorf("failed to get users by their names")
	}

	apiClient := gitlab.Client{} // Empty Client, it won't be used, just to satisfy the function signature
	_, _, err = ua.UsersFromReplaces(&apiClient, nil)
	if err == nil {
		t.Errorf("UsersFromReplaces() expected error to not be nil")
	}
	if want != err.Error() {
		t.Errorf("UsersFromReplace() expected error = %s, got = %w", want, err)
	}

	_, _, err = ua.UsersFromAddRemove(nil, nil, &apiClient, nil)
	if err == nil {
		t.Errorf("UsersFromReplaces() expected error to not be nil")
	}
	if want != err.Error() {
		t.Errorf("UsersFromReplace() expected error = %s, got = %w", want, err)
	}
}

func Test_UsersFromAddRemove(t *testing.T) {
	testCases := []struct {
		name           string
		users          []*gitlab.User          // Mock *gitlab.User received from api.UsersByNames
		merge          []*gitlab.BasicUser     // Mock `.Assignee field` from a merge request
		issue          []*gitlab.IssueAssignee // Mock `.Assignee field` from an issue
		expectedIDs    []int
		expectedAction []string
		ua             UserAssignments
		wantErr        string
	}{
		{
			name: "add foo (issue and merge request)",
			users: []*gitlab.User{
				{
					ID:       1,
					Username: "foo",
				},
			},
			expectedIDs:    []int{1},
			expectedAction: []string{"assigned \"@foo\""},
			ua:             UserAssignments{ToAdd: []string{"foo"}},
		},
		{
			name: "add foo, bar and baz (issue and merge request)",
			users: []*gitlab.User{
				{
					ID:       1,
					Username: "foo",
				},
				{
					ID:       235,
					Username: "bar",
				},
				{
					ID:       1500,
					Username: "baz",
				},
			},
			expectedIDs:    []int{1, 235, 1500},
			expectedAction: []string{"assigned \"@foo @bar @baz\""},
			ua:             UserAssignments{ToAdd: []string{"foo", "bar", "baz"}},
		},
		{
			name:  "remove foo (issue)",
			users: []*gitlab.User{},
			issue: []*gitlab.IssueAssignee{
				{
					ID:       1,
					Username: "foo",
				},
			},
			expectedIDs:    []int{0},
			expectedAction: []string{"unassigned \"@foo\""},
			ua:             UserAssignments{ToRemove: []string{"foo"}},
		},
		{
			name:  "remove foo and baz out of foo, bar and baz (issue)",
			users: []*gitlab.User{},
			issue: []*gitlab.IssueAssignee{
				{
					ID:       1,
					Username: "foo",
				},
				{
					ID:       2,
					Username: "bar",
				},
				{
					ID:       3,
					Username: "baz",
				},
			},
			expectedIDs:    []int{2},
			expectedAction: []string{"unassigned \"@foo @baz\""},
			ua:             UserAssignments{ToRemove: []string{"foo", "baz"}},
		},
		{
			name: "remove foo out of foo and baz and add bar (issue)",
			users: []*gitlab.User{
				{
					ID:       100,
					Username: "bar",
				},
			},
			issue: []*gitlab.IssueAssignee{
				{
					ID:       1,
					Username: "foo",
				},
				{
					ID:       500,
					Username: "baz",
				},
			},
			expectedIDs: []int{500, 100},
			expectedAction: []string{
				"unassigned \"@foo\"",
				"assigned \"@bar\"",
			},
			ua: UserAssignments{
				ToAdd:    []string{"bar"},
				ToRemove: []string{"foo"},
			},
		},
		{
			name:  "remove foo (merge request)",
			users: []*gitlab.User{},
			merge: []*gitlab.BasicUser{
				{
					ID:       1,
					Username: "foo",
				},
			},
			expectedIDs:    []int{0},
			expectedAction: []string{"unassigned \"@foo\""},
			ua:             UserAssignments{ToRemove: []string{"foo"}},
		},
		{
			name:  "remove foo and baz out of foo, bar and baz (merge request)",
			users: []*gitlab.User{},
			merge: []*gitlab.BasicUser{
				{
					ID:       1,
					Username: "foo",
				},
				{
					ID:       2,
					Username: "bar",
				},
				{
					ID:       3,
					Username: "baz",
				},
			},
			expectedIDs:    []int{2},
			expectedAction: []string{"unassigned \"@foo @baz\""},
			ua:             UserAssignments{ToRemove: []string{"foo", "baz"}},
		},
		{
			name: "remove foo out of foo and baz and add bar (merge request)",
			users: []*gitlab.User{
				{
					ID:       100,
					Username: "bar",
				},
			},
			merge: []*gitlab.BasicUser{
				{
					ID:       1,
					Username: "foo",
				},
				{
					ID:       500,
					Username: "baz",
				},
			},
			expectedIDs: []int{500, 100},
			expectedAction: []string{
				"unassigned \"@foo\"",
				"assigned \"@bar\"",
			},
			ua: UserAssignments{
				ToAdd:    []string{"bar"},
				ToRemove: []string{"foo"},
			},
		},
		{
			name: "try to pass both issue and merge request users",
			issue: []*gitlab.IssueAssignee{
				{
					ID:       1,
					Username: "foo",
				},
			},
			merge: []*gitlab.BasicUser{
				{
					ID:       5,
					Username: "bar",
				},
			},
			wantErr: "issueAssignes and mergeRequestAssignes can't both not be nil",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			api.UsersByNames = func(_ *gitlab.Client, _ []string) ([]*gitlab.User, error) {
				return tC.users, nil
			}
			var gotAction []string
			gotIDs, gotAction, err := tC.ua.UsersFromAddRemove(tC.issue, tC.merge, &gitlab.Client{}, gotAction)
			if err != nil {
				if tC.wantErr != "" && tC.wantErr != err.Error() {
					t.Errorf("UsersFromAddRemove() expected error = %s, got = %s", tC.wantErr, err)
				} else if tC.wantErr == "" {
					t.Errorf("UsersFromAddRemove() unexpected error = %s", err)
				}
			}
			assert.ElementsMatch(t, gotIDs, tC.expectedIDs)
			assert.ElementsMatch(t, gotAction, tC.expectedAction)
		})
	}
}

func Test_ParseMilestoneTitleIsID(t *testing.T) {
	title := "1"
	expectedMilestoneID := 1

	// Override function to return an error, it should never reach this
	api.MilestoneByTitle = func(client *gitlab.Client, projectID interface{}, name string) (*gitlab.Milestone, error) {
		return nil, fmt.Errorf("We shouldn't have reached here")
	}

	got, err := ParseMilestone(&gitlab.Client{}, glrepo.New("foo", "bar"), title)
	if err != nil {
		t.Errorf("ParseMilestone() unexpected error = %s", err)
	}
	if got != expectedMilestoneID {
		t.Errorf("ParseMilestone() got = %d, expected = %d", got, expectedMilestoneID)
	}
}

func Test_ParseMilestoneAPIFail(t *testing.T) {
	title := "AsLongAsItDoesn'tConvertToInt"
	want := "api call failed in api.MilestoneByTitle()"

	// Override function to return an error simulating an API call failure
	api.MilestoneByTitle = func(client *gitlab.Client, projectID interface{}, name string) (*gitlab.Milestone, error) {
		return nil, fmt.Errorf("api call failed in api.MilestoneByTitle()")
	}

	_, err := ParseMilestone(&gitlab.Client{}, glrepo.New("foo", "bar"), title)
	if err == nil {
		t.Errorf("ParseMilestone() expected error")
	}
	if want != err.Error() {
		t.Errorf("ParseMilestone() expected error = %s, got error = %s", want, err)
	}
}

func Test_ParseMilestoneTitleToID(t *testing.T) {
	milestoneTitle := "kind: testing"
	expectedID := 3

	// Override function so it returns the correct milestone
	api.MilestoneByTitle = func(client *gitlab.Client, projectID interface{}, name string) (*gitlab.Milestone, error) {
		return &gitlab.Milestone{
				Title: "kind: testing",
				ID:    3,
			},
			nil
	}

	got, err := ParseMilestone(&gitlab.Client{}, glrepo.New("foo", "bar"), milestoneTitle)
	if err != nil {
		t.Errorf("ParseMilestone() unexpected error = %s", err)
	}
	if got != expectedID {
		t.Errorf("ParseMilestone() expected = %d, got = %d", expectedID, got)
	}
}

func Test_PickMetadata(t *testing.T) {
	const (
		labelsLabel    = "labels"
		assigneeLabel  = "assignees"
		milestoneLabel = "milestones"
	)

	testCases := []struct {
		name     string
		values   []string
		expected []Action
	}{
		{
			name: "nothing picked",
		},
		{
			name:     "labels",
			values:   []string{labelsLabel},
			expected: []Action{AddLabelAction},
		},
		{
			name:     "assignees",
			values:   []string{assigneeLabel},
			expected: []Action{AddAssigneeAction},
		},
		{
			name:     "milestone",
			values:   []string{milestoneLabel},
			expected: []Action{AddMilestoneAction},
		},
		{
			name:     "labels and assignees",
			values:   []string{labelsLabel, assigneeLabel},
			expected: []Action{AddLabelAction, AddAssigneeAction},
		},
		{
			name:     "labels and milestone",
			values:   []string{labelsLabel, milestoneLabel},
			expected: []Action{AddLabelAction, AddMilestoneAction},
		},
		{
			name:     "assignees and milestone",
			values:   []string{assigneeLabel, milestoneLabel},
			expected: []Action{AddAssigneeAction, AddMilestoneAction},
		},
		{
			name:     "labels, assignees and milestone",
			values:   []string{labelsLabel, assigneeLabel, milestoneLabel},
			expected: []Action{AddLabelAction, AddAssigneeAction, AddMilestoneAction},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			as, restoreAsk := prompt.InitAskStubber()
			defer restoreAsk()

			as.Stub([]*prompt.QuestionStub{
				{
					Name:  "metadata",
					Value: tC.values,
				},
			})

			got, err := PickMetadata()
			if err != nil {
				t.Errorf("PickMetadata() unexpected error = %s", err)
			}
			assert.ElementsMatch(t, got, tC.expected)
		})
	}
}

func Test_AssigneesPrompt(t *testing.T) {
	testCases := []struct {
		name   string
		input  string
		output []string
	}{
		{
			name: "nothing",
		},
		{
			name:   "Single name",
			input:  "foo",
			output: []string{"foo"},
		},
		{
			name:   "2 or more names",
			input:  "foo,bar",
			output: []string{"foo", "bar"},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			as, restoreAsk := prompt.InitAskStubber()
			defer restoreAsk()

			as.Stub([]*prompt.QuestionStub{
				{
					Name:  "assignee",
					Value: tC.input,
				},
			})

			var got []string
			err := AssigneesPrompt(&got)
			if err != nil {
				t.Errorf("AssigneesPrompt() unexpected error = %s", err)
			}
			assert.ElementsMatch(t, got, tC.output)
		})
	}
}

func Test_MilestonesPrompt(t *testing.T) {
	mockMilestones := []*gitlab.Milestone{
		{
			Title: "New Release",
			ID:    5,
		},
		{
			Title: "Really big feature",
			ID:    240,
		},
		{
			Title: "Get rid of low quality code",
			ID:    650,
		},
	}

	// Override API.ListMilestones so it doesn't make any network calls
	api.ListMilestones = func(_ *gitlab.Client, _ interface{}, _ *gitlab.ListMilestonesOptions) ([]*gitlab.Milestone, error) {
		return mockMilestones, nil
	}

	// mock glrepo.Remote object
	repo := glrepo.New("foo", "bar")
	remote := &git.Remote{
		Name:     "test",
		Resolved: "base",
	}
	repoRemote := &glrepo.Remote{
		Remote: remote,
		Repo:   repo,
	}

	testCases := []struct {
		name       string
		input      string // Selected milestone
		expectedID int    // expected global ID from the milestone
	}{
		{
			name:       "match",
			input:      "New Release",
			expectedID: 5,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			as, restoreAsk := prompt.InitAskStubber()
			defer restoreAsk()

			as.Stub([]*prompt.QuestionStub{
				{
					Name:  "milestone",
					Value: tC.input,
				},
			})

			var got int
			var io utils.IOStreams

			err := MilestonesPrompt(&got, &gitlab.Client{}, repoRemote, &io)
			if err != nil {
				t.Errorf("MilestonesPrompt() unexpected error = %s", err)
			}
			if got != 0 && got != tC.expectedID {
				t.Errorf("MilestonesPrompt() expected = %d, got = %d", got, tC.expectedID)
			}
		})
	}
}

func Test_MilestonesPromptNoPrompts(t *testing.T) {
	// Override api.ListMilestones so it returns an empty slice, we are testing if MilestonesPrompt()
	// will print the correct message to `stderr` when it tries to get the list of Milestones in a
	// project but the project has no milestones
	api.ListMilestones = func(_ *gitlab.Client, _ interface{}, _ *gitlab.ListMilestonesOptions) ([]*gitlab.Milestone, error) {
		return []*gitlab.Milestone{}, nil
	}

	// mock glrepo.Remote object
	repo := glrepo.New("foo", "bar")
	remote := &git.Remote{
		Name:     "test",
		Resolved: "base",
	}
	repoRemote := &glrepo.Remote{
		Remote: remote,
		Repo:   repo,
	}

	var got int
	io, _, _, stderr := utils.IOTest()

	err := MilestonesPrompt(&got, &gitlab.Client{}, repoRemote, io)
	if err != nil {
		t.Errorf("MilestonesPrompt() unexpected error = %s", err)
	}
	assert.Equal(t, "There are no active milestones in this project\n", stderr.String())
}

func TestMilestonesPromptFailures(t *testing.T) {
	// Override api.ListMilestones so it returns an error, we are testing to see if error
	// handling from the usage of api.ListMilestones is correct
	api.ListMilestones = func(_ *gitlab.Client, _ interface{}, _ *gitlab.ListMilestonesOptions) ([]*gitlab.Milestone, error) {
		return nil, errors.New("api.ListMilestones() failed")
	}

	// mock glrepo.Remote object
	repo := glrepo.New("foo", "bar")
	remote := &git.Remote{
		Name:     "test",
		Resolved: "base",
	}
	repoRemote := &glrepo.Remote{
		Remote: remote,
		Repo:   repo,
	}

	var got int
	io, _, _, _ := utils.IOTest()

	err := MilestonesPrompt(&got, &gitlab.Client{}, repoRemote, io)
	if err == nil {
		t.Error("MilestonesPrompt() expected error")
	}
	assert.Equal(t, "api.ListMilestones() failed", err.Error())
}

func Test_IDsFromUsers(t *testing.T) {
	testCases := []struct {
		name  string
		users []*gitlab.User // Mock of the gitlab.User object
		IDs   []int          // IDs we expect from the users
	}{
		{
			name: "no users",
		},
		{
			name: "one user",
			users: []*gitlab.User{
				{
					ID: 1,
				},
			},
			IDs: []int{1},
		},
		{
			name: "multiple users",
			users: []*gitlab.User{
				{
					ID: 3,
				},
				{
					ID: 6,
				},
				{
					ID: 2,
				},
				{
					ID: 51,
				},
				{
					ID: 32,
				},
				{
					ID: 87,
				},
				{
					ID: 210,
				},
				{
					ID: 6493,
				},
				{
					ID: 50132,
				},
			},
			IDs: []int{
				50132,
				6493,
				210,
				87,
				32,
				51,
				2,
				3,
				6,
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			got := IDsFromUsers(tC.users)
			assert.ElementsMatch(t, got, tC.IDs)
		})
	}
}
