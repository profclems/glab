package cmdutils

import (
	"testing"

	"github.com/profclems/glab/pkg/api"
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
			apiClient := gitlab.Client{} // Empty Client, it won't be used, just to satisfy the function signature
			gotIDs, gotAction, err := ua.UsersFromReplaces(&apiClient, gotAction)
			if err != nil {
				t.Errorf("UsersFromReplaces() unexpected error = %s", err)
			}
			assert.ElementsMatch(t, gotIDs, tC.expectedIDs)
			assert.ElementsMatch(t, gotAction, tC.expectedAction)
		})
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
			apiClient := gitlab.Client{} // Empty Client, it won't be used, just to satisfy the function signature
			gotIDs, gotAction, err := tC.ua.UsersFromAddRemove(tC.issue, tC.merge, &apiClient, gotAction)
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
