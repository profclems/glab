package api

import (
	"fmt"

	"github.com/xanzy/go-gitlab"
	"golang.org/x/sync/errgroup"
)

// Describe namespace kinds which is either group or user
// See docs: https://docs.gitlab.com/ee/api/namespaces.html
const (
	NamespaceKindUser  = "user"
	NamespaceKindGroup = "group"
)

type Milestone struct {
	ID    int
	Title string
}

func NewProjectMilestone(m *gitlab.Milestone) *Milestone {
	return &Milestone{
		ID:    m.ID,
		Title: m.Title,
	}
}

func NewGroupMilestone(m *gitlab.GroupMilestone) *Milestone {
	return &Milestone{
		ID:    m.ID,
		Title: m.Title,
	}
}

type ListMilestonesOptions struct {
	IIDs                    []int
	State                   *string
	Title                   *string
	Search                  *string
	IncludeParentMilestones *bool
	PerPage                 int
	Page                    int
}

func (opts *ListMilestonesOptions) ListProjectMilestonesOptions() *gitlab.ListMilestonesOptions {
	projectOpts := &gitlab.ListMilestonesOptions{
		IIDs:   &opts.IIDs,
		State:  opts.State,
		Title:  opts.Title,
		Search: opts.Search,
	}
	projectOpts.PerPage = opts.PerPage
	projectOpts.Page = opts.Page
	return projectOpts
}

func (opts *ListMilestonesOptions) ListGroupMilestonesOptions() *gitlab.ListGroupMilestonesOptions {
	groupOpts := &gitlab.ListGroupMilestonesOptions{
		IIDs:                    &opts.IIDs,
		State:                   opts.State,
		Title:                   opts.Title,
		Search:                  opts.Search,
		IncludeParentMilestones: opts.IncludeParentMilestones,
	}
	groupOpts.PerPage = opts.PerPage
	groupOpts.Page = opts.Page
	return groupOpts
}

var ListGroupMilestones = func(client *gitlab.Client, groupID interface{}, opts *gitlab.ListGroupMilestonesOptions) ([]*gitlab.GroupMilestone, error) {
	if client == nil {
		client = apiClient.Lab()
	}

	if opts.PerPage == 0 {
		opts.PerPage = DefaultListLimit
	}

	milestone, _, err := client.GroupMilestones.ListGroupMilestones(groupID, opts)
	if err != nil {
		return nil, err
	}
	return milestone, nil
}

var ListProjectMilestones = func(client *gitlab.Client, projectID interface{}, opts *gitlab.ListMilestonesOptions) ([]*gitlab.Milestone, error) {
	if client == nil {
		client = apiClient.Lab()
	}

	if opts.PerPage == 0 {
		opts.PerPage = DefaultListLimit
	}

	milestone, _, err := client.Milestones.ListMilestones(projectID, opts)
	if err != nil {
		return nil, err
	}
	return milestone, nil
}

var ProjectMilestoneByTitle = func(client *gitlab.Client, projectID interface{}, name string) (*gitlab.Milestone, error) {
	opts := &gitlab.ListMilestonesOptions{Title: gitlab.String(name)}

	if client == nil {
		client = apiClient.Lab()
	}

	if opts.PerPage == 0 {
		opts.PerPage = DefaultListLimit
	}

	milestones, _, err := client.Milestones.ListMilestones(projectID, opts)
	if err != nil {
		return nil, err
	}

	if len(milestones) != 1 {
		return nil, fmt.Errorf("failed to find milestone by title: %s", name)
	}

	return milestones[0], nil
}

var ListAllMilestones = func(client *gitlab.Client, projectID interface{}, opts *ListMilestonesOptions) ([]*Milestone, error) {
	project, err := GetProject(client, projectID)
	if err != nil {
		return nil, err
	}

	errGroup := &errgroup.Group{}
	projectMilestones := []*gitlab.Milestone{}
	groupMilestones := []*gitlab.GroupMilestone{}

	errGroup.Go(func() error {
		var err error
		projectMilestones, err = ListProjectMilestones(client, projectID, opts.ListProjectMilestonesOptions())
		return err
	})

	if project.Namespace.Kind == NamespaceKindGroup {
		errGroup.Go(func() error {
			var err error
			groupMilestones, err = ListGroupMilestones(client, project.Namespace.ID, opts.ListGroupMilestonesOptions())
			return err
		})
	}

	if err := errGroup.Wait(); err != nil {
		return nil, fmt.Errorf("failed to get all project related milestones. %w", err)
	}

	milestones := make([]*Milestone, 0, len(projectMilestones)+len(groupMilestones))
	for _, v := range projectMilestones {
		milestones = append(milestones, NewProjectMilestone(v))
	}

	for _, v := range groupMilestones {
		milestones = append(milestones, NewGroupMilestone(v))
	}

	return milestones, nil
}
