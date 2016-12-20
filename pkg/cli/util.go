package cli

import (
	"strings"
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
	return t.Local().Format("2006-01-02 15:04:05")
}

func NameValueSliceToMap(params []string) map[string]string {
	if params == nil {
		return nil
	}
	envmap := make(map[string]string)
	for _, p := range params {
		if strings.Contains(p, "=") {
			v := strings.Split(p, "=")
			envmap[v[0]] = v[1]
		}
	}
	return envmap
}
