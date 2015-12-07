package compose

import (
	"github.com/gondor/depcon/pkg/logger"
)

const (
	DEFAULT_PROJECT string = "depcon_proj"
)

var log = logger.GetLogger("depcon.compose")

type Compose interface {
	Up(services ...string) error

	Kill(services ...string) error

	Logs(services ...string) error

	Delete(services ...string) error

	Build(services ...string) error

	Restart(services ...string) error

	Pull(services ...string) error

	Start(services ...string) error

	Stop(services ...string) error

	Port(index int, proto, service, port string) error

	PS(quiet bool) error
}
