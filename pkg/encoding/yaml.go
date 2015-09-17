package encoding

import (
	"github.com/ghodss/yaml"
	"io"
	"io/ioutil"
	"strings"
)

// An encoder that marshal's and unmarshal's YAML which implements the Encoder interface
type YAMLEncoder struct{}

func newYAMLEncoder() *YAMLEncoder {
	return &YAMLEncoder{}
}

func (e *YAMLEncoder) MarshalIndent(data interface{}) (string, error) {
	return e.Marshal(data)
}

func (e *YAMLEncoder) Marshal(data interface{}) (string, error) {

	if response, err := yaml.Marshal(data); err != nil {
		return "", err
	} else {
		return string(response), err
	}
}

func (e *YAMLEncoder) UnMarshal(r io.Reader, result interface{}) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(b, result); err != nil {
		return err
	}
	return nil
}

func (e *YAMLEncoder) UnMarshalStr(data string, result interface{}) error {
	return e.UnMarshal(strings.NewReader(data), result)
}
