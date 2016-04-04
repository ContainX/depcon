package marathon

import (
	"errors"
	"fmt"
	"github.com/gondor/depcon/pkg/encoding"
	"github.com/gondor/depcon/pkg/envsubst"
	"github.com/gondor/depcon/pkg/httpclient"
	"github.com/gondor/depcon/utils"
	"os"
	"time"
)

const (
	ActionRestart  = "restart"
	ActionVersions = "versions"
	PathTasks      = "tasks"
)

var (
	ErrorAppExists        = errors.New("The application already exists")
	ErrorGroupExists      = errors.New("The group already exists")
	ErrorInvalidAppId     = errors.New("The application identifier is invalid")
	ErrorInvalidGroupId   = errors.New("The group identifier is invalid")
	ErrorNoAppExists      = errors.New("The application does not exist.  Create an application before updating")
	ErrorGropAppExists    = errors.New("The group does not exist.  Create a group before updating")
	ErrorAppParamsMissing = errors.New("One or more ${PARAMS} that were defined in the app configuration could not be resolved.")
)

func (c *MarathonClient) CreateApplicationFromFile(filename string, opts *CreateOptions) (*Application, error) {
	app, options, err := c.ParseApplicationFromFile(filename, opts)
	if err != nil {
		return app, err
	}
	return c.CreateApplication(app, options.Wait, options.Force)
}

func (c *MarathonClient) ParseApplicationFromFile(filename string, opts *CreateOptions) (*Application, *CreateOptions, error) {
	log.Info("Creating Application from file: %s", filename)
	options := initCreateOptions(opts)

	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("Error opening filename %s, %s", filename, err.Error())
	}

	var encoder encoding.Encoder
	encoder, err = encoding.NewEncoderFromFileExt(filename)
	if err != nil {
		return nil, nil, err
	}

	parsed, missing := envsubst.SubstFileTokens(file, filename, options.EnvParams)

	if options.ErrorOnMissingParams && missing {
		return nil, nil, ErrorAppParamsMissing
	}

	app := new(Application)
	err = encoder.UnMarshalStr(parsed, &app)
	if err != nil {
		return nil, nil, err
	}
	return app, options, nil
}

func (c *MarathonClient) CreateApplication(app *Application, wait, force bool) (*Application, error) {
	log.Info("Creating Application '%s', wait: %v, force: %v", app.ID, wait, force)

	result := new(Application)
	resp := c.http.HttpPost(c.marathonUrl(API_APPS), app, result)
	if resp.Error != nil {
		if resp.Error == httpclient.ErrorMessage {
			if resp.Status == 409 {
				if force {
					return c.UpdateApplication(app, wait)
				}
				return nil, ErrorAppExists
			}
			if resp.Status == 422 {
				return nil, ErrorInvalidAppId
			}
			return nil, fmt.Errorf("Error occurred (Status %v) Body -> %s", resp.Status, resp.Content)
		}
		return nil, resp.Error
	}
	if wait {
		err := c.WaitForApplication(result.ID, c.determineTimeout(app))
		if err != nil {
			return result, err
		}
	}
	app, err := c.GetApplication(result.ID)
	if err == nil {
		return app, nil
	}
	return result, nil
}

func (c *MarathonClient) UpdateApplication(app *Application, wait bool) (*Application, error) {
	log.Info("Update Application '%s', wait = %v", app.ID, wait)
	result := new(DeploymentID)
	id := utils.TrimRootPath(app.ID)
	app.ID = ""
	resp := c.http.HttpPut(c.marathonUrl(API_APPS, id), app, result)

	if resp.Error != nil {
		if resp.Error == httpclient.ErrorMessage {
			if resp.Status == 422 {
				return nil, ErrorNoAppExists
			}
		}
		return nil, resp.Error
	}
	if wait {
		if err := c.WaitForDeployment(result.DeploymentID, c.determineTimeout(app)); err != nil {
			return nil, err
		}
		err := c.WaitForApplication(id, c.determineTimeout(app))
		if err != nil {
			return nil, err
		}
	}
	// Get the latest version of the application to return
	app, err := c.GetApplication(id)
	return app, err
}

