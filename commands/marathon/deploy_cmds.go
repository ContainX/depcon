package marathon

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/ContainX/depcon/marathon"
	"github.com/ContainX/depcon/pkg/cli"
	"github.com/ContainX/depcon/pkg/encoding"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Marathon deployment management",
	Long: `Manage deployments in a marathon cluster (eg. creating, listing, monitoring)

    See deploy's subcommands for available choices`,
}

var deployListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all deployments",
	Run: func(cmd *cobra.Command, args []string) {
		v, e := client(cmd).ListDeployments()
		cli.Output(templateFor(T_DEPLOYMENTS, v), e)
	},
}

var deployDeleteCmd = &cobra.Command{
	Use:   "delete [deploymentId]",
	Short: "Delete a deployment by [deploymentID]",
	Run: func(cmd *cobra.Command, args []string) {
		if cli.EvalPrintUsage(Usage(cmd), args, 1) {
			return
		}
		force, _ := cmd.Flags().GetBool(FORCE_FLAG)

		v, e := client(cmd).DeleteDeployment(args[0], force)
		cli.Output(templateFor(T_DEPLOYMENT_ID, v), e)
	},
}

var deleteIfDeployingCmd = &cobra.Command{
	Use:   "cancel-app [appid]",
	Short: "Conditional Match: Delete a deployment based on the specified [appid]",
	Run: func(cmd *cobra.Command, args []string) {
		if cli.EvalPrintUsage(Usage(cmd), args, 1) {
			return
		}
		v, e := client(cmd).CancelAppDeployment(args[0], false)
		if v != nil || e != nil {
			cli.Output(templateFor(T_DEPLOYMENT_ID, v), e)
		} else {
			if e != nil {
				cli.Output(templateFor(T_DEPLOYMENT_ID, v), e)
			}
		}

	},
}

var deployCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a new app or group by introspecting the incoming descriptor.  Useful for deployment pipelines",
	Long: `
Creates a new app or group by introspecting the incoming descriptor.  Useful for deployment pipelines.

This command has a small penalty of unmarshalling the descriptor twice. One for introspection and
the other for delegation to the origin (app or group)
	`,

	Run: deployAppOrGroup,
}

func init() {

	deployCreateCmd.Flags().BoolP(WAIT_FLAG, "w", false, "Wait for group to become healthy")
	deployCreateCmd.Flags().String(TEMPLATE_CTX_FLAG, DEFAULT_CTX, "Provides data per environment in JSON form to do a first pass parse of descriptor as template")
	deployCreateCmd.Flags().BoolP(FORCE_FLAG, "f", false, "Force deployment (updates application if it already exists)")
	deployCreateCmd.Flags().Bool(STOP_DEPLOYS_FLAG, false, "Stop an existing deployment for this app (if exists) and use this revision")
	deployCreateCmd.Flags().BoolP(IGNORE_MISSING, "i", false, `Ignore missing ${PARAMS} that are declared in app config that could not be resolved
                        CAUTION: This can be dangerous if some params define versions or other required information.`)
	deployCreateCmd.Flags().StringP(ENV_FILE_FLAG, "c", "", `Adds a file with a param(s) that can be used for substitution.
						These take precidence over env vars`)
	deployCreateCmd.Flags().StringSliceP(PARAMS_FLAG, "p", nil, `Adds a param(s) that can be used for substitution.
                  eg. -p MYVAR=value would replace ${MYVAR} with "value" in the application file.
                  These take precidence over env vars`)

	deployCreateCmd.Flags().Bool(DRYRUN_FLAG, false, "Preview the parsed template - don't actually deploy")

	deployCreateCmd.Flags().DurationP(TIMEOUT_FLAG, "t", time.Duration(0), "Max duration to wait for application health (ex. 90s | 2m). See docs for ordering")
	deployDeleteCmd.Flags().BoolP(FORCE_FLAG, "f", false, "If set to true, then the deployment is still canceled but no rollback deployment is created.")
	deployCmd.AddCommand(deployCreateCmd, deployListCmd, deployDeleteCmd, deleteIfDeployingCmd)
}

func deployAppOrGroup(cmd *cobra.Command, args []string) {

	if cli.EvalPrintUsage(Usage(cmd), args, 1) {
		return
	}

	filename := args[0]
	wait, _ := cmd.Flags().GetBool(WAIT_FLAG)
	force, _ := cmd.Flags().GetBool(FORCE_FLAG)
	paramsFile, _ := cmd.Flags().GetString(ENV_FILE_FLAG)
	params, _ := cmd.Flags().GetStringSlice(PARAMS_FLAG)
	ignore, _ := cmd.Flags().GetBool(IGNORE_MISSING)
	stop_deploy, _ := cmd.Flags().GetBool(STOP_DEPLOYS_FLAG)
	tempctx, _ := cmd.Flags().GetString(TEMPLATE_CTX_FLAG)
	dryrun, _ := cmd.Flags().GetBool(DRYRUN_FLAG)
	options := &marathon.CreateOptions{Wait: wait, Force: force, ErrorOnMissingParams: !ignore, StopDeploy: stop_deploy, DryRun: dryrun}

	descriptor := parseDescriptor(tempctx, filename)
	et, err := encoding.NewEncoderFromFileExt(filename)
	if err != nil {
		exitWithError(err)
	}

	ag := &marathon.AppOrGroup{}
	if err := et.UnMarshalStr(descriptor, ag); err != nil {
		exitWithError(err)
	}

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

	if ag.IsApplication() {
		result, e := client(cmd).CreateApplicationFromString(filename, descriptor, options)
		outputDeployment(result, e)
		cli.Output(templateFor(T_APPLICATION, result), e)
	} else {
		result, e := client(cmd).CreateGroupFromString(filename, descriptor, options)
		outputDeployment(result, e)

		if e != nil {
			cli.Output(nil, e)
		}

		arr := flattenGroup(result, []*marathon.Group{})
		cli.Output(templateFor(T_GROUPS, arr), e)
	}
}

func outputDeployment(result interface{}, e error) {
	if e != nil && e == marathon.ErrorAppExists {
		exitWithError(errors.New(fmt.Sprintf("%s, consider using the --force flag to update when an application exists", e.Error())))
	}

	if result == nil {
		if e != nil {
			exitWithError(e)
		}
		os.Exit(1)
	}
}

func parseDescriptor(tempctx, filename string) string {
	if TemplateExists(tempctx) {
		b := &bytes.Buffer{}

		r, err := LoadTemplateContext(tempctx)
		if err != nil {
			exitWithError(err)
		}

		if err := r.Transform(b, filename); err != nil {
			exitWithError(err)
		}
		return b.String()
	} else {
		if b, err := ioutil.ReadFile(filename); err != nil {
			exitWithError(err)
		} else {
			return string(b)

		}
	}
	return ""
}
