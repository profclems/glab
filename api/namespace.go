package api

import "github.com/xanzy/go-gitlab"

var ListNamespaces = func(client *gitlab.Client, opts *gitlab.ListNamespacesOptions) ([]*gitlab.Namespace, error) {
	if client == nil {
		client = apiClient.Lab()
	}
	namespaces, _, err := client.Namespaces.ListNamespaces(opts)
	if err != nil {
		return nil, err
	}
	return namespaces, nil
}
