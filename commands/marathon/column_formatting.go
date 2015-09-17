package marathon

import (
	"fmt"
	"github.com/gondor/depcon/marathon"
	"github.com/gondor/depcon/pkg/cli"
	"github.com/gondor/depcon/utils"
	"io"
	"strconv"
	"text/tabwriter"
	"time"
)

const (
	OBJ_FMT       string = "%s:\t%s\n"
	OBJ_FLFMT     string = "%s:\t%.2f\n"
	OBJ_VFMT      string = "%s:\t%v\n"
	OBJ_VAL_SFMT  string = "\t%s\n"
	OBJ_KVFMT     string = "%s:\t%-25s: %s\n"
	OBJ_VAL_KVFMT string = "\t%-25s: %s\n"
)

type Applications struct {
	*marathon.Applications
}

type Versions struct {
	*marathon.Versions
}

type Application struct {
	*marathon.Application
}

type DeploymentId struct {
	*marathon.DeploymentID
}

type Deployments struct {
	Deploys []*marathon.Deploy
}

type AppGroups struct {
	*marathon.Groups
}

type AppGroup struct {
	*marathon.Group
}

type Tasks struct {
	Tasks []*marathon.Task
}

type Task struct {
	*marathon.Task
}

type MarInfo struct {
	*marathon.MarathonInfo
}

type MarLeaderInfo struct {
	*marathon.LeaderInfo
}

type MarPing struct {
	*marathon.MarathonPing
}

type Message struct {
	*marathon.Message
}

type QueuedTasks struct {
	*marathon.Queue
}

func (a Applications) ToColumns(output io.Writer) error {
	w := cli.NewTabWriter(output)
	fmt.Fprintln(w, "\nID\tINSTANCES\tCPU\tMEM\tPORTS\tCONTAINER")
	for _, e := range a.Apps {
		var container string
		if e.Container != nil && e.Container.Docker != nil {
			container = e.Container.Docker.Image
		}
		fmt.Fprintf(w, "%s\t%d\t%.2f\t%.2f\t%s\t%s\n", e.ID, e.Instances, e.CPUs, e.Mem, utils.ConcatInts(e.Ports), container)
	}
	cli.FlushWriter(w)
	return nil
}

func (a Application) ToColumns(output io.Writer) error {
	w := cli.NewTabWriter(output)
	fmt.Fprintln(w, "")
	fmt.Fprintf(w, OBJ_FMT, "ID", a.ID)
	fmt.Fprintf(w, OBJ_FLFMT, "CPUs", a.CPUs)
	fmt.Fprintf(w, OBJ_FLFMT, "Memory", a.Mem)
	fmt.Fprintf(w, OBJ_FMT, "Ports", utils.ConcatInts(a.Ports))
	fmt.Fprintf(w, OBJ_FMT, "Instances", strconv.Itoa(a.Instances))
	fmt.Fprintf(w, OBJ_KVFMT, "Tasks", "Staged", strconv.Itoa(a.TasksStaged))
	fmt.Fprintf(w, OBJ_VAL_KVFMT, "Running", strconv.Itoa(a.TasksRunning))
	fmt.Fprintf(w, OBJ_VAL_KVFMT, "Healthy", strconv.Itoa(a.TasksHealthy))
	fmt.Fprintf(w, OBJ_VAL_KVFMT, "UnHealthy", strconv.Itoa(a.TasksUnHealthy))
	fmt.Fprintln(w, "")
	if a.Container != nil && a.Container.Docker != nil {
		fmt.Fprintf(w, OBJ_KVFMT, "Container", "Type", "Docker")
		fmt.Fprintf(w, OBJ_VAL_KVFMT, "Image", a.Container.Docker.Image)
		fmt.Fprintf(w, OBJ_VAL_KVFMT, "Network", a.Container.Docker.Network)
		fmt.Fprintln(w, "")
	}

	writeMapColumn(w, "Environment", a.Env)
	cli.FlushWriter(w)
	return nil
}

func (v Versions) ToColumns(output io.Writer) error {
	w := cli.NewTabWriter(output)
	fmt.Fprintln(w, "\nVERSIONS")
	for _, e := range v.Versions.Versions {
		fmt.Fprintf(w, "%s\n", e)
	}
	cli.FlushWriter(w)
	return nil
}

func (d DeploymentId) ToColumns(output io.Writer) error {
	w := cli.NewTabWriter(output)
	fmt.Fprintln(w, "\nDEPLOYMENT_ID\tVERSION")
	fmt.Fprintf(w, "%s\t%s\n", d.DeploymentID.DeploymentID, d.Version)
	cli.FlushWriter(w)
	return nil
}

func (d Deployments) ToColumns(output io.Writer) error {
	w := cli.NewTabWriter(output)
	fmt.Fprintln(w, "\nDEPLOYMENT_ID\tVERSION\tPROGRESS\tAPPS")
	for _, deploy := range d.Deploys {
		apps := utils.ConcatIdentifiers(deploy.AffectedApps)
		fmt.Fprintf(w, "%s\t%s\t%d/%d\t%v\n", deploy.DeployID, deploy.Version, deploy.CurrentStep, deploy.TotalSteps, apps)
	}
	cli.FlushWriter(w)
	return nil
}

func (t Tasks) ToColumns(output io.Writer) error {
	w := cli.NewTabWriter(output)
	fmt.Fprintln(w, "\nAPP_ID\tHOST\tVERSION\tSTARTED_AT\tTASK_ID")
	for _, e := range t.Tasks {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", e.AppID, e.Host, e.Version, cli.FormatDate(e.StartedAt), e.ID)
	}
	cli.FlushWriter(w)
	return nil
}

