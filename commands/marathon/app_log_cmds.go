package marathon

import (
	"fmt"
	"github.com/ContainX/depcon/pkg/cli"
	"github.com/ContainX/depcon/pkg/logger"
	ml "github.com/ContainX/go-mesoslog/mesoslog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/url"
	"strings"
)

const (
	STDERR_FLAG = "stderr"
	FOLLOW_FLAG = "follow"
	POLL_FLAG   = "poll"
)

var logCmd = &cobra.Command{
	Use:   "log [appId]",
	Short: "Log or Tail Mesos application logs",
	Long:  "Log or Tail Mesos application logs",
	Run:   showLogCmd,
}

var log = logger.GetLogger("depcon")

func init() {
	logCmd.Flags().BoolP(STDERR_FLAG, "s", false, "Show StdErr vs default StdOut log")
	logCmd.Flags().BoolP(FOLLOW_FLAG, "f", false, "Tail/Follow log")
	logCmd.Flags().IntP(POLL_FLAG, "p", 5, "Log poll time (duration) in seconds")
}

func showLogCmd(cmd *cobra.Command, args []string) {
	if cli.EvalPrintUsage(Usage(cmd), args, 1) {
		return
	}

	host := getMesosHost()
	logType := ml.STDOUT

	if stderr, _ := cmd.Flags().GetBool(STDERR_FLAG); stderr {
		logType = ml.STDERR
	}

	c, _ := ml.NewMesosClient(host, 5050)
	appId := getMesosAppIdentifier(cmd, c, args[0])

	if follow, _ := cmd.Flags().GetBool(FOLLOW_FLAG); follow {
		duration, _ := cmd.Flags().GetInt(POLL_FLAG)
		if duration < 1 {
			duration = 5

		}
		if err := c.TailLog(appId, logType, duration); err != nil {
			log.Fatal(err)
		}
		return
	}

	logs, err := c.GetLog(appId, logType, "")
	if err != nil {
		log.Fatal(err)
	}

	showBreaks := len(logs) > 1
	for _, log := range logs {
		if showBreaks {
			fmt.Printf("\n::: [ %s - Logs For: %s ] ::: \n", args[0], log.TaskID)
		}
		fmt.Printf("%s\n", log.Log)
		if showBreaks {
			fmt.Printf("\n!!! [ %s - End Logs For: %s ] !!! \n", args[0], log.TaskID)
		}
	}
}

func getMesosAppIdentifier(cmd *cobra.Command, c *ml.MesosClient, appId string) string {
	tasks, err := client(cmd).GetTasks(appId)
	if err != nil {
		log.Fatal(err)
	}

	if len(tasks) > 0 {
		name, err := c.GetAppNameForTaskID(tasks[0].ID)
		if err != nil {
			log.Fatal(err)
		}
		return name
	}

	log.Fatal("Currently no tasks found for application: %s\n", appId)
	return ""
}

func getMesosHost() string {
	envName := viper.GetString("env_name")
	mc := *configFile.Environments[envName].Marathon

	u, err := url.Parse(mc.HostUrl)
	if err != nil {
		log.Fatal(err)
	}
	if strings.Index(u.Host, ":") > 0 {
		return strings.Split(u.Host, ":")[0]
	}
	return u.Host
}
