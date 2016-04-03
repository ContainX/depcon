package marathon

import "time"

type AppById struct {
	App Application `json:"app"`
}

type Applications struct {
	Apps []Application `json:"apps"`
}

type Tasks struct {
	Tasks []*Task `json:"tasks"`
}

type Application struct {
	ID                    string              `json:"id,omitempty"`
	Cmd                   string              `json:"cmd,omitempty"`
	Args                  []string            `json:"args,omitempty"`
	AcceptedResourceRoles []string            `json:"acceptedResourceRoles,omitempty"`
	Constraints           [][]string          `json:"constraints,omitempty"`
	Container             *Container          `json:"container,omitempty"`
	CPUs                  float64             `json:"cpus,omitempty"`
	Disk                  float64             `json:"disk,omitempty"`
	Env                   map[string]string   `json:"env,omitempty"`
	Labels                map[string]string   `json:"labels,omitempty"`
	Executor              string              `json:"executor,omitempty"`
	HealthChecks          []*HealthCheck      `json:"healthChecks,omitempty"`
	Instances             int                 `json:"instances,omitempty"`
	Mem                   float64             `json:"mem,omitempty"`
	Tasks                 []*Task             `json:"tasks,omitempty"`
	Ports                 []int               `json:"ports,omitempty"`
	ServicePorts          []int               `json:"servicePorts,omitempty"`
	RequirePorts          bool                `json:"requirePorts,omitempty"`
	BackoffFactor         float64             `json:"backoffFactor,omitempty"`
	DeploymentID          []map[string]string `json:"deployments,omitempty"`
	Dependencies          []string            `json:"dependencies,omitempty"`
	TasksRunning          int                 `json:"tasksRunning,omitempty"`
	TasksStaged           int                 `json:"tasksStaged,omitempty"`
	TasksHealthy          int                 `json:"tasksHealthy,omitempty"`
	TasksUnHealthy        int                 `json:"tasksUnHealthy,omitempty"`
	User                  string              `json:"user,omitempty"`
	UpgradeStrategy       *UpgradeStrategy    `json:"upgradeStrategy,omitempty"`
	Uris                  []string            `json:"uris,omitempty"`
	Version               string              `json:"version,omitempty"`
}

type KillTasksScale struct {
	ids []string `json:"ids"`
}

type AppKillTasksOptions struct {
	Host  string `json:"host"`
	Scale bool   `json:bool`
}

type Versions struct {
	Versions []string
}

type DeploymentID struct {
	DeploymentID string `json:"deploymentId"`
	Version      string `json:"version"`
}

type Container struct {
	Type    string    `json:"type,omitempty"`
	Docker  *Docker   `json:"docker,omitempty"`
	Volumes []*Volume `json:"volumes,omitempty"`
}
type PortMapping struct {
	ContainerPort int    `json:"containerPort,omitempty"`
	HostPort      int    `json:"hostPort"`
	ServicePort   int    `json:"servicePort,omitempty"`
	Protocol      string `json:"protocol"`
}

type Parameters struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type Volume struct {
	ContainerPath string `json:"containerPath,omitempty"`
	HostPath      string `json:"hostPath,omitempty"`
	Mode          string `json:"mode,omitempty"`
}

type Docker struct {
	ForcePullImage bool           `json:"forcePullImage,omitempty"`
	Image          string         `json:"image,omitempty"`
	Network        string         `json:"network,omitempty"`
	Parameters     []*Parameters  `json:"parameters,omitempty"`
	PortMappings   []*PortMapping `json:"portMappings,omitempty"`
	Privileged     bool           `json:"privileged,omitempty"`
}

type UpgradeStrategy struct {
	MinimumHealthCapacity float64 `json:"minimumHealthCapacity,omitempty"`
	MaximumOverCapacity   float64 `json:"maximumOverCapacity,omitempty"`
}

type HealthCheck struct {
	Protocol               string `json:"protocol,omitempty"`
	Path                   string `json:"path,omitempty"`
	GracePeriodSeconds     int    `json:"gracePeriodSeconds,omitempty"`
	IntervalSeconds        int    `json:"intervalSeconds,omitempty"`
	PortIndex              int    `json:"portIndex,omitempty"`
	MaxConsecutiveFailures int    `json:"maxConsecutiveFailures,omitempty"`
	TimeoutSeconds         int    `json:"timeoutSeconds,omitempty"`
}

