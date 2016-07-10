package marathon

import (
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"text/template"

	"github.com/ContainX/depcon/pkg/encoding"
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
)

const (
	DefaultEnv = "-"
)

// Templated based Functions
var Funcs = template.FuncMap{
	"default": func(args ...interface{}) interface{} {
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
	},
	"isEnv": func(value string) bool {
		if len(value) > 0 {
			current := strings.ToLower(viper.GetString(ENV_NAME))
			return current == strings.ToLower(value)
		}
		return false
	},
}

type TemplateContext struct {
	Environments map[string]*TemplateEnvironment `json:"environments,omitempty"`
}

type TemplateEnvironment struct {
	Apps map[string]map[string]interface{} `json:"apps,omitempty"`
}

func (ctx *TemplateContext) Transform(writer io.Writer, descriptor string) error {
	var t *template.Template

	if b, err := ioutil.ReadFile(descriptor); err != nil {
		return err
	} else {
		var e error
		t = template.New(descriptor).Funcs(Funcs)
		t, e = t.Parse(string(b))
		if e != nil {
			return e
		}
		if matches, err := filepath.Glob("./**/*.tmpl"); err == nil && len(matches) > 0 {
			if t, e = t.ParseFiles(matches...); err != nil {
				return err
			}
		}
	}
	environment := viper.GetString(ENV_NAME)
	m := ctx.mergeAppWithDefault(strings.ToLower(environment))

	if err := t.Execute(writer, m); err != nil {
		return err
	}
	return nil
}

// Validates the specified app is declared within the current envirnoment. If it is any values missing from
// specific environment are propagated from the default environment
func (ctx *TemplateContext) mergeAppWithDefault(env string) map[string]map[string]interface{} {
	if _, exists := ctx.Environments[DefaultEnv]; !exists {
		if _, cok := ctx.Environments[env]; !cok {
			return make(map[string]map[string]interface{})
		}
		return ctx.Environments[env].Apps
	}

	defm := ctx.Environments[DefaultEnv].Apps

	if ctx.Environments[env] == nil {
		return defm
	}

	envm := ctx.Environments[env].Apps

	merged := make(map[string]map[string]interface{})

	for app, props := range envm {
		merged[app] = props
	}

	for app, props := range defm {
		if _, exists := merged[app]; !exists {
			merged[app] = props
		} else {
			for k, v := range props {
				if _, ok := merged[app][k]; !ok {
					merged[app][k] = v
				}
			}

		}
	}
	return merged
}

func recovery() {
	recover()
}

func TemplateExists(filename string) bool {

	if len(filename) > 0 {
		if _, err := os.Stat(filename); err == nil {
			return true
		}
	}
	return false
}

func LoadTemplateContext(filename string) (*TemplateContext, error) {
	ctx, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	encoder, err := encoding.NewEncoder(encoding.JSON)
	if err != nil {
		return nil, err
	}

	result := &TemplateContext{Environments: make(map[string]*TemplateEnvironment)}

	if err := encoder.UnMarshal(ctx, result); err != nil {
		return nil, err
	}
	return result, nil
}
