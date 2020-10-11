package tableprinter

import (
	"fmt"

	"github.com/gosuri/uitable"
)

var (
	DefaultSeparator = "\t"
	DefaultMaxColWidth uint = 70
)

// TablePrinter represents a decorator that renders the data formatted in a tabular form.
type TablePrinter struct {
	// List of columns to display
	Header []string
	// Total number of records. Needed if AddRowFunc is used
	TotalRows int
	// Wrap when set to true wraps the contents of the columns when the length exceeds the MaxColWidth
	Wrap bool
	// MaxColWidth is the maximum allowed width for cells in the table
	MaxColWidth uint
	// Separator is the seperator for columns in the table. Default is "\t
	Separator string
	// Rows is the collection of rows in the table
	Rows [][]*TableCell
}

type TableCell struct {
	Text interface{}
}

func NewTablePrinter() TablePrinter {
	return TablePrinter{
		Separator:   DefaultSeparator,
		MaxColWidth: DefaultMaxColWidth,
		Wrap:        false,
	}
}

func (t *TablePrinter) AddCell(s interface{})  {
	if t.Rows == nil {
		t.Rows = make([][]*TableCell, 1)
	}
	rowI := len(t.Rows) - 1
	field := TableCell{
		Text: s,
	}
	t.Rows[rowI] = append(t.Rows[rowI], &field)
}

func (t *TablePrinter) appendCellToIndex(s interface{}, index int)  {
	if t.Rows == nil {
		t.Rows = make([][]*TableCell, 1)
	}
	rowI := len(t.Rows) - 1
	last := len(t.Rows[rowI]) - 1

	if last <= index {
		t.AddCell(s)
		return
	}

	t.Rows[rowI] = append(t.Rows[rowI], t.Rows[rowI][last])
	copy(t.Rows[rowI][(index+1):], t.Rows[rowI][index:last])

	field := TableCell{
		Text: s,
	}
	t.Rows[rowI][index] = &field
}

// AddCellf formats according to a format specifier and adds cell to row
func (t *TablePrinter) AddCellf(s string, f ...interface{})  {
	t.AddCell(fmt.Sprintf(s, f...))
}

func (t *TablePrinter) AddRow(str ...interface{})  {
	for _, s := range str {
		t.AddCell(s)
	}
	t.EndRow()
}

// AddCellI appends a cell to the given index
func (t *TablePrinter) AddCellI(index int, str interface{})  {
	t.appendCellToIndex(str, index)
}

func (t *TablePrinter) AddRowFunc(f func(int, int) string) {
	for ri := 0; ri < t.TotalRows; ri++ {
		row := make([]interface{}, t.TotalRows)
		for ci := range row {
			row[ci] = f(ri, ci)
		}
		t.AddRow(row)
		t.EndRow()
	}
}

func (t *TablePrinter) EndRow() {
	t.Rows = append(t.Rows, []*TableCell{})
}

// Bytes returns the []byte value of table
func (t *TablePrinter) Bytes() []byte {
	return []byte(t.String())
}

// String returns the string value of table. Alternative to Render()
func (t *TablePrinter) String() string {
	return t.Render()
}

// purgeRow removes empty rows
func (t *TablePrinter) purgeRow()  {
	newSlice := make([][]*TableCell, 0, len(t.Rows))
	for _, item := range t.Rows {
		if len(item) > 0 {
			newSlice = append(newSlice, item)
		}
	}
	t.Rows = newSlice
}

// Render builds and returns the string representation of the table
func (t *TablePrinter) Render() string {
	table := uitable.New()

	if t.MaxColWidth != 0 {
		table.MaxColWidth = t.MaxColWidth
	}

	if t.Separator != "" {
		table.Separator = t.Separator
	}

	t.purgeRow() // remove empty rows
	rLen := len(t.Rows)
	fmt.Println(rLen)

	if rLen > 0 {
		if len(t.Header) > 0 {
			header := make([]interface{}, len(t.Header), len(t.Header))
			for ci, c := range t.Header {
				header[ci] = c
			}

			table.AddRow(header...)
		}

		for _, row := range t.Rows {
			rowData := make([]interface{}, len(row), len(row))
			for i, r := range row {
				if len(rowData) <= i {
					rowData[i-1] = r.Text
					continue
				}
				rowData[i] = r.Text
			}
			table.AddRow(rowData...)
		}

		return table.String()
	}

	return ""
}
