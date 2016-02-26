package compose

import (
	"github.com/docker/libcompose/logger"
	depconLog "github.com/gondor/depcon/pkg/logger"
)

var log = depconLog.GetLogger("depcon.compose")

type ComposeLogger struct {
}

func (n *ComposeLogger) Out(b []byte) {
	log.Info("%v", b)
}

func (n *ComposeLogger) Err(b []byte) {
	log.Error("%v", b)
}

func (n *ComposeLogger) Create(name string) logger.Logger {
	return &ComposeLogger{}
}
