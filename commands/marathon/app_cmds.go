package marathon

import (
	"errors"
	"fmt"
	"github.com/ContainX/depcon/marathon"
	"github.com/ContainX/depcon/pkg/cli"
	"github.com/ContainX/depcon/pkg/encoding"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	HOST_FLAG  = "host"
	SCALE_FLAG = "scale"
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
	Run:   createApp,
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
	Use:   "list",
	Short: "List all applications",
	Run: func(cmd *cobra.Command, args []string) {
		v, e := client(cmd).ListApplications()
		cli.Output(templateFor(T_APPLICATIONS, v), e)
	},
}

var appGetCmd = &cobra.Command{
	Use:   "get [applicationId]",
	Short: "Gets an application details by Id",
	Long:  `Retrieves the specified [appliationId] application`,
	Run: func(cmd *cobra.Command, args []string) {
		if cli.EvalPrintUsage(cmd.Usage, args, 1) {
			return
		}
		v, e := client(cmd).GetApplication(args[0])
		cli.Output(templateFor(T_APPLICATION, v), e)
	},
}

var appVersionsCmd = &cobra.Command{
	Use:   "versions [applicationId]",
	Short: "Gets the versions that have been deployed with Marathon for [applicationId]",
	Long:  `Retrieves the list of versions for [appliationId] application`,
	Run: func(cmd *cobra.Command, args []string) {
		if cli.EvalPrintUsage(cmd.Usage, args, 1) {
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
	appCmd.AddCommand(appListCmd, appGetCmd, logCmd, appCreateCmd, appUpdateCmd, appDestroyCmd, appRollbackCmd, bgCmd, appRestartCmd, appScaleCmd, appVersionsCmd, appConvertFileCmd)

	// Create Flags
	appCreateCmd.Flags().BoolP(FORCE_FLAG, "f", false, "Force deployment (updates application if it already exists)")
	appCreateCmd.Flags().BoolP(IGNORE_MISSING, "i", false, `Ignore missing ${PARAMS} that are declared in app config that could not be resolved
                        CAUTION: This can be dangerous if some params define versions or other required information.`)
	appCreateCmd.Flags().StringP(ENV_FILE_FLAG, "c", "", `Adds a file with a param(s) that can be used for substitution.
						These take precidence over env vars`)
	appCreateCmd.Flags().StringSliceP(PARAMS_FLAG, "p", nil, `Adds a param(s) that can be used for substitution.
                  eg. -p MYVAR=value would replace ${MYVAR} with "value" in the application file.
                  These take precidence over env vars`)

	applyCommonAppFlags(appCreateCmd, appUpdateCPUCmd, appUpdateMemoryCmd, appRollbackCmd, appDestroyCmd, appRestartCmd, appScaleCmd)
}

func createApp(cmd *cobra.Command, args []string) {
	if cli.EvalPrintUsage(cmd.Usage, args, 1) {
		return
	}

	wait, _ := cmd.Flags().GetBool(WAIT_FLAG)
	force, _ := cmd.Flags().GetBool(FORCE_FLAG)
	paramsFile, _ := cmd.Flags().GetString(ENV_FILE_FLAG)
	params, _ := cmd.Flags().GetStringSlice(PARAMS_FLAG)
	ignore, _ := cmd.Flags().GetBool(IGNORE_MISSING)
	options := &marathon.CreateOptions{Wait: wait, Force: force, ErrorOnMissingParams: !ignore}

	if paramsFile != "" {
		envParams, _ := parseParamsFile(paramsFile)
		options.EnvParams = envParams
	} else {
		options.EnvParams = make(map[string]string)
	}

	if params != nil {
		for _, p := range params {
			if strings.Contains(p, "=") {
				v := strings.Split(p, "=")
				options.EnvParams[v[0]] = v[1]
			}
		}
	}

	result, e := client(cmd).CreateApplicationFromFile(args[0], options)
	if e != nil && e == marathon.ErrorAppExists {
		cli.Output(nil, errors.New(fmt.Sprintf("%s, consider using the --force flag to update when an application exists", e.Error())))
		os.Exit(1)
	}
	if result == nil {
		if e != nil {

			fmt.Printf("[ERROR] %s\n", e.Error())
		}
		os.Exit(1)
	}
	cli.Output(templateFor(T_APPLICATIONS, result), e)
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
	if cli.EvalPrintUsage(cmd.Usage, args, 1) {
		os.Exit(1)
	}

	force, _ := cmd.Flags().GetBool(FORCE_FLAG)

	v, e := client(cmd).RestartApplication(args[0], force)
	cli.Output(templateFor(T_DEPLOYMENT_ID, v), e)
	waitForDeploymentIfFlagged(cmd, v.DeploymentID)
}

func destroyApp(cmd *cobra.Command, args []string) {
	if cli.EvalPrintUsage(cmd.Usage, args, 1) {
		os.Exit(1)
	}

	v, e := client(cmd).DestroyApplication(args[0])
	cli.Output(templateFor(T_DEPLOYMENT_ID, v), e)
	waitForDeploymentIfFlagged(cmd, v.DeploymentID)
}

func scaleApp(cmd *cobra.Command, args []string) {
	if cli.EvalPrintUsage(cmd.Usage, args, 2) {
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

func updateAppCPU(cmd *cobra.Command, args []string) {
	if cli.EvalPrintUsage(cmd.Usage, args, 2) {
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
	cli.Output(templateFor(T_APPLICATIONS, v), e)
}

func updateAppMemory(cmd *cobra.Command, args []string) {
	if cli.EvalPrintUsage(cmd.Usage, args, 2) {
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
	cli.Output(templateFor(T_APPLICATIONS, v), e)
}

func rollbackAppVersion(cmd *cobra.Command, args []string) {
	if cli.EvalPrintUsage(cmd.Usage, args, 1) {
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
	cli.Output(templateFor(T_APPLICATIONS, v), e)
}

func convertFile(cmd *cobra.Command, args []string) {
	if cli.EvalPrintUsage(cmd.Usage, args, 2) {
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
