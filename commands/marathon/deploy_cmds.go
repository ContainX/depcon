package marathon

import (
	"github.com/ContainX/depcon/pkg/cli"
	"github.com/spf13/cobra"
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
		if cli.EvalPrintUsage(cmd.Usage, args, 1) {
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
		if cli.EvalPrintUsage(cmd.Usage, args, 1) {
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

	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {

	deployCreateCmd.Flags().String(TEMPLATE_CTX_FLAG, "", "Provides data per environment in JSON form to do a first pass parse of descriptor as template")
	deployCreateCmd.Flags().BoolP(FORCE_FLAG, "f", false, "Force deployment (updates application if it already exists)")
	deployCreateCmd.Flags().Bool(STOP_DEPLOYS_FLAG, false, "Stop an existing deployment for this app (if exists) and use this revision")
	deployCreateCmd.Flags().BoolP(IGNORE_MISSING, "i", false, `Ignore missing ${PARAMS} that are declared in app config that could not be resolved
                        CAUTION: This can be dangerous if some params define versions or other required information.`)
	deployCreateCmd.Flags().StringP(ENV_FILE_FLAG, "c", "", `Adds a file with a param(s) that can be used for substitution.
						These take precidence over env vars`)
	deployCreateCmd.Flags().StringSliceP(PARAMS_FLAG, "p", nil, `Adds a param(s) that can be used for substitution.
                  eg. -p MYVAR=value would replace ${MYVAR} with "value" in the application file.
                  These take precidence over env vars`)

	deployDeleteCmd.Flags().BoolP(FORCE_FLAG, "f", false, "If set to true, then the deployment is still canceled but no rollback deployment is created.")
	deployCmd.AddCommand(deployListCmd, deployDeleteCmd, deleteIfDeployingCmd)
}
