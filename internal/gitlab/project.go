package gitlab

import "github.com/xanzy/go-gitlab"

func GetProject(gLab *gitlab.Client, projectID interface{}) (*gitlab.Project, error) {
	if gLab == nil {
		gLab = gLabClient
	}
	opts := &gitlab.GetProjectOptions{
		Statistics:           gitlab.Bool(true),
		License:              gitlab.Bool(true),
		WithCustomAttributes: gitlab.Bool(true),
	}
	project, _, err := gLab.Projects.GetProject(projectID, opts)
	if err != nil {
		return nil, err
	}
	return project, nil
}
func GetRepository(gLab *gitlab.Client, projectID interface{}) (*gitlab.Project, error) {
	if gLab == nil {
		gLab = gLabClient
	}
	opts := &gitlab.GetProjectOptions{
		Statistics:           gitlab.Bool(true),
		License:              gitlab.Bool(true),
		WithCustomAttributes: gitlab.Bool(true),
	}
	project, _, err := gLab.Projects.GetProject(projectID, opts)
	if err != nil {
		return nil, err
	}
	return project, nil
}

func DeleteProject(gLab *gitlab.Client, projectID interface{}) (*gitlab.Response, error) {
	if gLab == nil {
		gLab = gLabClient
	}
	project, err := gLab.Projects.DeleteProject(projectID)
	if err != nil {
		return nil, err
	}
	return project, nil
}

func CreateProject(gLab *gitlab.Client, opts *gitlab.CreateProjectOptions) (*gitlab.Project, error) {
	if gLab == nil {
		gLab = gLabClient
	}
	project, _, err := gLab.Projects.CreateProject(opts)
	if err != nil {
		return nil, err
	}
	return project, nil
}

func GetGroup(gLab *gitlab.Client, groupID interface{}) (*gitlab.Group, error) {
	if gLab == nil {
		gLab = gLabClient
	}
	group, _, err := gLab.Groups.GetGroup(groupID)
	if err != nil {
		return nil, err
	}
	return group, nil
}
