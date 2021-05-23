package iostreams

import "fmt"

// Log prints output to StdErr
func (s *IOStreams) Log(msg ...interface{}) {
	fmt.Fprintln(s.StdErr, msg...)
}

// Logf formats according to a format specifier and writes to StdErr
func (s *IOStreams) Logf(format string, a ...interface{}) {
	fmt.Fprintf(s.StdErr, format, a...)
}

// LogInfo is just like Log but prints output to StdOut
func (s *IOStreams) LogInfo(a ...interface{}) {
	fmt.Fprintln(s.StdOut, a...)
}

// LogInfof formats according to a format specifier and writes to StdOut
func (s *IOStreams) LogInfof(format string, a ...interface{}) {
	fmt.Fprintf(s.StdOut, format, a...)
}
