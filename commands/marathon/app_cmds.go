package marathon

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ContainX/depcon/marathon"
	"github.com/ContainX/depcon/pkg/cli"
	"github.com/ContainX/depcon/pkg/encoding"
	"github.com/spf13/cobra"
)

const (
	HOST_FLAG         = "host"
	SCALE_FLAG        = "scale"
	FORMAT_FLAG       = "format"
	TEMPLATE_CTX_FLAG = "tempctx"
	DEFAULT_CTX       = "template-context.json"
	STOP_DEPLOYS_FLAG = "stop-deploys"
)

var appCmd = &cobra.Command{
	Use:   "app",
	Short: "Marathon application management",
	Long: `Manage applications in a marathon cluster (eg. creating, listing, details)

    See app's subcommands for available choices`,
}

var appCreateCmd = &cobra.Command{
	Use:   "create [file(.json | .yaml)]",
	Short: "Create a new application with the [file(.json | .yaml)]",
	Long:  "Creates a new App in the cluster.  This is an alias for the 'deploy create' command",
	Run:   deployAppOrGroup,
}

var appUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates a running application.  See subcommands for available choices",
}

var appUpdateCPUCmd = &cobra.Command{
	Use:   "cpu [applicationId] [cpu_shares]",
	Short: "Updates [applicationId] to have [cpu_shares]",
	Run:   updateAppCPU,
}

var appUpdateMemoryCmd = &cobra.Command{
	Use:   "mem [applicationId] [amount]",
	Short: "Updates [applicationId] to have [amount] of memory in MB",
	Run:   updateAppMemory,
}

var appListCmd = &cobra.Command{
	Use:   "list (optional filtering - label=mylabel | id=/services | cmd=java ...)",
	Short: "List all applications",
	Run: func(cmd *cobra.Command, args []string) {
		filter := ""
		if len(args) > 0 {
			filter = args[0]
		}
		v, e := client(cmd).ListApplicationsWithFilters(filter)

		cli.Output(templateFor(templateFormat(T_APPLICATIONS, cmd), v), e)
	},
}

var appGetCmd = &cobra.Command{
	Use:   "get [applicationId]",
	Short: "Gets an application details by Id",
	Long:  `Retrieves the specified [appliationId] application`,
	Run: func(cmd *cobra.Command, args []string) {
		if cli.EvalPrintUsage(Usage(cmd), args, 1) {
			return
		}
		v, e := client(cmd).GetApplication(args[0])
		cli.Output(templateFor(templateFormat(T_APPLICATION, cmd), v), e)
	},
}

var appVersionsCmd = &cobra.Command{
	Use:   "versions [applicationId]",
	Short: "Gets the versions that have been deployed with Marathon for [applicationId]",
	Long:  `Retrieves the list of versions for [appliationId] application`,
	Run: func(cmd *cobra.Command, args []string) {
		if cli.EvalPrintUsage(Usage(cmd), args, 1) {
			return
		}
		v, e := client(cmd).ListVersions(args[0])
		cli.Output(templateFor(T_VERSIONS, v), e)
	},
}

var appDestroyCmd = &cobra.Command{
	Use:   "destroy [applicationId]",
	Short: "Remove an application [applicationId] and all of it's instances",
	Long:  `Removes the specified [appliationId] application`,
	Run:   destroyApp,
}

var appRestartCmd = &cobra.Command{
	Use:   "restart [applicationId]",
	Short: "Restarts an application by Id",
	Long:  `Restarts the specified [appliationId] application`,
	Run:   restartApp,
}

var appScaleCmd = &cobra.Command{
	Use:   "scale [applicationId] [instances]",
	Short: "Scales [appliationId] to total [instances]",
	Run:   scaleApp,
}

var appPauseCmd = &cobra.Command{
	Use:   "pause [applicationId]",
	Short: "Suspends the [applicationId",
	Run:   pauseApp,
}

var appRollbackCmd = &cobra.Command{
	Use:   "rollback [applicationId] (version)",
	Short: "Rolls an [appliationId] to a specific (version : optional)",
	Long:  `Rolls an [appliationId] to a specific [version] - See: "depcon app versions" for a list of versions`,
	Run:   rollbackAppVersion,
}

var appConvertFileCmd = &cobra.Command{
	Use:   "convert [from.(json | yaml)] [to.(json | yaml)]",
	Short: "Utilty to convert an application file from json to yaml or yaml to json.",
	Run:   convertFile,
}

func init() {
	appUpdateCmd.AddCommand(appUpdateCPUCmd, appUpdateMemoryCmd)
	appCmd.AddCommand(appListCmd, appGetCmd, logCmd, appCreateCmd, appUpdateCmd, appDestroyCmd, appRollbackCmd, bgCmd, appRestartCmd, appScaleCmd, appPauseCmd, appVersionsCmd, appConvertFileCmd)

	// Create Flags
	addDeployCreateFlags(appCreateCmd)

	appListCmd.Flags().String(FORMAT_FLAG, "", "Custom output format. Example: '{{range .Apps}}{{ .Container.Docker.Image }}{{end}}'")
	appGetCmd.Flags().String(FORMAT_FLAG, "", "Custom output format. Example: '{{ .ID }}'")
	applyCommonAppFlags(appUpdateCPUCmd, appUpdateMemoryCmd, appRollbackCmd, appDestroyCmd, appRestartCmd, appScaleCmd, appPauseCmd)
}

