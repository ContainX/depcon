package bluegreen

import (
	"errors"
	"fmt"
	"github.com/gondor/depcon/marathon"
	"github.com/gondor/depcon/pkg/logger"
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

func (c *BGClient) DeployBlueGreen(app *marathon.Application) error {

	if app.Labels == nil || len(app.Labels) == 0 {
		return ErrorNoLabels
	}

	if !labelExists(app, DeployGroup) {
		return fmt.Errorf(LabelFormatErr, DeployGroup)
	}

	if !labelExists(app, DeployGroupAltPort) {
		return fmt.Errorf(LabelFormatErr, DeployGroupAltPort)
	}

	group := app.Labels[DeployGroup]
	groupAltPort, err := strconv.Atoi(app.Labels[DeployGroupAltPort])
	if err != nil {
		return err
	}

	app.Labels[ProxyAppId] = app.ID
	servicePort := findServicePort(app)

	if servicePort <= 0 {
		return ErrorNoServicePortSet
	}

	state, err := c.bgAppInfo(group, groupAltPort)
	if err != nil {
		return err
	}

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

	return nil
}

func (c *BGClient) startDeployment(app *marathon.Application, state appState) bool {
	if !state.resuming {
	}
	if state.existingApp != nil {
		return c.checkIfTasksDrained(app, state.existingApp, time.Now())
	}
	return false
}

func (c *BGClient) bgAppInfo(deployGroup string, deployGroupAltPort int) (appState, error) {
	apps, err := c.marathon.ListApplications()

	if err != nil {
		return err
	}

	var existingApp marathon.Application

	colour := ColourBlue
	nextPort := deployGroupAltPort
	resume := false

	for _, app := range apps.Apps {
		if len(app.Labels) <= 0 {
			continue
		}
		if labelExists(app, DeployGroup) && labelExists(app, DeployGroupColour) && app.Labels[DeployGroup] == deployGroupAltPort {
			if existingApp != nil {
				if c.opts.Resume {
					log.Info("Found previous deployment -- resuming")
					resume = true
					if deployStartTimeCompare(existingApp, app) == -1 {
						break
					}
				} else {
					return errors.New("There appears to be an existing deployment in progress")
				}
			}
			prev_colour := app.Labels[DeployGroupColour]
			prev_port := app.Ports[0]
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
	return &appState{
		existingApp: existingApp,
		nextPort:    nextPort,
		colour:      colour,
		resuming:    resume,
	}, nil
}
