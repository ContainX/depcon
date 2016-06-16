package marathon

import (
	"github.com/ContainX/depcon/marathon"
	"github.com/ContainX/depcon/pkg/cli"
	"github.com/ContainX/depcon/utils"
	"io"
	"text/template"
)

const (
	T_APPLICATIONS = `
{{ "ID" }}	{{ "INSTANCES" }}	{{ "CPU" }}	{{ "MEM" }}	{{ "PORTS" }}	{{ "CONTAINER" }}	{{ "VERSION" }}
{{range .Apps}}{{ .ID }}	{{ .Instances }}	{{ .CPUs | floatToString }}	{{ .Mem | floatToString }}	{{ .Ports | intConcat }}	{{ .Container | dockerImage }}	{{ .Version }}
{{end}}`

	T_APPLICATION = `
{{ "ID" }}	{{ .ID }}
{{ "CPUs:" }}	{{ .CPUs | floatToString }}
{{ "Memory:" }}	{{ .Mem | floatToString }}
{{ "Ports:" }}	{{ .Ports | intConcat }}
{{ "Instances:" }}	{{ .Instances | intToString }}
{{ "Version:" }}	{{ .Version }}
{{ "Tasks:" }}	{{ "Staged" | pad }} {{ .TasksStaged | intToString }}
	{{ "Running" | pad }} {{ .TasksRunning | intToString }}
	{{ "Healthy" | pad }} {{ .TasksHealthy | intToString }}
	{{ "UnHealthy" | pad }} {{ .TasksUnHealthy | intToString }}

{{ if hasDocker .Container }}
{{ "Container: " }}	{{ "Type" | pad }} {{ "Docker" }}
	{{ "Image" | pad }} {{ .Container.Docker.Image }}
	{{ "Network" | pad }} {{ .Container.Docker.Network }}
{{- end}}
{{ "Environment:" }}
{{ range $key, $value := .Env }}		{{ $key | pad }} {{ $value }}
{{end}}
{{ "Labels:" }}
{{ range $key, $value := .Labels }}		{{ $key | pad }} {{ $value }}
{{end}}
`

	T_VERSIONS = `
{{ "VERSIONS" }}
{{ range .Versions }}{{ . }}
{{end}}`

	T_DEPLOYMENT_ID = `
{{ "DEPLOYMENT_ID" }}	{{ "VERSION" }}
{{ .DeploymentID }}	{{ .Version }}`

	T_TASKS = `
{{ "APP_ID" }}	{{ "HOST" }}	{{ "VERSION" }}	{{ "STARTED" }}	{{ "TASK_ID" }}
{{ range . }}{{ .AppID }}	{{ .Host }}	{{ .Version }}	{{ .StartedAt | fdate }}	{{ .ID }}
{{end}}`

	T_TASK = `
{{ "ID:" }}	{{ .ID }}
{{ "AppID:" }}	{{ .AppID }}
{{ "Staged:" }}	{{ .StagedAt | fdate }}
{{ "Started:" }}	{{ .StartedAt | fdate }}
{{ "Host:"	}}	{{ .Host }}
{{ "Ports:" }}	{{ .Ports | intConcat }}
`
	T_DEPLOYMENTS = `
{{ "DEPLOYMENT_ID" }}	{{ "VERSION" }} 	{{ "PROGRESS" }}	{{ "APPS" }}
{{ range . }}{{ .DeployID }}	{{ .Version }}	{{ .CurrentStep | intToString }}/{{ .TotalSteps | intToString }}	{{ .AffectedApps | idConcat }}
{{end}}`
	T_LEADER_INFO = `
{{ "Leader:" }}	{{ .Leader }}
`

	T_PING = `
{{ "HOST" }}	{{ "DURATION" }}
{{ .Host }}	{{ .Elapsed | msDur }}
`

	T_MARATHON_INFO = `
{{ "INFO" }}
{{ "Name:" }}	{{ .Name }}
{{ "Version:" }}	{{ .Version }}
{{ "FrameworkId:" }}	{{ .FrameworkId }}
{{ "Leader:" }}	{{ .Leader }}

{{ "HTTP CONFIG" }}
{{ "HTTP Port:" }}	{{ .HttpConfig.HttpPort | valString }}
{{ "HTTPS Port:" }}	{{ .HttpConfig.HttpsPort | valString }}

{{ "MARATHON CONFIG" }}
{{ "Checkpoint:" }}	{{ .MarathonConfig.Checkpoint | valString }}
{{ "Executor:" }}	{{ .MarathonConfig.Executor }}
{{ "HA:" }}	{{ .MarathonConfig.Ha | valString }}
{{ "Master:" }}	{{ .MarathonConfig.Master }}
{{ "Failover Timeout:" }}	{{ .MarathonConfig.FailoverTimeout | valString }}
{{ "Local Port (Min):" }}	{{ .MarathonConfig.LocalPortMin | valString }}
{{ "Local Port (Max):" }}	{{ .MarathonConfig.LocalPortMax | valString }}

{{ "ZOOKEEPER CONFIG" }}
{{ "ZK:" }}	{{ .ZookeeperConfig.Zk }}
{{ "Timeout:" }}	{{ .ZookeeperConfig.ZkTimeout | valString }}
`
	T_QUEUED_TASKS = `
{{ "APP_ID" }}	{{ "VERSION" }}	{{ "OVERDUE" }}
{{ range .Queue }}{{ .App.ID }}	{{ .App.Version }}	{{ .Delay.overdue | valString }}
{{end}}`

	T_MESSAGE = `
{{ "Message:" }}	{{ .Message }}
`
	T_GROUPS = `
{{ "ID" }}	{{ "VERSION" }}	{{ "GROUPS" }}	{{ "APPS" }}
{{ range . }}{{ .GroupID }}	{{ .Version }}	{{ .Groups | len | valString }}	{{ .Apps | len | valString }}
{{end}}`
)

type Templated struct {
	cli.FormatData
}

func templateFor(template string, data interface{}) Templated {
	return Templated{cli.FormatData{Template: template, Data: data, Funcs: buildFuncMap()}}
}

func (d Templated) ToColumns(output io.Writer) error {
	return d.FormatData.ToColumns(output)
}

func (d Templated) Data() cli.FormatData {
	return d.FormatData
}

func buildFuncMap() template.FuncMap {
	funcMap := template.FuncMap{
		"intConcat":   utils.ConcatInts,
		"idConcat":    utils.ConcatIdentifiers,
		"dockerImage": dockerImageOrEmpty,
		"hasDocker":   hasDocker,
	}
	return funcMap
}

func hasDocker(c *marathon.Container) bool {
	return c != nil && c.Docker != nil
}

func dockerImageOrEmpty(c *marathon.Container) string {
	if c != nil && c.Docker != nil {
		return c.Docker.Image
	}
	return ""
}
