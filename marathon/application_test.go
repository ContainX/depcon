package marathon

import (
	"github.com/ContainX/depcon/pkg/mockrest"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	AppsFolder   = "testdata/apps/"
	CommonFolder = "testdata/common/"
)

func TestCreateApplicationFromFile(t *testing.T) {

	envParams := make(map[string]string, 1)
	envParams["NODE_EXPORTER_VERSION"] = "1"

	opts := &CreateOptions{Wait: false, Force: true, ErrorOnMissingParams: true, EnvParams: envParams}

	c := MarathonClient{}
	app, err := c.ParseApplicationFromFile(AppsFolder+"app_params.json", opts)
	if err != nil {
		log.Panicf("Expected success %v", err)
	}
	log.Debug("%v", app)
	if app.Labels["tags"] != "prom-metrics" {
		log.Panic("Expected Labels parsed correctly")
	}
	if app.Instances != int(2) {
		log.Panic("Expected Instances parsed correctly")
	}
}

func TestListApplications(t *testing.T) {
	s := mockrest.StartNewWithFile(AppsFolder + "list_apps_response.json")
	defer s.Stop()

	c := NewMarathonClient(s.URL, "", "", "")
	apps, err := c.ListApplications()

	assert.Nil(t, err, "Error response was not expected")
	assert.Equal(t, "/myapp", apps.Apps[0].ID)
}

func TestGetApplication(t *testing.T) {
	s := mockrest.StartNewWithFile(AppsFolder + "get_app_response.json")
	defer s.Stop()

	c := NewMarathonClient(s.URL, "", "", "")
	app, err := c.GetApplication("storage/redis-x")

	assert.Nil(t, err, "Error response was not expected")
	assert.Equal(t, "/storage/redis-x", app.ID)
	assert.Equal(t, 1, len(app.Ports))
	assert.Equal(t, "redis", app.Container.Docker.Image)
	assert.Equal(t, "cache", app.Labels["role"])
}

func TestHasApplication(t *testing.T) {
	s := mockrest.StartNewWithFile(AppsFolder + "get_app_response.json")
	defer s.Stop()

	c := NewMarathonClient(s.URL, "", "", "")
	ok, err := c.HasApplication("/storage/redis-x")

	assert.Nil(t, err, "Error response was not expected")
	assert.Equal(t, true, ok)
}

func TestHasApplicationInvalid(t *testing.T) {
	s := mockrest.StartNewWithStatusCode(404)
	defer s.Stop()

	c := NewMarathonClient(s.URL, "", "", "")
	ok, _ := c.HasApplication("/storage/redis-invalid")

	assert.Equal(t, false, ok)
}

func TestDestroyApplication(t *testing.T) {
	s := mockrest.StartNewWithFile(CommonFolder + "deployid_response.json")
	defer s.Stop()

	c := NewMarathonClient(s.URL, "", "", "")
	depId, err := c.DestroyApplication("/someapp")
	assert.Nil(t, err, "Error response was not expected")
	assert.Equal(t, "5ed4c0c5-9ff8-4a6f-a0cd-f57f59a34b43", depId.DeploymentID)
}

func TestRestartApplication(t *testing.T) {
	s := mockrest.StartNewWithFile(CommonFolder + "deployid_response.json")
	defer s.Stop()

	c := NewMarathonClient(s.URL, "", "", "")
	depId, err := c.DestroyApplication("/someapp")
	assert.Nil(t, err, "Error response was not expected")
	assert.Equal(t, "5ed4c0c5-9ff8-4a6f-a0cd-f57f59a34b43", depId.DeploymentID)
}

func TestScaleApplication(t *testing.T) {
	s := mockrest.StartNewWithFile(CommonFolder + "deployid_response.json")
	defer s.Stop()

	c := NewMarathonClient(s.URL, "", "", "")
	depId, err := c.ScaleApplication("/someapp", 5)
	assert.Nil(t, err, "Error response was not expected")
	assert.Equal(t, "5ed4c0c5-9ff8-4a6f-a0cd-f57f59a34b43", depId.DeploymentID)
}

func TestCreateApplicationInvalidAppId(t *testing.T) {
	s := mockrest.StartNewWithStatusCode(422)
	defer s.Stop()

	c := NewMarathonClient(s.URL, "", "", "")
	_, err := c.CreateApplication(NewApplication("/someapp"), false, false)

	assert.NotNil(t, err, "Expecting Error")
}

func TestNewApplication(t *testing.T) {
	app := NewApplication("/some/application")
	assert.Equal(t, "/some/application", app.ID)
}
