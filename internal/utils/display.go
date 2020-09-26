package utils

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gosuri/uitable"
)

// ListInfo represents the parameters required to display a list result.
type ListInfo struct {
	// Name of the List to be used in constructing Description and EmptyMessage if not provided.
	Name string
	// List of columns to display
	Columns []string
	// Total number of record. Ideally size of the List.
	Total int
	// Function to pick a cell value from cell index
	GetCellValue func(int, int) interface{}
	// Optional. Description of the List. If not provided, default one constructed from list Name.
	Description string
	// Optional. EmptyMessage to display when List is empty. If not provided, default one constructed from list Name.
	EmptyMessage string
	// TableWrap wraps the contents when the column length exceeds the maximum width
	TableWrap bool
}

// Prints the list data on console
func DisplayList(lInfo ListInfo, projectID string) *uitable.Table {
	table := uitable.New()
	table.MaxColWidth = 70
	table.Wrap = lInfo.TableWrap

	if lInfo.Total > 0 {
		description := lInfo.Description
		if description == "" {
			description = fmt.Sprintf("Showing %s %d of %d on %s\n", lInfo.Name, lInfo.Total, lInfo.Total, projectID)
		}
		table.AddRow(description)
		header := make([]interface{}, len(lInfo.Columns))
		for ci, c := range lInfo.Columns {
			header[ci] = c
		}
		table.AddRow(header...)

		for ri := 0; ri < lInfo.Total; ri++ {
			row := make([]interface{}, len(lInfo.Columns))
			for ci := range lInfo.Columns {
				row[ci] = lInfo.GetCellValue(ri, ci)
			}
			table.AddRow(row...)
		}
	} else {
		emptyMessage := lInfo.EmptyMessage
		if emptyMessage == "" {
			emptyMessage = fmt.Sprintf("No %s available on %s", lInfo.Name, projectID)
		}
		table.AddRow(emptyMessage)
	}

	return table
}

var lineRE = regexp.MustCompile(`(?m)^`)

func Indent(s, indent string) string {
	if len(strings.TrimSpace(s)) == 0 {
		return s
	}
	return lineRE.ReplaceAllLiteralString(s, indent)
}
