package cli

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// Specializes in formatting a type into Column format
type Formatter interface {
	ToColumns(output io.Writer) error
}

// Handles writing the formatted type into the desired output and global formatting
type CLIWriter struct {
	FormatWriter func(f Formatter)
	ErrorWriter  func(err error)
}

var writer *CLIWriter

func Register(w *CLIWriter) {
	writer = w
}

func Output(f Formatter, err error) {
	if err == nil {
		writer.FormatWriter(f)
	} else {
		writer.ErrorWriter(err)
	}
}

func FlushWriter(w *tabwriter.Writer) {
	fmt.Fprintln(w, "")
	w.Flush()
}

func NewTabWriter(output io.Writer) *tabwriter.Writer {
	w := new(tabwriter.Writer)
	w.Init(output, 0, 8, 2, '\t', 0)
	return w
}
