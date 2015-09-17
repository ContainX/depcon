package cli

import (
	"time"
)

func EvalPrintUsage(usage_func func() error, args []string, minlen int) bool {
	if len(args) < minlen {
		usage_func()
		return true
	}
	return false
}

func FormatDate(date string) string {
	t, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return date
	}
	return t.Format("2006-01-_2 15:04:05")
}
