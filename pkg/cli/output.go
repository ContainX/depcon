package cli

import (
	"fmt"
	"io"
	"strconv"
	"text/tabwriter"
	"text/template"
	"time"
)

// Specializes in formatting a type into Column format
type Formatter interface {
	ToColumns(output io.Writer) error
	Data() FormatData
}

type FormatData struct {
	Data     interface{}
	Template string
	Funcs    template.FuncMap
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

func (d FormatData) ToColumns(output io.Writer) error {
	w := NewTabWriter(output)
	t := template.New("output").Funcs(buildFuncMap(d.Funcs))
	t, _ = t.Parse(d.Template)
	if err := t.Execute(w, d.Data); err != nil {
		return err
	}
	FlushWriter(w)
	return nil
}

func buildFuncMap(userFuncs template.FuncMap) template.FuncMap {
	funcMap := template.FuncMap{
		"floatToString": floatToString,
		"intToString":   strconv.Itoa,
		"valString":     valueToString,
		"pad":           padString,
		"fdate":         FormatDate,
		"msDur":         durationToMilliseconds,
		"boolToYesNo":   boolToYesNo,
	}

	if userFuncs != nil {
		for k, v := range userFuncs {
			funcMap[k] = v
		}
	}

	return funcMap
}

func durationToMilliseconds(t time.Duration) string {
	return fmt.Sprintf("%d ms", t.Nanoseconds()/int64(time.Millisecond))
}

func padString(s string) string {
	return fmt.Sprintf("%-25s:", s)
}

func floatToString(f float64) string {
	return fmt.Sprintf("%.2f", f)
}

func valueToString(v interface{}) string {
	return fmt.Sprintf("%v", v)
}

func boolToYesNo(b bool) string {
	if b {
		return "Y"
	}
	return "N"
}
