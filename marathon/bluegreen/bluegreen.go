package bluegreen

import (
	"github.com/gondor/depcon/marathon"
	"github.com/gondor/depcon/pkg/httpclient"
	"time"
)

type BlueGreen interface {

	// Starts a blue green deployment.  If the application exists then the deployment will slowly
	// release the new version, draining connections from the HAProxy balancer during the process
	// {filename} - the file name of the json | yaml application
	// {opts} - blue/green options
	DeployBlueGreenFromFile(filename string) (*marathon.Application, error)

	// Starts a blue green deployment.  If the application exists then the deployment will slowly
	// release the new version, draining connections from the HAProxy balancer during the process
	// {app} - the application to deploy/update
	// {opts} - blue/green options
	DeployBlueGreen(app *marathon.Application) (*marathon.Application, error)
}

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
	// if true will attempt to wait until the NEW application or group is running
	Wait bool
	// If true an error will be returned on params defined in the configuration file that
	// could not resolve to user input and environment variables
	ErrorOnMissingParams bool
	// Additional environment params - looks at this map for token substitution which takes
	// priority over matching environment variables
	EnvParams map[string]string
}

type BGClient struct {
	marathon marathon.Marathon
	opts     *BlueGreenOptions
	http     *httpclient.HttpClient
}

type appState struct {
	colour      string
	nextPort    int
	existingApp *marathon.Application
	resuming    bool
}

type proxyInfo struct {
	hmap          map[string]int
	backends      [][]string
	instanceCount int
}

func NewBlueGreenClient(marathon marathon.Marathon, opts *BlueGreenOptions) BlueGreen {
	c := new(BGClient)
	c.marathon = marathon
	c.opts = opts
	c.http = httpclient.DefaultHttpClient()
	return c
}

func NewBlueGreenOptions() *BlueGreenOptions {
	opts := &BlueGreenOptions{}
	opts.InitialInstances = 1
	opts.ProxyWaitTimeout = time.Duration(300) * time.Second
	opts.StepDelay = time.Duration(6) * time.Second
	return opts
}
