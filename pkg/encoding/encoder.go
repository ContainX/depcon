// YAML and JSON encoding
package encoding

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type EncoderType int

const (
	JSON EncoderType = 1 + iota
	YAML
)

var ErrorInvalidExtension = errors.New("File extension must be [.json | .yml | .yaml]")

type Encoder interface {
	MarshalIndent(data interface{}) (string, error)

	Marshal(data interface{}) (string, error)

	UnMarshal(r io.Reader, result interface{}) error

	UnMarshalStr(data string, result interface{}) error
}

func NewEncoder(encoder EncoderType) (Encoder, error) {
	switch encoder {
	case JSON:
		return newJSONEncoder(), nil
	case YAML:
		return newYAMLEncoder(), nil
	default:
		panic(fmt.Errorf("Unsupported encoder type"))
	}
}

func NewEncoderFromFileExt(filename string) (Encoder, error) {

	if et, err := EncoderTypeFromExt(filename); err != nil {
		return nil, err
	} else {
		return NewEncoder(et)
	}
}

func EncoderTypeFromExt(filename string) (EncoderType, error) {
	switch filepath.Ext(filename) {
	case ".yml", ".yaml":
		return YAML, nil
	case ".json":
		return JSON, nil
	}
	return JSON, ErrorInvalidExtension

}

func ConvertFile(infile, outfile string, dataType interface{}) error {
	var fromEnc, toEnc Encoder
	var encErr error

	if fromEnc, encErr = NewEncoderFromFileExt(infile); encErr != nil {
		return encErr
	}

	if toEnc, encErr = NewEncoderFromFileExt(outfile); encErr != nil {
		return encErr
	}

	file, err := os.Open(infile)
	if err != nil {
		return err
	}

	if err := fromEnc.UnMarshal(file, dataType); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(outfile), 0700); err != nil {
		return err
	}
	f, err := os.OpenFile(outfile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	if data, err := toEnc.MarshalIndent(dataType); err != nil {
		return err
	} else {
		f.WriteString(data)
	}
	return nil
}
