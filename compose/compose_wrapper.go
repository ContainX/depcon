package compose

import (
	"errors"
	"fmt"
	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/project"
	"github.com/gondor/depcon/pkg/envsubst"
	"io/ioutil"
	"os"
	"strings"
)

const (
	DOCKER_TLS_VERIFY string = "DOCKER_TLS_VERIFY"
)

var (
	ErrorParamsMissing = errors.New("One or more ${PARAMS} that were defined in the compose file could not be resolved.")
)

type ComposeWrapper struct {
	context *Context
	project *project.Project
}

func NewCompose(context *Context) Compose {
	c := new(ComposeWrapper)
	c.context = context
	project, err := c.createDockerContext()
	if err != nil {
		log.Fatal(err)
	}
	c.project = project
	for k, v := range c.project.Configs {
		log.Error("%s = \n", k)
		log.Error("%s, %v, %v, %v, %v\n\n", v.Image, v.Links, v.Ports, v.Environment, v)
	}
	return c
}

func (c *ComposeWrapper) Up(services ...string) error {
	return c.project.Up(services...)
}

func (c *ComposeWrapper) Kill(services ...string) error {
	return c.project.Kill(services...)
}

func (c *ComposeWrapper) Build(services ...string) error {
	return c.project.Build(services...)
}

func (c *ComposeWrapper) Restart(services ...string) error {
	return c.project.Restart(services...)
}

func (c *ComposeWrapper) Pull(services ...string) error {
	return c.project.Pull(services...)
}

func (c *ComposeWrapper) Delete(services ...string) error {
	return c.project.Delete(services...)
}

func (c *ComposeWrapper) Logs(services ...string) error {
	return c.project.Log(services...)
}

func (c *ComposeWrapper) Start(services ...string) error {
	return c.execStartStop(true, services...)
}

func (c *ComposeWrapper) Stop(services ...string) error {
	return c.execStartStop(false, services...)
}

func (c *ComposeWrapper) execStartStop(start bool, services ...string) error {
	if start {
		return c.project.Start(services...)
	}
	return c.project.Down(services...)
}

func (c *ComposeWrapper) Port(index int, proto, service, port string) error {

	s, err := c.project.CreateService(service)
	if err != nil {
		return err
	}

	containers, err := s.Containers()
	if err != nil {
		return err
	}

	if index < 1 || index > len(containers) {
		fmt.Errorf("Invalid index %d", index)
	}

	output, err := containers[index-1].Port(fmt.Sprintf("%s/%s", port, proto))
	if err != nil {
		return err
	}
	fmt.Println(output)
	return nil
}

func (c *ComposeWrapper) PS(quiet bool) error {
	allInfo := project.InfoSet{}

	for name := range c.project.Configs {
		service, err := c.project.CreateService(name)
		if err != nil {
			return err
		}

		info, err := service.Info(quiet)
		if err != nil {
			return err
		}

		allInfo = append(allInfo, info...)
	}
	os.Stdout.WriteString(allInfo.String(!quiet))
	return nil
}

func (c *ComposeWrapper) createDockerContext() (*project.Project, error) {

	clientFactory, err := docker.NewDefaultClientFactory(docker.ClientOpts{})
	if err != nil {
		log.Fatal(err)
	}

	tlsVerify := os.Getenv(DOCKER_TLS_VERIFY)

	if tlsVerify == "1" {
		clientFactory, err = docker.NewDefaultClientFactory(docker.ClientOpts{
			TLS:       true,
			TLSVerify: true,
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	if c.context.EnvParams != nil && len(c.context.EnvParams) > 0 {
		file, err := os.Open(c.context.ComposeFile)
		if err != nil {
			return nil, fmt.Errorf("Error opening filename %s, %s", c.context.ComposeFile, err.Error())
		}
		parsed, missing := envsubst.SubstFileTokens(file, c.context.ComposeFile, c.context.EnvParams)
		log.Debug("Map: %v\nParsed: %s\n", c.context.EnvParams, parsed)

		if c.context.ErrorOnMissingParams && missing {
			return nil, ErrorParamsMissing
		}
		file, err = ioutil.TempFile("", "depcon")
		if err != nil {
			return nil, err
		}
		err = ioutil.WriteFile(file.Name(), []byte(parsed), os.ModeTemporary)
		if err != nil {
			return nil, err
		}
		c.context.ComposeFile = file.Name()
	}
	return docker.NewProject(&docker.Context{
		Context: project.Context{
			ComposeFiles: strings.Split(c.context.ComposeFile, ","),
			ProjectName:  c.context.ProjectName,
		},
		ClientFactory: clientFactory,
	})
}
