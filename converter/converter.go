package converter

import (
	"github.com/gondor/depcon/marathon"
)

type Converter interface {

	// Converts a Docker Compose file to a Marathon Applications structure where
	// each compose service becomes a Marathon application
	//
	// {composeFile} - the docker compose file path/name
	ComposeToMarathon(composeFile string) (*marathon.Applications, error)

}


