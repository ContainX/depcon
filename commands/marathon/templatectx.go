package marathon

import (
	"io"
	"io/ioutil"
	"os"
	"text/template"

	"fmt"
	"github.com/ContainX/depcon/pkg/encoding"
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
)

const (
	DefaultEnv      = "-"
	DefaultRootPath = "."
	ContextErrFmt   = "Error parsing template context: %s - %s"
)

// Template based Functions
var Funcs = FuncMap()

type TemplateContext struct {
	Environments map[string]*TemplateEnvironment `json:"environments,omitempty"`
}

type TemplateEnvironment struct {
	Apps map[string]map[string]interface{} `json:"apps,omitempty"`
}

func (ctx *TemplateContext) Transform(writer io.Writer, descriptor, rootDir string) error {
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

		if rootDir == "" {
			rootDir = DefaultRootPath
		}

		if matches, err := filepath.Glob(fmt.Sprintf("%s/**/*.tmpl", rootDir)); err == nil && len(matches) > 0 {
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
	// Return empty context if non-exists
	if !TemplateExists(filename) {
		fmt.Println("Ignoring context file (not found) - template-context.json")
		return &TemplateContext{Environments: make(map[string]*TemplateEnvironment)}, nil
	}

	ctx, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	encoder, err := encoding.NewEncoder(encoding.JSON)
	if err != nil {
		return nil, fmt.Errorf(ContextErrFmt, filename, err.Error())
	}

	result := &TemplateContext{Environments: make(map[string]*TemplateEnvironment)}

	if err := encoder.UnMarshal(ctx, result); err != nil {
		return nil, fmt.Errorf(ContextErrFmt, filename, err.Error())
	}
	return result, nil
}
