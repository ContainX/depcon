package marathon

import (
	"testing"
)

var c MarathonClient = MarathonClient{}

func TestCreateApplicationFromFile(t *testing.T) {

	envParams := make(map[string]string, 1)
	envParams["NODE_EXPORTER_VERSION"] = "1"

	opts := &CreateOptions{Wait: false, Force: true, ErrorOnMissingParams: true, EnvParams: envParams}
	app, _, err := c.ParseApplicationFromFile("resources/marathon.json", opts)
	if err != nil {
		log.Panic("Expected success %v", err)
	}
	log.Debug("%v", app)
	if app.Labels["tags"] != "prom-metrics" {
		log.Panic("Expected Labels parsed correctly")
	}
	if app.Instances != int(2) {
		log.Panic("Expected Instances parsed correctly")
	}

}
