package tableprinter

import (
	"bytes"
	"testing"
)

func Test_ttyTablePrinter_truncate(t *testing.T) {
	buf := bytes.Buffer{}
	tp := NewTablePrinter()
	tp.SetSeparator(" ")
	tp.SetTerminalWidth(5)

	tp.AddCell("1")
	tp.AddCell("hello")
	tp.EndRow()
	tp.AddCell("2")
	tp.AddCell("world")
	tp.EndRow()

	buf.Write(tp.Bytes())

	expected := "1 h...\n2 w...\n"
	if buf.String() != expected {
		t.Errorf("expected: %q, got: %q", expected, buf.Bytes())
	}
}
