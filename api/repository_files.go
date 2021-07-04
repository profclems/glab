package api

import "github.com/xanzy/go-gitlab"

// GetFile retrieves a file from repository. Note that file content is Base64 encoded.
var GetFile = func(client *gitlab.Client, projectID interface{}, path string, ref string) (*gitlab.File, error) {
	if client == nil {
		client = apiClient.Lab()
	}

	fileOpts := &gitlab.GetFileOptions{
		Ref: &ref,
	}
	file, _, err := client.RepositoryFiles.GetFile(projectID, path, fileOpts)

	if err != nil {
		return nil, err
	}

	return file, nil
}
