package tableprinter

import (
	"bytes"
	"testing"
)

func Test_ttyTablePrinter_truncate(t *testing.T) {
	buf := bytes.Buffer{}
	tp := NewTablePrinter()
	tp.SetTTYSeparator(" ")
	tp.SetTerminalWidth(5)
	tp.SetIsTTY(true)

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

func Test_nonTTYTablePrinter_truncate(t *testing.T) {
	buf := bytes.Buffer{}
	tp := NewTablePrinter()
	tp.SetTerminalWidth(5)
	tp.SetIsTTY(false)

	tp.AddCell("1")
	tp.AddCell("hello")
	tp.EndRow()
	tp.AddCell("2")
	tp.AddCell("world")
	tp.EndRow()

	buf.Write(tp.Bytes())

	expected := "1\thello\n2\tworld\n"
	if buf.String() != expected {
		t.Errorf("expected: %q, got: %q", expected, buf.Bytes())
	}
}
