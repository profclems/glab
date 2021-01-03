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
		Name:           strings.TrimSpace(listName),
		ListActionType: "list",
		Page:           1,
	}
}

func (opts *ListTitleOptions) Describe() string {
	var pageNumInfo string
	var pageInfo string

	if opts.Total != 0 {
		opts.Name = pluralizeName(opts.Total, opts.Name)
		pageNumInfo = fmt.Sprintf("%d of %d", opts.CurrentPageTotal, opts.Total)
	} else {
		opts.Name = pluralizeName(opts.CurrentPageTotal, opts.Name)
		pageNumInfo = fmt.Sprintf("%d", opts.CurrentPageTotal)
	}

	if opts.Page != 0 {
		pageInfo = fmt.Sprintf("(Page %d)", opts.Page)
	}

	if opts.ListActionType == "search" {
		if opts.CurrentPageTotal > 0 {
			return fmt.Sprintf("Showing %s %s in %s that match your search %s\n", pageNumInfo, opts.Name,
				opts.RepoName, pageInfo)
		}

		return fmt.Sprintf("No %s match your search in %s\n", opts.Name, opts.RepoName)
	}

	if opts.CurrentPageTotal > 0 {
		return fmt.Sprintf("Showing %s %s on %s %s\n", pageNumInfo, opts.Name, opts.RepoName, pageInfo)
	}

	emptyMessage := opts.EmptyMessage
	if emptyMessage == "" {
		return fmt.Sprintf("No %s available on %s", opts.Name, opts.RepoName)
	}

	return emptyMessage
}

func pluralizeName(num int, thing string) string {
	return strings.TrimPrefix(Pluralize(num, thing), fmt.Sprintf("%d ", num))
}
