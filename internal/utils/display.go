package utils

import (
	"fmt"
	"regexp"
	"strings"
)

var lineRE = regexp.MustCompile(`(?m)^`)

func Indent(s, indent string) string {
	if strings.TrimSpace(s) == "" {
		return s
	}
	return lineRE.ReplaceAllLiteralString(s, indent)
}

type ListTitleOptions struct {
	// Name of the List to be used in constructing Description and EmptyMessage if not provided.
	Name string
	// Page represents the page number of the current page
	Page int
	// CurrentPageTotal is the total number of items in current page
	CurrentPageTotal int
	// Total number of records. Default is the total number of rows.
	// Can be set to be greater than the total number of rows especially, if the list is paginated
	Total int
	// RepoName represents the name of the project or repository
	RepoName string
	// ListActionType should be either "search" or "list". Default is list
	ListActionType string
	// Optional. EmptyMessage to display when List is empty. If not provided, default one constructed from list Name.
	EmptyMessage string
}

func NewListTitle(listName string) ListTitleOptions {
	return ListTitleOptions{
		Name:           listName,
		ListActionType: "list",
	}
}

func (opts *ListTitleOptions) Describe() string {
	if opts.ListActionType == "search" {
		if opts.CurrentPageTotal > 0 {
			return fmt.Sprintf("Showing %d of %d %s in %s that match your search", opts.CurrentPageTotal,
				opts.Total, opts.Name, opts.RepoName)
		}

		return fmt.Sprintf("No %s match your search in %s", opts.Name, opts.RepoName)
	}

	if opts.CurrentPageTotal > 0 {
		return fmt.Sprintf("Showing %s %d of %d on %s\n", opts.Name, opts.CurrentPageTotal, opts.Total, opts.RepoName)
	}

	emptyMessage := opts.EmptyMessage
	if emptyMessage == "" {
		return fmt.Sprintf("No %s available on %s", opts.Name, opts.RepoName)
	}

	return emptyMessage
}
