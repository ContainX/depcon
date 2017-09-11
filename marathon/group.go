package marathon

import (
	"fmt"
	"github.com/ContainX/depcon/pkg/encoding"
	"github.com/ContainX/depcon/pkg/envsubst"
	"github.com/ContainX/depcon/pkg/httpclient"
	"io"
	"os"
	"strings"
	"time"
)

func (c *MarathonClient) CreateGroupFromString(filename string, grpstr string, opts *CreateOptions) (*Group, error) {
	et, err := encoding.EncoderTypeFromExt(filename)
	if err != nil {
		return nil, err
	}
	group, err := c.ParseGroupFromString(strings.NewReader(grpstr), et, opts)

	if err != nil {
		return group, err
	}

	if opts.StopDeploy {
		if deployment, err := c.CancelAppDeployment(group.GroupID, true); err == nil && deployment != nil {
			c.logOutput(log.Info, "Previous deployment found..  cancelling and waiting until complete.")
			c.WaitForDeployment(deployment.DeploymentID, time.Second*30)
		}
	}

	return c.CreateGroup(group, opts.Wait, opts.Force)
}

func (c *MarathonClient) CreateGroupFromFile(filename string, opts *CreateOptions) (*Group, error) {
	log.Info("Creating Group from file: %s", filename)

	group, err := c.ParseGroupFromFile(filename, opts)
	if err != nil {
		return group, err
	}

	if opts.StopDeploy {
		if deployment, err := c.CancelAppDeployment(group.GroupID, true); err == nil && deployment != nil {
			log.Info("Previous deployment found..  cancelling and waiting until complete.")
			c.WaitForDeployment(deployment.DeploymentID, time.Second*30)
		}
	}

	return c.CreateGroup(group, opts.Wait, opts.Force)
}

func (c *MarathonClient) ParseGroupFromFile(filename string, opts *CreateOptions) (*Group, error) {
	log.Info("Creating Group from file: %s", filename)

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Error opening filename %s, %s", filename, err.Error())
	}

	if et, err := encoding.EncoderTypeFromExt(filename); err != nil {
		return nil, err
	} else {
		return c.ParseGroupFromString(file, et, opts)
	}
}

func (c *MarathonClient) ParseGroupFromString(r io.Reader, et encoding.EncoderType, opts *CreateOptions) (*Group, error) {

	options := initCreateOptions(opts)

	var encoder encoding.Encoder
	var err error

	encoder, err = encoding.NewEncoder(et)
	if err != nil {
		return nil, err
	}

	parsed, missing := envsubst.SubstFileTokens(r, options.EnvParams)

	if opts.ErrorOnMissingParams && missing {
		return nil, ErrorAppParamsMissing
	}

	if opts.DryRun {
		fmt.Printf("Create Group :: DryRun :: Template Output\n\n%s", parsed)
		os.Exit(0)
	}

	group := new(Group)
	err = encoder.UnMarshalStr(parsed, &group)
	if err != nil {
		return nil, err
	}
	return group, nil
}

func (c *MarathonClient) CreateGroup(group *Group, wait, force bool) (*Group, error) {
	c.logOutput(log.Info, "Creating Group '%s', wait: %v, force: %v", group.GroupID, wait, force)
	result := new(DeploymentID)
	resp := c.http.HttpPost(c.marathonUrl(API_GROUPS), group, result)
	if resp.Error != nil {
		if resp.Error == httpclient.ErrorMessage {
			if resp.Status == 409 {
				if force {
					return c.UpdateGroup(group, wait)
				}
				return nil, ErrorGroupExists
			}
			if resp.Status == 422 {
				return nil, ErrorInvalidGroupId
			}
			return nil, fmt.Errorf("Error occurred (Status %v) Body -> %s", resp.Status, resp.Content)
		}
		return nil, resp.Error
	}
	if wait {
		if err := c.WaitForDeployment(result.DeploymentID, time.Duration(500)*time.Second); err != nil {
			return nil, err
		}
	}
	return group, nil
}

func (c *MarathonClient) UpdateGroup(group *Group, wait bool) (*Group, error) {
	log.Info("Update Group '%s', wait = %v", group.GroupID, wait)
	result := new(DeploymentID)
	resp := c.http.HttpPut(c.marathonUrl(API_GROUPS), group, result)

	if resp.Error != nil {
		if resp.Error == httpclient.ErrorMessage {
			if resp.Status == 422 {
				return nil, ErrorGroupAppExists
			}
		}
		return nil, resp.Error
	}
	if wait {
		if err := c.WaitForDeployment(result.DeploymentID, c.determineTimeout(nil)); err != nil {
			return nil, err
		}
	}
	// Get the latest version of the application to return
	return c.GetGroup(group.GroupID)
}

func (c *MarathonClient) ListGroups() (*Groups, error) {
	groups := new(Groups)

	resp := c.http.HttpGet(c.marathonUrl(API_GROUPS), groups)
	if resp.Error != nil {
		return nil, resp.Error
	}
	return groups, nil
}

func (c *MarathonClient) GetGroup(id string) (*Group, error) {
	group := new(Group)
	resp := c.http.HttpGet(c.marathonUrl(API_GROUPS, id), group)
	if resp.Error != nil {
		return nil, resp.Error
	}
	return group, nil
}

func (c *MarathonClient) DestroyGroup(id string) (*DeploymentID, error) {
	deploymentId := new(DeploymentID)
	resp := c.http.HttpDelete(fmt.Sprintf("%s?force=true", c.marathonUrl(API_GROUPS, id)), nil, deploymentId)
	if resp.Error != nil {
		return nil, resp.Error
	}
	return deploymentId, nil
}