func exitWithError(err error) {
	cli.Output(nil, err)
	os.Exit(1)
}

func parseParamsFile(filename string) (map[string]string, error) {
	paramsFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(paramsFile)
	if err != nil {
		return nil, err
	}
	data := string(bytes)
	params := strings.Split(data, "\n")

	envmap := make(map[string]string)
	for _, p := range params {
		if strings.Contains(p, "=") {
			v := strings.Split(p, "=")
			envmap[v[0]] = v[1]
		}
	}
	return envmap, nil
}

func restartApp(cmd *cobra.Command, args []string) {
	if cli.EvalPrintUsage(Usage(cmd), args, 1) {
		os.Exit(1)
	}

	force, _ := cmd.Flags().GetBool(FORCE_FLAG)

	v, e := client(cmd).RestartApplication(args[0], force)
	cli.Output(templateFor(T_DEPLOYMENT_ID, v), e)
	waitForDeploymentIfFlagged(cmd, v.DeploymentID)
}

func destroyApp(cmd *cobra.Command, args []string) {
	if cli.EvalPrintUsage(Usage(cmd), args, 1) {
		os.Exit(1)
	}

	v, e := client(cmd).DestroyApplication(args[0])
	cli.Output(templateFor(T_DEPLOYMENT_ID, v), e)
	waitForDeploymentIfFlagged(cmd, v.DeploymentID)
}

func scaleApp(cmd *cobra.Command, args []string) {
	if cli.EvalPrintUsage(Usage(cmd), args, 2) {
		os.Exit(1)
	}

	instances, err := strconv.Atoi(args[1])
	if err != nil {
		cli.Output(nil, err)
		os.Exit(1)
	}
	v, e := client(cmd).ScaleApplication(args[0], instances)
	cli.Output(templateFor(T_DEPLOYMENT_ID, v), e)
	waitForDeploymentIfFlagged(cmd, v.DeploymentID)
}

func pauseApp(cmd *cobra.Command, args []string) {
	if cli.EvalPrintUsage(Usage(cmd), args, 1) {
		os.Exit(1)
	}

	v, e := client(cmd).PauseApplication(args[0])
	cli.Output(templateFor(T_DEPLOYMENT_ID, v), e)
	waitForDeploymentIfFlagged(cmd, v.DeploymentID)
}

func updateAppCPU(cmd *cobra.Command, args []string) {
	if cli.EvalPrintUsage(Usage(cmd), args, 2) {
		os.Exit(1)
	}

	wait, _ := cmd.Flags().GetBool(WAIT_FLAG)
	cpu, err := strconv.ParseFloat(args[1], 64)

	if err != nil {
		cli.Output(nil, err)
		os.Exit(1)
	}
	update := marathon.NewApplication(args[0]).CPU(cpu)
	v, e := client(cmd).UpdateApplication(update, wait)
	cli.Output(templateFor(T_APPLICATION, v), e)
}

func updateAppMemory(cmd *cobra.Command, args []string) {
	if cli.EvalPrintUsage(Usage(cmd), args, 2) {
		os.Exit(1)
	}

	wait, _ := cmd.Flags().GetBool(WAIT_FLAG)
	mem, err := strconv.ParseFloat(args[1], 64)

	if err != nil {
		cli.Output(nil, err)
		os.Exit(1)
	}
	update := marathon.NewApplication(args[0]).Memory(mem)
	v, e := client(cmd).UpdateApplication(update, wait)
	cli.Output(templateFor(T_APPLICATION, v), e)
}

func rollbackAppVersion(cmd *cobra.Command, args []string) {
	if cli.EvalPrintUsage(Usage(cmd), args, 1) {
		os.Exit(1)
	}

	wait, _ := cmd.Flags().GetBool(WAIT_FLAG)
	version := ""

	if len(args) > 1 {
		version = args[1]
	} else {
		versions, e := client(cmd).ListVersions(args[0])
		if e == nil && len(versions.Versions) > 1 {
			version = versions.Versions[1]
		}
	}
	update := marathon.NewApplication(args[0]).RollbackVersion(version)
	v, e := client(cmd).UpdateApplication(update, wait)
	cli.Output(templateFor(T_APPLICATION, v), e)
}

func convertFile(cmd *cobra.Command, args []string) {
	if cli.EvalPrintUsage(Usage(cmd), args, 2) {
		os.Exit(1)
	}
	if err := encoding.ConvertFile(args[0], args[1], &marathon.Application{}); err != nil {
		cli.Output(nil, err)
		os.Exit(1)
	}
	fmt.Printf("Source file %s has been re-written into new format in %s\n\n", args[0], args[1])
}

func waitForDeploymentIfFlagged(cmd *cobra.Command, depId string) {
	if found, err := cmd.Flags().GetBool(WAIT_FLAG); err == nil && found {
		client(cmd).WaitForDeployment(depId, time.Duration(80)*time.Second)
	}
}

func applyCommonAppFlags(cmd ...*cobra.Command) {
	for _, c := range cmd {
		c.Flags().BoolP(WAIT_FLAG, "w", false, "Wait for application to become healthy")
		c.Flags().DurationP(TIMEOUT_FLAG, "t", time.Duration(0), "Max duration to wait for application health (ex. 90s | 2m). See docs for ordering")
	}
}

func templateFormat(template string, cmd *cobra.Command) string {
	t := template
	tv, _ := cmd.Flags().GetString(FORMAT_FLAG)
	if len(tv) > 0 {
		t = tv
	}
	return t
}
