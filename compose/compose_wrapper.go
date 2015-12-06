package compose

import (
	"fmt"
	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/project"
	"os"
)

type ComposeWrapper struct {
	context *Context
}

func NewCompose(context *Context) Compose {
	c := new(ComposeWrapper)
	c.context = context
	return c
}

func (c *ComposeWrapper) Up(envParams map[string]string, ErrorOnMissingParams bool, services ...string) error {

	project, err := c.createDockerContext()

	if err != nil {
		return err
	}
	return project.Up(services...)
}

func (c *ComposeWrapper) Kill(services ...string) error {
	project, err := c.createDockerContext()

	if err != nil {
		return err
	}
	return project.Kill(services...)
}

func (c *ComposeWrapper) Build(services ...string) error {
	project, err := c.createDockerContext()

	if err != nil {
		return err
	}
	return project.Build(services...)
}

func (c *ComposeWrapper) Restart(services ...string) error {
	project, err := c.createDockerContext()

	if err != nil {
		return err
	}
	return project.Restart(services...)
}

func (c *ComposeWrapper) Pull(services ...string) error {
	project, err := c.createDockerContext()

	if err != nil {
		return err
	}
	return project.Pull(services...)
}

func (c *ComposeWrapper) Delete(services ...string) error {
	project, err := c.createDockerContext()

	if err != nil {
		return err
	}
	return project.Delete(services...)
}

func (c *ComposeWrapper) Logs(services ...string) error {
	project, err := c.createDockerContext()

	if err != nil {
		return err
	}
	return project.Log(services...)
}

func (c *ComposeWrapper) Start(services ...string) error {
	return c.execStartStop(true, services...)
}

func (c *ComposeWrapper) Stop(services ...string) error {
	return c.execStartStop(false, services...)
}

func (c *ComposeWrapper) execStartStop(start bool, services ...string) error {
	project, err := c.createDockerContext()

	if err != nil {
		return err
	}
	if start {
		return project.Start(services...)
	}
	return project.Down(services...)
}

func (c *ComposeWrapper) Port(index int, proto, service, port string) error {
	project, err := c.createDockerContext()
	if err != nil {
		return err
	}

	s, err := project.CreateService(service)
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
	p, err := c.createDockerContext()
	allInfo := project.InfoSet{}

	if err != nil {
		return err
	}

	for name := range p.Configs {
		service, err := p.CreateService(name)
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

	tlsVerify := os.Getenv("DOCKER_TLS_VERIFY")

	if tlsVerify == "1" {
		clientFactory, err = docker.NewDefaultClientFactory(docker.ClientOpts{
			TLS:       true,
			TLSVerify: true,
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	return docker.NewProject(&docker.Context{
		Context: project.Context{
			ComposeFile: c.context.ComposeFile,
			ProjectName: c.context.ProjectName,
		},
		ClientFactory: clientFactory,
	})
}
