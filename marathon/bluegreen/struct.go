package bluegreen

import (
	"github.com/gondor/depcon/marathon"
	"github.com/gondor/depcon/pkg/httpclient"
	"time"
)

type BlueGreenOptions struct {
	// The max time to wait on HAProxy to drain connections (in seconds)
	ProxyWaitTimeout time.Duration
	// Initial number of app instances to create
	InitialInstances int
	// Delay (in seconds) to wait between each successive deployment step
	StepDelay time.Duration
	// Resume from previous deployment
	Resume bool
	// Marathon-LB stats endpoint - ex: http://host:9090
	LoadBalancer string
}

type BGClient struct {
	marathon marathon.Marathon
	opts     BlueGreenOptions
	http     httpclient.HttpClient
}

type appState struct {
	colour      string
	nextPort    int
	existingApp marathon.Application
	resuming    bool
}

type proxyInfo struct {
	hmap          map[string]int
	backends      []string
	instanceCount int
}
