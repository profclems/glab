package tableprinter

import (
	"fmt"
	"strings"

	"github.com/profclems/glab/pkg/text"
	"github.com/spf13/cast"
)

var tp *TablePrinter

func init() {
	tp = &TablePrinter{
		TotalRows:       0,
		Wrap:            false,
		MaxColWidth:     0,
		TTYSeparator:    "\t",
		NonTTYSeparator: "\t",
		TerminalWidth:   80,
	}
}

// TablePrinter represents a decorator that renders the data formatted in a tabular form.
type TablePrinter struct {
	// Total number of records. Needed if AddRowFunc is used
	TotalRows int
	// Wrap when set to true wraps the contents of the columns when the length exceeds the MaxColWidth
	Wrap bool
	// MaxColWidth is the maximum allowed width for cells in the table
	MaxColWidth int
	// TTYSeparator is the separator for columns in the table on TTYs. Default is "\t"
	TTYSeparator string
	// NonTTYSeparator is the separator for columns in the table on non-TTYs. Default is "\t"
	NonTTYSeparator string
	// Rows is the collection of rows in the table
	Rows []*TableRow
	// TerminalWidth is the max width of the terminal
	TerminalWidth int
	// IsTTY indicates whether output is a TTY or non-TTY
	IsTTY bool
}

type TableCell struct {
	// Value in the cell
	Value interface{}
	// Width is the width of the cell
	Width int
	// Wrap when true wraps the contents of the cell when the length exceeds the width
	Wrap bool

	isaTTY bool
}

type TableRow struct {
	Cells []*TableCell
	// Separator is the seperator for columns in the table. Default is " "
	Separator string
}

func NewTablePrinter() *TablePrinter {
	t := &TablePrinter{
		TTYSeparator:    tp.TTYSeparator,
		NonTTYSeparator: tp.NonTTYSeparator,
		MaxColWidth:     tp.MaxColWidth,
		Wrap:            false,
		TerminalWidth:   tp.TerminalWidth,
		IsTTY:           tp.IsTTY,
	}

	return t
}

func (t *TablePrinter) Separator() string {
	if t.IsTTY {
		return t.TTYSeparator
	}

	return t.NonTTYSeparator
}

// SetTerminalWidth sets the maximum width for the terminal
func SetTerminalWidth(width int) { tp.SetTerminalWidth(width) }

func (t *TablePrinter) SetTerminalWidth(width int) {
	t.TerminalWidth = width
}

// SetIsTTY sets the IsTTY variable which indicates whether terminal
// output is a TTY or nonTTY
func SetIsTTY(isTTY bool) { tp.SetIsTTY(isTTY) }

func (t *TablePrinter) SetIsTTY(isTTY bool) {
	t.IsTTY = isTTY
}

// SetTTYSeparator sets the separator for the columns in the table for TTYs
func SetTTYSeparator(s string) { tp.SetTTYSeparator(s) }

func (t *TablePrinter) SetTTYSeparator(s string) {
	t.TTYSeparator = s
}

// SetNonTTYSeparator sets the separator for the columns in the table for non-ttys
func SetNonTTYSeparator(s string) { tp.SetNonTTYSeparator(s) }

func (t *TablePrinter) SetNonTTYSeparator(s string) {
	t.NonTTYSeparator = s
}

func (t *TablePrinter) makeRow() {
	if t.Rows == nil {
		t.Rows = make([]*TableRow, 1)
		t.Rows[0] = &TableRow{}
	}
}

func (t *TablePrinter) AddCell(s interface{}) {
	t.makeRow()
	rowI := len(t.Rows) - 1
	row := t.Rows[rowI]

	cell := &TableCell{
		Value:  s,
		isaTTY: t.IsTTY,
	}

	row.Separator = t.Separator()

	row.Cells = append(row.Cells, cell)
}

// AddCellf formats according to a format specifier and adds cell to row
func (t *TablePrinter) AddCellf(s string, f ...interface{}) {
	t.AddCell(fmt.Sprintf(s, f...))
}

func (t *TablePrinter) AddRow(str ...interface{}) {
	for _, s := range str {
		t.AddCell(s)
	}
	t.EndRow()
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
	t.Rows = append(t.Rows, &TableRow{Cells: make([]*TableCell, 1)})
}

// Bytes returns the []byte value of table
func (t *TablePrinter) Bytes() []byte {
	return []byte(t.String())
}

// String returns the string value of table. Alternative to Render()
func (t *TablePrinter) String() string {
	return t.Render()
}

