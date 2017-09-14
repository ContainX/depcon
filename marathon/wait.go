package marathon

import (
	"github.com/ContainX/depcon/pkg/logger"
	"github.com/ContainX/depcon/utils"
	"time"
)

var logWait = logger.GetLogger("depcon.deploy.wait")

func (c *MarathonClient) WaitForApplication(id string, timeout time.Duration) error {
	t_now := time.Now()
	t_stop := t_now.Add(timeout)

	c.logWaitApplication(id)
	for {
		if time.Now().After(t_stop) {
			return ErrorTimeout
		}

		app, err := c.GetApplication(id)
		if err == nil {
			if app.DeploymentID == nil || len(app.DeploymentID) <= 0 {
				logWait.Infof("Application deployment has completed for %s, elapsed time %s", id, utils.ElapsedStr(time.Since(t_now)))
				if app.HealthChecks != nil && len(app.HealthChecks) > 0 {
					err := c.WaitForApplicationHealthy(id, timeout)
					if err != nil {
						logWait.Errorf("Error waiting for application '%s' to become healthy: %s", id, err.Error())
					}
				} else {
					logWait.Warningf("No health checks defined for '%s', skipping waiting for healthy state", id)
				}
				return nil
			}
		}
		c.logWaitApplication(id)
		time.Sleep(time.Duration(2) * time.Second)
	}
}

func (c *MarathonClient) WaitForApplicationHealthy(id string, timeout time.Duration) error {
	t_now := time.Now()
	t_stop := t_now.Add(timeout)
	duration := time.Duration(2) * time.Second
	for {
		if time.Now().After(t_stop) {
			return ErrorTimeout
		}
		app, err := c.GetApplication(id)
		if err != nil {
			return err
		}
		total := app.TasksStaged + app.TasksRunning
		diff := total - app.TasksHealthy
		if diff == 0 {
			logWait.Infof("%v of %v expected instances are healthy.  Elapsed health check time of %s", app.TasksHealthy, total, utils.ElapsedStr(time.Since(t_now)))
			return nil
		}
		logWait.Infof("%v healthy instances.  Waiting for %v total instances. Retrying check in %v seconds", app.TasksHealthy, total, duration)
		time.Sleep(duration)
	}
}

func (c *MarathonClient) WaitForDeployment(id string, timeout time.Duration) error {

	t_now := time.Now()
	t_stop := t_now.Add(timeout)

	c.logWaitDeployment(id)

	for {
		if time.Now().After(t_stop) {
			return ErrorTimeout
		}
		if found, _ := c.HasDeployment(id); !found {
			c.logOutput(logWait.Infof, "Deployment has completed for %s, elapsed time %s", id, utils.ElapsedStr(time.Since(t_now)))
			return nil
		}
		c.logWaitDeployment(id)
		time.Sleep(time.Duration(2) * time.Second)
	}
}

func (c *MarathonClient) logWaitDeployment(id string) {
	c.logOutput(logWait.Infof, "Waiting for deployment %s", id)
}

func (c *MarathonClient) logWaitApplication(id string) {
	c.logOutput(logWait.Infof, "Waiting for application deployment to complete for %s", id)
}
