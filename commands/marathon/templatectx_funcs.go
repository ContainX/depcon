package marathon

import (
	"fmt"
	"github.com/spf13/viper"
	"math"
	"os"
	"reflect"
	"strconv"
	"strings"
	"text/template"
)

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"default":    dfault,
		"dict":       dict,
		"isEnv":      isEnv,
		"isNotEnv":   isNotEnv,
		"empty":      empty,
		"max":        max,
		"min":        min,
		"toInt64":    toInt64,
		"trim":       strings.TrimSpace,
		"upper":      strings.ToUpper,
		"lower":      strings.ToLower,
		"title":      strings.Title,
		"trimAll":    func(a, b string) string { return strings.Trim(b, a) },
		"trimSuffix": func(a, b string) string { return strings.TrimSuffix(b, a) },
		"trimPrefix": func(a, b string) string { return strings.TrimPrefix(b, a) },
		"split":      split,
		"env":        func(s string) string { return os.Getenv(s) },
		"expandenv":  func(s string) string { return os.ExpandEnv(s) },
		"N":          N,
	}
}

func isEnv(value string) bool {
	if len(value) > 0 {
		current := strings.ToLower(viper.GetString(ENV_NAME))
		return current == strings.ToLower(value)
	}
	return false
}

func isNotEnv(value string) bool {
	return !isEnv(value)
}

func dfault(args ...interface{}) interface{} {
	arg := args[0]
	if len(args) < 2 {
		return arg
	}
	value := args[1]

	defer recovery()

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String, reflect.Slice, reflect.Array, reflect.Map:
		if v.Len() == 0 {
			return arg
		}
	case reflect.Bool:
		if !v.Bool() {
			return arg
		}
	default:
		return value
	}

	return value
}

// empty returns true if the given value has the zero value for its type.
func empty(given interface{}) bool {
	g := reflect.ValueOf(given)
	if !g.IsValid() {
		return true
	}

	// Basically adapted from text/template.isTrue
	switch g.Kind() {
	default:
		return g.IsNil()
	case reflect.Array, reflect.Slice, reflect.Map, reflect.String:
		return g.Len() == 0
	case reflect.Bool:
		return g.Bool() == false
	case reflect.Complex64, reflect.Complex128:
		return g.Complex() == 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return g.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return g.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return g.Float() == 0
	case reflect.Struct:
		return false
	}
}

func max(a interface{}, i ...interface{}) int64 {
	aa := toInt64(a)
	for _, b := range i {
		bb := toInt64(b)
		if bb > aa {
			aa = bb
		}
	}
	return aa
}

func min(a interface{}, i ...interface{}) int64 {
	aa := toInt64(a)
	for _, b := range i {
		bb := toInt64(b)
		if bb < aa {
			aa = bb
		}
	}
	return aa
}

// toInt64 converts integer types to 64-bit integers
func toInt64(v interface{}) int64 {

	if str, ok := v.(string); ok {
		iv, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return 0
		}
		return iv
	}

	val := reflect.Indirect(reflect.ValueOf(v))
	switch val.Kind() {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		return val.Int()
	case reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return int64(val.Uint())
	case reflect.Uint, reflect.Uint64:
		tv := val.Uint()
		if tv <= math.MaxInt64 {
			return int64(tv)
		}
		// TODO: What is the sensible thing to do here?
		return math.MaxInt64
	case reflect.Float32, reflect.Float64:
		return int64(val.Float())
	case reflect.Bool:
		if val.Bool() == true {
			return 1
		}
		return 0
	default:
		return 0
	}
}

func split(sep, orig string) map[string]string {
	parts := strings.Split(orig, sep)
	res := make(map[string]string, len(parts))
	for i, v := range parts {
		res["_"+strconv.Itoa(i)] = v
	}
	return res
}

func dict(v ...interface{}) map[string]interface{} {
	dict := map[string]interface{}{}
	lenv := len(v)
	for i := 0; i < lenv; i += 2 {
		key := strval(v[i])
		if i+1 >= lenv {
			dict[key] = ""
			continue
		}
		dict[key] = v[i+1]
	}
	return dict
}

func strval(v interface{}) string {
	switch v := v.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case error:
		return v.Error()
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}

func N(start, end int) (stream chan int) {
	stream = make(chan int)
	go func() {
		for i := start; i <= end; i++ {
			stream <- i
		}
		close(stream)
	}()
	return
}
