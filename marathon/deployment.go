package marathon

import (
	"errors"
	"fmt"
	"github.com/ContainX/depcon/pkg/httpclient"
	"strings"
)

func (c *MarathonClient) ListDeployments() ([]*Deploy, error) {
	var deploys []*Deploy
	resp := c.http.HttpGet(c.marathonUrl(API_DEPLOYMENTS), &deploys)
	if resp.Error != nil {
		return nil, resp.Error
	}
	return deploys, nil
}

func (c *MarathonClient) HasDeployment(id string) (bool, error) {
	deployments, err := c.ListDeployments()
	if err != nil {
		return false, err
	}
	for _, deployment := range deployments {
		if deployment.DeployID == id {
			return true, nil
		}
	}
	return false, nil
}

func (c *MarathonClient) DeleteDeployment(id string, force bool) (*DeploymentID, error) {
	deploymentID := new(DeploymentID)
	uri := fmt.Sprintf("%s?force=%v", c.marathonUrl(API_DEPLOYMENTS, id), force)
	resp := c.http.HttpDelete(uri, nil, deploymentID)
	if resp.Error != nil {
		if resp.Error == httpclient.ErrorNotFound {
			return nil, errors.New(fmt.Sprintf("Deployment '%s' was not found", id))
		}
		return nil, resp.Error
	}
	return deploymentID, nil
}

func (c *MarathonClient) CancelAppDeployment(appId string, matchPrefix bool) (*DeploymentID, error) {
	if deployments, err := c.ListDeployments(); err == nil {
		for _, value := range deployments {
			for _, id := range value.AffectedApps {
				if doesIDMatch(appId, id, matchPrefix) {
					log.Info("Removing matched deployment: %s for app: %s", value.DeployID, id)
					return c.DeleteDeployment(value.DeployID, true)
				}
			}
		}
	} else {
		return nil, err
	}
	return nil, nil
}

func doesIDMatch(appId, otherId string, matchPrefix bool) bool {
	if matchPrefix {
		return strings.HasPrefix(otherId, appId)
	}
	return appId == otherId
}
