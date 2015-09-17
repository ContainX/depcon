package encoding

import (
	"encoding/json"
	"io"
	"strings"
)

// An encoder that marshal's and unmarshal's Json which implements the Encoder interface
type JSONEncoder struct{}

func newJSONEncoder() *JSONEncoder {
	return &JSONEncoder{}
}

func (e *JSONEncoder) MarshalIndent(data interface{}) (string, error) {
	if response, err := json.MarshalIndent(data, "", "   "); err != nil {
		return "", err
	} else {
		return string(response), err
	}
}

func (e *JSONEncoder) Marshal(data interface{}) (string, error) {
	if response, err := json.Marshal(data); err != nil {
		return "", err
	} else {
		return string(response), err
	}
}

func (e *JSONEncoder) UnMarshal(r io.Reader, result interface{}) error {
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(result); err != nil {
		return err
	}
	return nil
}

func (e *JSONEncoder) UnMarshalStr(data string, result interface{}) error {
	return e.UnMarshal(strings.NewReader(data), result)
}
