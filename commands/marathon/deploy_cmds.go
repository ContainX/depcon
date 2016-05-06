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
		cli.Output(Templated{Template: T_DEPLOYMENTS, Data: v}, e)
	},
}

var deployDeleteCmd = &cobra.Command{
	Use:   "delete [deploymentId]",
	Short: "Delete a deployment by [deploymentID]",
	Run: func(cmd *cobra.Command, args []string) {
		if cli.EvalPrintUsage(cmd.Usage, args, 1) {
			return
		}
		v, e := client(cmd).DeleteDeployment(args[0])
		cli.Output(Templated{Template: T_DEPLOYMENT_ID, Data: v}, e)
	},
}

func init() {
	deployCmd.AddCommand(deployListCmd, deployDeleteCmd)
}