func (c *MarathonClient) ListApplications() (*Applications, error) {
	log.Debug("Enter: ListApplications")

	apps := new(Applications)

	resp := c.http.HttpGet(c.marathonUrl(API_APPS), apps)
	if resp.Error != nil {
		return nil, resp.Error
	}
	return apps, nil
}

func (c *MarathonClient) GetApplication(id string) (*Application, error) {
	log.Debug("Enter: GetApplication: %s", id)
	app := new(AppById)
	resp := c.http.HttpGet(c.marathonUrl(API_APPS, id), app)
	if resp.Error != nil {
		return nil, resp.Error
	}
	return &app.App, nil
}

func (c *MarathonClient) HasApplication(id string) (bool, error) {
	app, err := c.GetApplication(id)

	if err != nil {
		if err == httpclient.ErrorNotFound {
			return false, nil
		}
		return false, err
	}
	return app != nil, nil
}

func (c *MarathonClient) DestroyApplication(id string) (*DeploymentID, error) {
	log.Info("Deleting Application '%s'", id)
	deploymentId := new(DeploymentID)

	resp := c.http.HttpDelete(c.marathonUrl(API_APPS, id), nil, deploymentId)
	if resp.Error != nil {
		return nil, resp.Error
	}
	return deploymentId, nil
}

func (c *MarathonClient) RestartApplication(id string, force bool) (*DeploymentID, error) {
	log.Info("Restarting Application '%s', force: %v", id, force)

	deploymentId := new(DeploymentID)

	uri := fmt.Sprintf("%s?force=%v", c.marathonUrl(API_APPS, id, ActionRestart), force)
	resp := c.http.HttpPost(uri, nil, deploymentId)
	if resp.Error != nil {
		return nil, resp.Error
	}
	return deploymentId, nil
}

func (c *MarathonClient) ScaleApplication(id string, instances int) (*DeploymentID, error) {
	log.Info("Scale Application '%s' to %v instances", id, instances)

	update := new(Application)
	update.ID = id
	update.Instances = instances
	deploymentID := new(DeploymentID)
	resp := c.http.HttpPut(c.marathonUrl(API_APPS, id), &update, deploymentID)
	if resp.Error != nil {
		return nil, resp.Error
	}
	return deploymentID, nil
}

func (c *MarathonClient) ListVersions(id string) (*Versions, error) {
	versions := new(Versions)
	resp := c.http.HttpGet(c.marathonUrl(API_APPS, id, ActionVersions), versions)
	if resp.Error != nil {
		return nil, resp.Error
	}
	return versions, nil

}

func NewApplication(id string) *Application {
	application := new(Application)
	application.ID = id
	return application
}

// The number of instances that the application should run
// {count} - number of instances
func (app *Application) Count(count int) *Application {
	app.Instances = count
	return app
}

// The amount of memory in MB to assign per instance
// {memory} - memory in MB
func (app *Application) Memory(memory float64) *Application {
	app.Mem = memory
	return app
}

// The amount of CPU shares to assign per instance
// {cpu} - the CPU shares
func (app *Application) CPU(cpu float64) *Application {
	app.CPUs = cpu
	return app
}

// Rolls back an application to a specific version
// {version} - the version to rollback
func (app *Application) RollbackVersion(version string) *Application {
	app.Version = version
	return app
}

func (c *MarathonClient) determineTimeout(app *Application) time.Duration {
	if c.opts != nil && c.opts.WaitTimeout > 0 {
		return c.opts.WaitTimeout
	}

	if app == nil {
		return DefaultTimeout
	}

	max := DefaultTimeout

	if len(app.HealthChecks) > 0 {
		for _, h := range app.HealthChecks {
			grace := time.Duration(h.GracePeriodSeconds) * time.Second
			if grace > max {
				max = grace
			}
		}
		log.Debug("determineTimeout:  Max is %d\n", max)
		return max
	}
	return DefaultTimeout
}
