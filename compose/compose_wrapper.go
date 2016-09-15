package compose

import (
	"errors"
	"fmt"
	"github.com/ContainX/depcon/pkg/envsubst"
	"github.com/docker/distribution/context"
	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/docker/ctx"
	"github.com/docker/libcompose/project"
	"github.com/docker/libcompose/project/options"
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
	project project.APIProject
}

func NewCompose(context *Context) Compose {
	c := new(ComposeWrapper)
	c.context = context
	project, err := c.createDockerContext()
	if err != nil {
		log.Fatal(err)
	}
	c.project = project
	return c
}

func (c *ComposeWrapper) Up(services ...string) error {
	options := options.Up{Create: options.Create{}}
	return c.project.Up(context.Background(), options, services...)
}

func (c *ComposeWrapper) Kill(services ...string) error {
	return c.project.Kill(context.Background(), "SIGKILL", services...)
}

func (c *ComposeWrapper) Build(services ...string) error {
	options := options.Build{}
	return c.project.Build(context.Background(), options, services...)
}

func (c *ComposeWrapper) Restart(services ...string) error {
	timeout := 10
	return c.project.Restart(context.Background(), timeout, services...)
}

func (c *ComposeWrapper) Pull(services ...string) error {
	return c.project.Pull(context.Background(), services...)
}

func (c *ComposeWrapper) Delete(services ...string) error {
	options := options.Delete{}
	return c.project.Delete(context.Background(), options, services...)
}

func (c *ComposeWrapper) Logs(services ...string) error {
	return c.project.Log(context.Background(), true, services...)
}

func (c *ComposeWrapper) Start(services ...string) error {
	return c.execStartStop(true, services...)
}

func (c *ComposeWrapper) Stop(services ...string) error {
	return c.execStartStop(false, services...)
}

func (c *ComposeWrapper) execStartStop(start bool, services ...string) error {
	if start {
		return c.project.Start(context.Background(), services...)
	}
	options := options.Down{}
	return c.project.Down(context.Background(), options, services...)
}

func (c *ComposeWrapper) Port(index int, proto, service, port string) error {

	output, err := c.project.Port(context.Background(), index, proto, service, port)
	if err != nil {
		return err
	}
	fmt.Println(output)
	return nil
}

func (c *ComposeWrapper) PS(quiet bool) error {
	if allInfo, err := c.project.Ps(context.Background()); err == nil {
		os.Stdout.WriteString(allInfo.String([]string{"Name", "Command", "State", "Ports"}, !quiet))
	}
	return nil
}

func (c *ComposeWrapper) createDockerContext() (project.APIProject, error) {

	if c.context.EnvParams != nil && len(c.context.EnvParams) > 0 {
		file, err := os.Open(c.context.ComposeFile)
		if err != nil {
			return nil, fmt.Errorf("Error opening filename %s, %s", c.context.ComposeFile, err.Error())
		}
		parsed, missing := envsubst.SubstFileTokens(file, c.context.EnvParams)
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
	return docker.NewProject(&ctx.Context{
		Context: project.Context{
			ComposeFiles: strings.Split(c.context.ComposeFile, ","),
			ProjectName:  c.context.ProjectName,
		},
	}, nil)
}
