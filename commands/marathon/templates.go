package marathon

import (
	"fmt"
	"github.com/ContainX/depcon/marathon"
	"github.com/ContainX/depcon/utils"
	"strconv"
	"text/template"
	"github.com/ContainX/depcon/pkg/cli"
	"time"
	"io"
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
{{end}}`

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
	Data     interface{}
	Template string
}

func (d Templated) ToColumns(output io.Writer) error {
	w := cli.NewTabWriter(output)
	t := template.New("output").Funcs(buildFuncMap())
	t, _ = t.Parse(d.Template)
	if err := t.Execute(w, d.Data); err != nil {
		return err
	}
	cli.FlushWriter(w)
	return nil
}

func buildFuncMap() template.FuncMap {
	funcMap := template.FuncMap{
		"intConcat":     utils.ConcatInts,
		"idConcat":	 utils.ConcatIdentifiers,
		"dockerImage":   dockerImageOrEmpty,
		"floatToString": floatToString,
		"intToString":   strconv.Itoa,
		"valString":	 valueToString,
		"pad":           padString,
		"hasDocker":     hasDocker,
		"fdate":	 cli.FormatDate,
		"msDur":	 durationToMilliseconds,
	}
	return funcMap
}

func durationToMilliseconds(t time.Duration) string {
	return fmt.Sprintf("%d ms", (t.Nanoseconds() / int64(time.Millisecond)))
}

func padString(s string) string {
	return fmt.Sprintf("%-25s:", s)
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

func floatToString(f float64) string {
	return fmt.Sprintf("%.2f", f)
}

func valueToString(v interface{}) string {
	return fmt.Sprintf("%v", v)
}