// String returns the string representation of the row
func (r *TableRow) String() string {
	// get the max number of lines for each cell
	var lc int // line count
	for _, cell := range r.Cells {
		if clc := len(strings.Split(cell.String(), "\n")); clc > lc {
			lc = clc
		}
	}

	// allocate a two-dimentional array of cells for each line and add size them
	cells := make([][]*TableCell, lc)
	for x := 0; x < lc; x++ {
		cells[x] = make([]*TableCell, len(r.Cells))
		for y := 0; y < len(r.Cells); y++ {
			cells[x][y] = &TableCell{Width: r.Cells[y].Width}
		}
	}

	// insert each line in a cell as new cell in the cells array
	for y, cell := range r.Cells {
		lines := strings.Split(cell.String(), "\n")
		for x, line := range lines {
			cells[x][y].Value = line
		}
	}

	// format each line
	lines := make([]string, lc)
	for x := range lines {
		line := make([]string, len(cells[x]))
		for y := range cells[x] {
			line[y] = cells[x][y].String()
		}
		lines[x] = text.Join(line, r.Separator)
	}
	return strings.Join(lines, "\n")
}

// purgeRow removes nil cells and rows
func (t *TablePrinter) purgeRow() {
	newSlice := make([]*TableRow, 0, len(t.Rows))
	for _, row := range t.Rows {
		var newRow *TableRow
		if len(row.Cells) > 0 && row.Cells != nil {
			var newCells []*TableCell
			for _, cell := range row.Cells {
				if cell != nil {
					newCells = append(newCells, cell)
				}
			}
			newRow = &TableRow{Cells: newCells}
		}

		if newRow != nil {
			newSlice = append(newSlice, newRow)
		}
	}
	t.Rows = newSlice
}

// Render builds and returns the string representation of the table
func (t *TablePrinter) Render() string {
	if len(t.Rows) == 0 {
		return ""
	}
	// remove nil cells and rows
	t.purgeRow()

	colWidths := t.colWidths()

	var lines []string
	for _, row := range t.Rows {
		row.Separator = t.Separator()
		for i, cell := range row.Cells {
			cell.Width = colWidths[i]
			cell.Wrap = t.Wrap
		}
		lines = append(lines, row.String())
	}
	return text.Join(lines, "\n")
}

// LineWidth returns the max width of all the lines in a cell
func (c *TableCell) LineWidth() int {
	width := 0
	for _, s := range strings.Split(c.String(), "\n") {
		w := text.StringWidth(s)
		if w > width {
			width = w
		}
	}
	return width
}

// String returns the string formatted representation of the cell
func (c *TableCell) String() string {
	if c == nil {
		return ""
	}
	if c.Value == nil {
		return text.PadLeft(" ", c.Width, ' ')
	}

	s := cast.ToString(c.Value)
	if c.Width > 0 && c.isaTTY {
		if c.Wrap && len(s) > c.Width {
			return text.WrapString(s, c.Width)
		} else {
			return text.Truncate(s, c.Width)
		}
	}
	return s
}

// colWidths determine the width for each column (cell in a row)
func (t *TablePrinter) colWidths() []int {
	var colWidths []int
	for _, row := range t.Rows {
		for i, cell := range row.Cells {
			// resize colwidth array
			if i+1 > len(colWidths) {
				colWidths = append(colWidths, 0)
			}
			cellwidth := cell.LineWidth()
			if t.MaxColWidth != 0 && cellwidth > t.MaxColWidth {
				cellwidth = t.MaxColWidth
			}

			if cellwidth > colWidths[i] {
				colWidths[i] = cellwidth
			}
		}
	}
	numCols := len(colWidths)
	separatorWidth := (numCols - 1) * len(t.Separator())
	totalWidth := separatorWidth
	for _, width := range colWidths {
		totalWidth += width
	}

	if t.MaxColWidth == 0 && totalWidth > t.TerminalWidth {
		availWidth := t.TerminalWidth - colWidths[0] - separatorWidth
		// add extra space from columns that are already narrower than threshold
		for col := 1; col < numCols; col++ {
			availColWidth := availWidth / (numCols - 1)
			if extra := availColWidth - colWidths[col]; extra > 0 {
				availWidth += extra
			}
		}
		// cap all but first column to fit available terminal width
		for col := 1; col < numCols; col++ {
			availColWidth := availWidth / (numCols - 1)
			if colWidths[col] > availColWidth {
				colWidths[col] = availColWidth
			}
		}
	}

	return colWidths
}
