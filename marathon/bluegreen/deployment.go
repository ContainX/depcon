package bluegreen

import (
	"errors"
	"fmt"
	"github.com/ContainX/depcon/marathon"
	"github.com/ContainX/depcon/pkg/logger"
	"os"
	"strconv"
	"time"
)

const (
	DeployGroup           = "HAPROXY_DEPLOYMENT_GROUP"
	DeployGroupAltPort    = "HAPROXY_DEPLOYMENT_ALT_PORT"
	DeployGroupColour     = "HAPROXY_DEPLOYMENT_COLOUR"
	DeployProxyPort       = "HAPROXY_0_PORT"
	DeployTargetInstances = "HAPROXY_DEPLOYMENT_TARGET_INSTANCES"
	DeployStartedAt       = "HAPROXY_DEPLOYMENT_STARTED_AT"
	ProxyAppId            = "HAPROXY_APP_ID"
	ColourBlue            = "blue"
	ColourGreen           = "green"
)

var (
	ErrorNoLabels         = errors.New("No labels found. Please define the HAPROXY_DEPLOYMENT_GROUP and HAPROXY_DEPLOYMENT_ALT_PORT label")
	ErrorNoServicePortSet = errors.New("No service port set")
	LabelFormatErr        = "Please define the %s label"
	log                   = logger.GetLogger("depcon.marathon.bg")
)

func (c *BGClient) DeployBlueGreenFromFile(filename string) (*marathon.Application, error) {

	log.Debug("Enter DeployBlueGreenFromFile")

	parseOpts := &marathon.CreateOptions{
		ErrorOnMissingParams: c.opts.ErrorOnMissingParams,
		EnvParams:            c.opts.EnvParams,
	}
	app, _, err := c.marathon.ParseApplicationFromFile(filename, parseOpts)
	if err != nil {
		return nil, err
	}
	return c.DeployBlueGreen(app)
}

func (c *BGClient) DeployBlueGreen(app *marathon.Application) (*marathon.Application, error) {

	log.Debug("Enter DeployBlueGreen")

	// Before we return the client lets make sure the LoadBalancer is properly defined
	c.isProxyAlive()

	if app.Labels == nil || len(app.Labels) == 0 {
		return nil, ErrorNoLabels
	}

	if !labelExists(app, DeployGroup) {
		return nil, fmt.Errorf(LabelFormatErr, DeployGroup)
	}

	if !labelExists(app, DeployGroupAltPort) {
		return nil, fmt.Errorf(LabelFormatErr, DeployGroupAltPort)
	}

	group := app.Labels[DeployGroup]
	groupAltPort, err := strconv.Atoi(app.Labels[DeployGroupAltPort])
	if err != nil {
		return nil, err
	}

	app.Labels[ProxyAppId] = app.ID
	servicePort := findServicePort(app)

	if servicePort <= 0 {
		return nil, ErrorNoServicePortSet
	}

	state, err := c.bgAppInfo(group, groupAltPort)
	if err != nil {
		return nil, err
	}

	app = c.updateServicePort(app, state.nextPort)

	app.ID = formatIdentifier(app.ID, state.colour)

	if state.existingApp != nil {
		app.Instances = c.opts.InitialInstances
		app.Labels[DeployTargetInstances] = strconv.Itoa(state.existingApp.Instances)
	} else {
		app.Labels[DeployTargetInstances] = strconv.Itoa(app.Instances)
	}

	app.Labels[DeployGroupColour] = state.colour
	app.Labels[DeployStartedAt] = time.Now().Format(time.RFC3339)
	app.Labels[DeployProxyPort] = strconv.Itoa(servicePort)

	if c.opts.DryRun {
		return app, nil
	}

	c.startDeployment(app, state)

	return c.marathon.GetApplication(app.ID)
}

func (c *BGClient) updateServicePort(app *marathon.Application, port int) *marathon.Application {
	log.Debug("Entering updateServicePort, port=%d", port)
	if app.Container != nil && app.Container.Docker != nil {
		if app.Container.Docker.PortMappings != nil && len(app.Container.Docker.PortMappings) > 0 {
			app.Container.Docker.PortMappings[0].ServicePort = port
		}
		if len(app.Ports) > 0 {
			app.Ports[0] = port
		}
	}
	return app
}

func (c *BGClient) startDeployment(app *marathon.Application, state *appState) bool {
	log.Debug("startDeployment: resuming: %v", state.resuming)
	if !state.resuming {
		a, err := c.marathon.CreateApplication(app, true, false)
		if err != nil {
			log.Error("Unable to create application: %s", err.Error())
			os.Exit(1)
		}
		app = a
	}
	if state.existingApp != nil {
		return c.checkIfTasksDrained(app, state.existingApp, time.Now())
	}
	return false
}

func (c *BGClient) bgAppInfo(deployGroup string, deployGroupAltPort int) (*appState, error) {
	apps, err := c.marathon.ListApplications()

	if err != nil {
		return nil, err
	}

	var existingApp marathon.Application

	colour := ColourBlue
	nextPort := deployGroupAltPort
	resume := false
	exists := false

	for _, app := range apps.Apps {
		log.Debug("bgAppInfo: loop %s", app.ID)
		if len(app.Labels) <= 0 {
			continue
		}
		if labelExists(&app, DeployGroup) && labelExists(&app, DeployGroupColour) && app.Labels[DeployGroup] == deployGroup {
			if exists {
				if c.opts.Resume {
					log.Info("Found previous deployment -- resuming")
					resume = true
					if deployStartTimeCompare(&existingApp, &app) == -1 {
						break
					}
				} else {
					return nil, errors.New("There appears to be an existing deployment in progress")
				}
			}
			prev_colour := app.Labels[DeployGroupColour]
			prev_port := app.Ports[0]

			log.Debug("bgAppInfo: assigning %s to existing app: %s = %s", app.ID, app.Labels[DeployGroup], deployGroup)
			existingApp = app
			exists = true

			if prev_port == deployGroupAltPort {
				nextPort, _ = strconv.Atoi(app.Labels[DeployProxyPort])
			} else {
				nextPort = deployGroupAltPort
			}

			if prev_colour == ColourBlue {
				colour = ColourGreen
			} else {
				colour = ColourBlue
			}
		}
	}

	as := &appState{
		nextPort: nextPort,
		colour:   colour,
		resuming: resume,
	}

	if exists {
		as.existingApp = &existingApp
		log.Debug("bgAppInfo: Returning %s, np: %d, clr: %s", sprintApp(as.existingApp), as.nextPort, as.colour)
	} else {
		log.Debug("bgAppInfo: Returning np: %d, clr: %s", as.nextPort, as.colour)
	}
	return as, nil
}

func sprintApp(a *marathon.Application) string {
	return fmt.Sprintf("[id: %s, i: %d, lbls: %v]", a.ID, a.Instances, a.Labels)
}
