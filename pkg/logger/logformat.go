// +build !windows

package logger

///
// Color overrides and formatting changes for +nix environments
//

import (
	"github.com/op/go-logging"
	"os"
)

var format = logging.MustStringFormatter(
	"%{color}%{time:2006-01-02 15:04:05} %{level:.7s} [%{module}]:%{color:reset} %{message}",
)

func init() {
	backend := logging.NewLogBackend(os.Stdout, "", 0)
	backendFmt := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backendFmt)
}