type HealthCheckResult struct {
	Alive               bool   `json:"alive"`
	ConsecutiveFailures int    `json:"consecutiveFailures"`
	FirstSuccess        string `json:"firstSuccess"`
	LastFailure         string `json:"lastFailure"`
	LastSuccess         string `json:"lastSuccess"`
	TaskID              string `json:"taskId"`
}

type Task struct {
	AppID             string               `json:"appId"`
	Host              string               `json:"host"`
	ID                string               `json:"id"`
	HealthCheckResult []*HealthCheckResult `json:"healthCheckResults"`
	Ports             []int                `json:"ports"`
	ServicePorts      []int                `json:"servicePorts"`
	StagedAt          string               `json:"stagedAt"`
	StartedAt         string               `json:"startedAt"`
	Version           string               `json:"version"`
}

type QueuedTask struct {
	App   *Application    `json:"app"`
	Delay map[string]bool `json:"delay"`
}

type Queue struct {
	Queue []QueuedTask `json:"queue"`
}

type Which struct {
	Leader string `json:"leader"`
}

type Message struct {
	Message string `json:"message"`
}

type Deploys []Deploy

type Deploy struct {
	AffectedApps   []string `json:"affectedApps"`
	DeployID       string   `json:"id"`
	Steps          [][]Step `json:"steps"`
	CurrentActions []Step   `json:"currentActions"`
	Version        string   `json:"version"`
	CurrentStep    int      `json:"currentStep"`
	TotalSteps     int      `json:"totalSteps"`
}

type Step struct {
	Action string `json:"action"`
	App    string `json:"app"`
}

type Group struct {
	GroupID      string         `json:"id"`
	Version      string         `json:"version,omitempty"`
	Apps         []*Application `json:"apps,omitempty"`
	Dependencies []string       `json:"dependencies,omitempty"`
	Groups       []*Group       `json:"groups,omitempty"`
}

type Groups struct {
	GroupID      string         `json:"id"`
	Version      string         `json:"version"`
	Apps         []*Application `json:"apps"`
	Dependencies []string       `json:"dependencies"`
	Groups       []*Group       `json:"groups"`
}

type AppRestartOptions struct {
	Force bool `json:"force"`
}

type MarathonPing struct {
	Host    string
	Elapsed time.Duration
}

type MarathonInfo struct {
	EventSubscriber struct {
		HttpEndpoints []string `json:"http_endpoints"`
		Type          string   `json:"type"`
	} `json:"event_subscriber"`
	FrameworkId string `json:"frameworkId"`
	HttpConfig  struct {
		AssetsPath interface{} `json:"assets_path"`
		HttpPort   float64     `json:"http_port"`
		HttpsPort  float64     `json:"https_port"`
	} `json:"http_config"`
	Leader         string `json:"leader"`
	MarathonConfig struct {
		Checkpoint                 bool    `json:"checkpoint"`
		Executor                   string  `json:"executor"`
		FailoverTimeout            float64 `json:"failover_timeout"`
		Ha                         bool    `json:"ha"`
		Hostname                   string  `json:"hostname"`
		LocalPortMax               float64 `json:"local_port_max"`
		LocalPortMin               float64 `json:"local_port_min"`
		Master                     string  `json:"master"`
		MesosRole                  string  `json:"mesos_role"`
		MesosUser                  string  `json:"mesos_user"`
		ReconciliationInitialDelay float64 `json:"reconciliation_initial_delay"`
		ReconciliationInterval     float64 `json:"reconciliation_interval"`
		TaskLaunchTimeout          float64 `json:"task_launch_timeout"`
	} `json:"marathon_config"`
	Name            string `json:"name"`
	Version         string `json:"version"`
	ZookeeperConfig struct {
		Zk              string `json:"zk"`
		ZkFutureTimeout struct {
			Duration float64 `json:"duration"`
		} `json:"zk_future_timeout"`
		ZkHosts   string  `json:"zk_hosts"`
		ZkPath    string  `json:"zk_path"`
		ZkState   string  `json:"zk_state"`
		ZkTimeout float64 `json:"zk_timeout"`
	} `json:"zookeeper_config"`
}

type LeaderInfo struct {
	Leader string `json:"leader"`
}