func (t Task) ToColumns(output io.Writer) error {
	w := cli.NewTabWriter(output)

	fmt.Fprintf(w, OBJ_FMT, "ID", t.ID)
	fmt.Fprintf(w, OBJ_FMT, "AppID", t.AppID)
	fmt.Fprintf(w, OBJ_FMT, "Version", t.Version)
	fmt.Fprintf(w, OBJ_FMT, "Staged", cli.FormatDate(t.StagedAt))
	fmt.Fprintf(w, OBJ_FMT, "Started", cli.FormatDate(t.StartedAt))
	fmt.Fprintf(w, OBJ_FMT, "Host", t.Host)
	fmt.Fprintf(w, OBJ_FMT, "Ports", utils.ConcatInts(t.Ports))

	cli.FlushWriter(w)
	return nil
}

func (q QueuedTasks) ToColumns(output io.Writer) error {
	w := cli.NewTabWriter(output)
	fmt.Fprintln(w, "\nAPP_ID\tVERSION\tOVERDUE")
	for _, e := range q.Queue.Queue {
		fmt.Fprintf(w, "%s\t%s\t%v\n", e.App.ID, e.App.Version, e.Delay["overdue"])
	}
	cli.FlushWriter(w)
	return nil
}

func (g AppGroup) ToColumns(output io.Writer) error {
	w := cli.NewTabWriter(output)
	fmt.Fprintln(w, "\nID\tVERSION\tGROUPS\tAPPS")
	populateGroupAndChildren(w, g.Group)
	cli.FlushWriter(w)
	return nil
}

func (g AppGroups) ToColumns(output io.Writer) error {
	w := cli.NewTabWriter(output)
	fmt.Fprintln(w, "\nID\tVERSION\tGROUPS\tAPPS")
	for _, group := range g.Groups.Groups {
		populateGroupAndChildren(w, group)
	}
	cli.FlushWriter(w)
	return nil
}

func populateGroupAndChildren(w *tabwriter.Writer, g *marathon.Group) {
	fmt.Fprintf(w, "%s\t%s\t%v\t%v\n", g.GroupID, g.Version, len(g.Groups), len(g.Apps))
	// iterate child groups
	for _, cg := range g.Groups {
		populateGroupAndChildren(w, cg)
	}
}

func (m MarInfo) ToColumns(output io.Writer) error {
	w := cli.NewTabWriter(output)
	fmt.Fprintln(w, "\nINFO\t ")
	fmt.Fprintf(w, OBJ_FMT, "Name", m.Name)
	fmt.Fprintf(w, OBJ_FMT, "Version", m.Version)
	fmt.Fprintf(w, OBJ_FMT, "FrameworkId", m.FrameworkId)
	fmt.Fprintf(w, OBJ_FMT, "Leader", m.Leader)
	fmt.Fprintln(w, " \t \nHTTP CONFIG\t ")
	fmt.Fprintf(w, OBJ_VFMT, "HTTP_Port", m.HttpConfig.HttpPort)
	fmt.Fprintf(w, OBJ_VFMT, "HTTPS_Port", m.HttpConfig.HttpsPort)
	fmt.Fprintln(w, " \t \nMARATHON CONFIG\t ")
	fmt.Fprintf(w, OBJ_VFMT, "Checkpoint", m.MarathonConfig.Checkpoint)
	fmt.Fprintf(w, OBJ_FMT, "Executor", m.MarathonConfig.Executor)
	fmt.Fprintf(w, OBJ_VFMT, "HA", m.MarathonConfig.Ha)
	fmt.Fprintf(w, OBJ_FMT, "Master", m.MarathonConfig.Master)
	fmt.Fprintf(w, OBJ_VFMT, "Failover_Timeout", m.MarathonConfig.FailoverTimeout)
	fmt.Fprintf(w, OBJ_VFMT, "LocalPort_Min", m.MarathonConfig.LocalPortMin)
	fmt.Fprintf(w, OBJ_VFMT, "LocalPort_Max", m.MarathonConfig.LocalPortMax)
	fmt.Fprintln(w, " \t \nZOOKEEPER CONFIG\t ")
	fmt.Fprintf(w, OBJ_FMT, "ZK", m.ZookeeperConfig.Zk)
	fmt.Fprintf(w, OBJ_VFMT, "ZK_Timeout", m.ZookeeperConfig.ZkTimeout)

	cli.FlushWriter(w)
	return nil
}

func (m MarLeaderInfo) ToColumns(output io.Writer) error {
	w := cli.NewTabWriter(output)
	fmt.Fprintln(w, "")
	fmt.Fprintf(w, OBJ_FMT, "Leader", m.Leader)
	cli.FlushWriter(w)
	return nil
}

func (m MarPing) ToColumns(output io.Writer) error {
	w := cli.NewTabWriter(output)
	fmt.Fprintln(w, "\nHOST\tDURATION")
	fmt.Fprintf(w, "%s\t%v ms\n", m.Host, (m.Elapsed.Nanoseconds() / int64(time.Millisecond)))
	cli.FlushWriter(w)
	return nil
}

func (m Message) ToColumns(output io.Writer) error {
	w := cli.NewTabWriter(output)
	fmt.Fprintln(w, "")
	fmt.Fprintf(w, OBJ_FMT, "Message", m.Message.Message)
	cli.FlushWriter(w)
	return nil
}

func writeMapColumn(w *tabwriter.Writer, name string, m map[string]string) {
	var i int = 0
	for k, v := range m {
		if i == 0 {
			fmt.Fprintf(w, OBJ_KVFMT, name, k, v)
		} else {
			fmt.Fprintf(w, OBJ_VAL_KVFMT, k, v)
		}
		i++
	}
}
