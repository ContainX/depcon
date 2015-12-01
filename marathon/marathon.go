// Marathon API
package marathon

import (
	"github.com/shirkevich/depcon/pkg/httpclient"
	"github.com/gondor/depcon/pkg/logger"
	"github.com/gondor/depcon/utils"
	"time"
)

const (
	DEFAULT_EVENTS_URL = "event"

	/* --- api related constants --- */
	API_VERSION      = "v2"
	API_SUBSCRIPTION = API_VERSION + "/eventSubscriptions"
	API_APPS         = API_VERSION + "/apps"
	API_TASKS        = API_VERSION + "/tasks"
	API_DEPLOYMENTS  = API_VERSION + "/deployments"
	API_GROUPS       = API_VERSION + "/groups"
	API_QUEUE        = API_VERSION + "/queue"
	API_INFO         = API_VERSION + "/info"
	API_LEADER       = API_VERSION + "/leader"
	API_PING         = "ping"
	API_LOGGING      = "logging"
	API_HELP         = "help"
	API_METRICS      = "metrics"

	DefaultTimeout = time.Duration(90) * time.Second
)

// Common package logger
var log = logger.GetLogger("depcon.marathon")

type CreateOptions struct {
	// if true will attempt to wait until the new application or group is running
	Wait bool
	// if true and an application/group already exists an update will be performed.
	// if false and an application/group exists an error will be returned
	Force bool
	// If true an error will be returned on params defined in the configuration file that
	// could not resolve to user input and environment variables
	ErrorOnMissingParams bool
	// Additional environment params - looks at this map for token substitution which takes
	// priority over matching environment variables
	EnvParams map[string]string
}

type Marathon interface {

	/** Application API */

	// Creates a new Application from a file and replaces tokenized variables
	// with resolved environment values
	//
	// {filename} - the application file of type [ json | yaml ]
	// {opts}     - create application options
	CreateApplicationFromFile(filename string, opts *CreateOptions) (*Application, error)

	// Creates a new Application
	// {app}   - the application structure containing configuration
	// {wait}  - if true will attempt to wait until the new application is running
	// {force} - if true and a application already exists an update will be performed.
	//         - if false and a application exists an error will be returned
	CreateApplication(app *Application, wait, force bool) (*Application, error)

	// Updates an Application
	// {app} - the application structure containing configuration
	// {wait} - if true will attempt to wait until the application updated is running
	UpdateApplication(app *Application, wait bool) (*Application, error)

	// List all applications on a Marathon cluster
	ListApplications() (*Applications, error)

	// Get an Application by Id
	// {id} - application identifier
	GetApplication(id string) (*Application, error)

	// Determines if the application exists
	// {id} - the application identifier
	HasApplication(id string) (bool, error)

	// Removes an Application by Id and all of it's running instances
	// {id} - application identifier
	DestroyApplication(id string) (*DeploymentID, error)

	// Restarts an Application
	// {id} - application identifier
	// {force} - forces a restart if true
	RestartApplication(id string, force bool) (*DeploymentID, error)

	// Scale an Application by Id and Instances
	// {id} - application identifier
	// {instances} - instances to scale to
	ScaleApplication(id string, instances int) (*DeploymentID, error)

	// List application versions that have been deployed to Marathon
	// {id} - the application identifier
	ListVersions(id string) (*Versions, error)

	// Attempts to wait for an application to be running
	// {id} - the application id
	// {timeout} - the max time to wait
	WaitForApplication(id string, timeout time.Duration) error

	// Attempts to wait for an application to be running and healthy (all health checks for all tasks passing)
	// {id} - the application id
	// {timeout} - the max time to wait
	WaitForApplicationHealthy(id string, timeout time.Duration) error

	/** Deployment API */

	// Determines whether a deployment for the specified Id exists
	// {id} - deployment identifier
	HasDeployment(id string) (bool, error)

	// List the current deployments
	ListDeployments() ([]*Deploy, error)

	// Deletes a deployment
	// {id} - deployment identifier
	DeleteDeployment(id string) (*DeploymentID, error)

	// Waits for a deployment to finish for max timeout duration
	WaitForDeployment(id string, timeout time.Duration) error

	/** Group API */

	// Creates a new group from a file and replaces tokenized variables
	// with resolved environment values
	//
	// {filename} - the group file of type [ json | yaml ]
	// {opts}     - create application options
	CreateGroupFromFile(filename string, opts *CreateOptions) (*Group, error)

	// Creates a new Group
	// {group} - the group structure containing configuration
	// {wait}  - if true will attempt to wait until the new group is running
	// {force} - if true and a group already exists an update will be performed.
	//         - if false and a group exists an error will be returned
	CreateGroup(group *Group, wait, force bool) (*Group, error)

	// List all groups
	ListGroups() (*Groups, error)

	// Get a Group by Id
	// {id} - group identifier
	GetGroup(id string) (*Group, error)

	// Removes a Group by Id and all of it's related resources (application instances)
	// {id} - group identifier
	DestroyGroup(id string) (*DeploymentID, error)

	/** Task API */

	// List all running tasks
	ListTasks() ([]*Task, error)

	// List all tasks for an application
	// {id} - the application identifier
	GetTasks(id string) ([]*Task, error)

	// Kills application tasks for the app identifier
	// {id} - the application identifier
	// {host} - host to kill tasks on or empty (default)
	// {scale} - Scale the app down (i.e. decrement its instances setting by the number of tasks killed), false is default
	KillAppTasks(id string, host string, scale bool) ([]*Task, error)

	// Kill the task with ID taskId that belongs to an application
	// {taskId} - the task id
	// {scale}  - Scale the app down (ie. decrement it's instances setting by the number of tasks killed). Default: false
	KillAppTask(taskId string, scale bool) (*Task, error)

	// List Queue - tasks currently pending
	ListQueue() (*Queue, error)

	/** Marathon Server Info API */

	// Pings the Marathon host via the /ping endpoint
	Ping() (*MarathonPing, error)

	// Get info about the Marathon Instance
	GetMarathonInfo() (*MarathonInfo, error)

	// Get the current Marathon leader
	GetCurrentLeader() (*LeaderInfo, error)

	// Abdicates the current leader
	AbdicateLeader() (*Message, error)
}

type MarathonClient struct {
	http httpclient.HttpClient
	host string
}

func NewMarathonClient(host, username, password string) Marathon {
	httpConfig := httpclient.NewDefaultConfig()
	httpConfig.HttpUser = username
	httpConfig.HttpPass = password

	httpClient := httpclient.NewHttpClient(*httpConfig)

	c := new(MarathonClient)
	c.http = *httpClient
	c.host = host
	return c
}

func (c *MarathonClient) marathonUrl(elements ...string) string {
	return utils.BuildPath(c.host, elements)
}

func initCreateOptions(opts *CreateOptions) *CreateOptions {
	if opts == nil {
		return &CreateOptions{}
	}
	return opts
}
