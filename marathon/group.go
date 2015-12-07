package marathon

import (
	"fmt"
	"github.com/gondor/depcon/pkg/encoding"
	"github.com/gondor/depcon/pkg/envsubst"
	"os"
	"time"
)

func (c *MarathonClient) CreateGroupFromFile(filename string, opts *CreateOptions) (*Group, error) {
	log.Info("Creating Group from file: %s", filename)

	options := initCreateOptions(opts)

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Error opening filename %s, %s", filename, err.Error())
	}

	var encoder encoding.Encoder
	encoder, err = encoding.NewEncoderFromFileExt(filename)
	if err != nil {
		return nil, err
	}

	parsed, missing := envsubst.SubstFileTokens(file, filename, options.EnvParams)

	if options.ErrorOnMissingParams && missing {
		return nil, ErrorAppParamsMissing
	}

	group := new(Group)
	err = encoder.UnMarshalStr(parsed, &group)
	if err != nil {
		return nil, err
	}
	return c.CreateGroup(group, options.Wait, options.Force)
}

func (c *MarathonClient) CreateGroup(group *Group, wait, force bool) (*Group, error) {
	log.Info("Creating Group '%s', wait: %v, force: %v", group.GroupID, wait, force)
	result := new(DeploymentID)
	resp := c.http.HttpPost(c.marathonUrl(API_GROUPS), group, result)
	if resp.Error != nil {
		return nil, resp.Error
	}
	if wait {
		if err := c.WaitForDeployment(result.DeploymentID, time.Duration(500)*time.Second); err != nil {
			return nil, err
		}
	}
	return group, nil
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
