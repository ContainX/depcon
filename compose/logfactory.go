package compose

import (
	depconLog "github.com/ContainX/depcon/pkg/logger"
	"github.com/docker/libcompose/logger"
	"io"
	"os"
)

var log = depconLog.GetLogger("depcon.compose")

type ComposeLogger struct {
}

func (n *ComposeLogger) Out(b []byte) {
	log.Infof("%v", b)
}

func (n *ComposeLogger) Err(b []byte) {
	log.Errorf("%v", b)
}

func (n *ComposeLogger) ErrWriter() io.Writer {
	return os.Stderr
}

func (n *ComposeLogger) OutWriter() io.Writer {
	return os.Stdout
}

func (n *ComposeLogger) Create(name string) logger.Logger {
	return &ComposeLogger{}
}
