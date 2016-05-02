package bluegreen

import (
	"fmt"
	"github.com/ContainX/depcon/marathon"
	"strconv"
	"strings"
)

func deployStartTimeCompare(existingApp, currentApp *marathon.Application) int {
	if labelExists(existingApp, DeployStartedAt) && labelExists(currentApp, DeployStartedAt) {
		return strings.Compare(existingApp.Labels[DeployStartedAt], currentApp.Labels[DeployStartedAt])
	}
	return 0
}

func findServicePort(app *marathon.Application) int {
	if app.Container != nil && len(app.Container.Docker.PortMappings) > 0 {
		return app.Container.Docker.PortMappings[0].ServicePort
	}
	if len(app.Ports) > 0 {
		return app.Ports[0]
	}
	return 0
}

func labelExists(app *marathon.Application, label string) bool {
	_, exist := app.Labels[label]
	return exist
}

func formatIdentifier(appId, colour string) string {
	id := fmt.Sprintf("%s-%s", appId, colour)
	if []rune(id)[0] != '/' {
		id = "/" + id
	}
	return id
}

func intOrZero(s string) int {
	if v, err := strconv.Atoi(s); err != nil {
		return 0
	} else {
		return v
	}
}
