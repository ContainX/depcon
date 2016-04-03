// Common utilities
package utils

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	urlPrefix []string = []string{"http://", "https://"}
)

func TrimRootPath(id string) string {
	if strings.HasPrefix(id, "/") {
		return strings.TrimPrefix(id, "/")
	}
	return id
}

func Contains(elements []string, value string) bool {
	for _, element := range elements {
		if element == value {
			return true
		}
	}
	return false
}

func BuildPath(host string, elements []string) string {
	var buffer bytes.Buffer
	buffer.WriteString(TrimRootPath(host))
	for _, e := range elements {
		buffer.WriteString("/")
		buffer.WriteString(e)
	}
	return buffer.String()
}

func ConcatInts(iarr []int) string {
	var b bytes.Buffer
	for idx, i := range iarr {
		if idx > 0 {
			b.WriteString(" ,")
		}
		b.WriteString(strconv.Itoa(i))
	}
	return b.String()
}

/** Concats a string array of ids in the form of /id with an output string of id, id2, etc */
func ConcatIdentifiers(ids []string) string {
	if ids == nil {
		return ""
	}
	var b bytes.Buffer
	for idx, id := range ids {
		if idx > 0 {
			b.WriteString(" ,")
		}
		b.WriteString(TrimRootPath(id))
	}
	return b.String()
}

func HasURLScheme(url string) bool {
	return HasPrefix(url, urlPrefix...)
}

func HasPrefix(str string, prefixes ...string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(str, prefix) {
			return true
		}
	}
	return false
}

func ElapsedStr(d time.Duration) string {
	return fmt.Sprintf("%0.2f sec(s)", d.Seconds())
}

func IntInSlice(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func MapStringKeysToSlice(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
