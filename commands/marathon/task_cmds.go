package marathon

import (
	"fmt"
	"github.com/ContainX/depcon/pkg/cli"
	"github.com/spf13/cobra"
)

var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "Marathon task management",
	Long: `Manage tasks in a marathon cluster (eg. creating, listing, monitoring, kill)

    See tasks's subcommands for available choices`,
}

var taskListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks",
	Run: func(cmd *cobra.Command, args []string) {
		v, e := client(cmd).ListTasks()
		cli.Output(templateFor(T_TASKS, v), e)
	},
}

var taskQueueCmd = &cobra.Command{
	Use:   "queue",
	Short: "List all queued tasks",
	Run: func(cmd *cobra.Command, args []string) {
		v, e := client(cmd).ListQueue()
		cli.Output(templateFor(T_QUEUED_TASKS, v), e)
	},
}

var appTaskGetCmd = &cobra.Command{
	Use:   "get [applicationId]",
	Short: "List tasks for the application [applicationId]",
	Run:   appTasks,
}

var appTaskKillallCmd = &cobra.Command{
	Use:   "killall [applicationId]",
	Short: "Kill tasks belonging to [applicationId]",
	Run:   appKillAllTasks,
}

var appTaskKillCmd = &cobra.Command{
	Use:   "kill [taskId]",
	Short: "Kill a task [taskId] that belongs to an application",
	Run:   appKillTask,
}

func init() {
	taskCmd.AddCommand(taskListCmd, appTaskGetCmd, appTaskKillCmd, appTaskKillallCmd, taskQueueCmd)

	// Task List Flags
	appTaskGetCmd.Flags().BoolP(DETAIL_FLAG, "d", false, "Prints each task instance in detailed form vs. table summary")
	// Task Kill Flags
	appTaskKillallCmd.Flags().String(HOST_FLAG, "", "Kill only those tasks running on host [host]. Default: none.")
	appTaskKillallCmd.Flags().Bool(SCALE_FLAG, false, "Scale the app down (i.e. decrement its instances setting by the number of tasks killed)")
	appTaskKillCmd.Flags().Bool(SCALE_FLAG, false, "Scale the app down (i.e. decrement its instances setting by the number of tasks killed)")
}

func appTasks(cmd *cobra.Command, args []string) {
	if cli.EvalPrintUsage(Usage(cmd), args, 1) {
		return
	}

	detailed, _ := cmd.Flags().GetBool(DETAIL_FLAG)

	v, e := client(cmd).GetTasks(args[0])

	if detailed && e == nil {
		fmt.Println("")
		for _, t := range v {
			fmt.Printf("::: Task: %s\n\n", t.ID)
			cli.Output(templateFor(T_TASKS, v), e)
		}
	} else {
		cli.Output(templateFor(T_TASKS, v), e)
	}

}

func appKillAllTasks(cmd *cobra.Command, args []string) {
	if cli.EvalPrintUsage(Usage(cmd), args, 1) {
		return
	}

	host, _ := cmd.Flags().GetString(HOST_FLAG)
	scale, _ := cmd.Flags().GetBool(SCALE_FLAG)

	v, e := client(cmd).KillAppTasks(args[0], host, scale)
	cli.Output(templateFor(T_TASKS, v), e)
}

func appKillTask(cmd *cobra.Command, args []string) {
	if cli.EvalPrintUsage(Usage(cmd), args, 1) {
		return
	}
	scale, _ := cmd.Flags().GetBool(SCALE_FLAG)
	v, e := client(cmd).KillAppTask(args[0], scale)
	cli.Output(templateFor(T_TASK, v), e)
}
