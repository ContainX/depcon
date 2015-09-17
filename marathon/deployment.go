package marathon

import (
	"github.com/gondor/depcon/pkg/httpclient"
	"errors"
	"fmt"
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

func (c *MarathonClient) DeleteDeployment(id string) (*DeploymentID, error) {
	deploymentID := new(DeploymentID)
	resp := c.http.HttpDelete(c.marathonUrl(API_DEPLOYMENTS, id), nil, deploymentID)
	if resp.Error != nil {
		if resp.Error == httpclient.ErrorNotFound {
			return nil, errors.New(fmt.Sprintf("Deployment '%s' was not found", id))
		}
		return nil, resp.Error
	}
	return deploymentID, nil

}
