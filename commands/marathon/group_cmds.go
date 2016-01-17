package marathon

import (
	"fmt"
	"github.com/gondor/depcon/marathon"
	"github.com/gondor/depcon/pkg/cli"
	"github.com/spf13/cobra"
	"strings"
)

var groupCmd = &cobra.Command{
	Use:   "group",
	Short: "Marathon application groups",
	Long: `Manage application groups in a marathon cluster (eg. creating, listing, details)

    See group's subcommands for available choices`,
}

var groupListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all applications",
	Run:   listGroups,
}

var groupGetCmd = &cobra.Command{
	Use:   "get [groupId]",
	Short: "Get a group details by [groupId]",
	Run:   getGroup,
}

var groupDestroyCmd = &cobra.Command{
	Use:   "destroy [groupId]",
	Short: "Removes a group by [groupId] and all of it's resources (nested groups and app instances)",
	Run:   destroyGroup,
}

var groupCreateCmd = &cobra.Command{
	Use:   "create [file(.json | .yaml)]",
	Short: "Create a new Group with the [file(.json | .yaml)]",
	Run:   createGroup,
}

func init() {
	groupCmd.AddCommand(groupListCmd, groupGetCmd, groupCreateCmd, groupDestroyCmd)

	// Destroy Flags
	groupDestroyCmd.Flags().BoolP(WAIT_FLAG, "w", false, "Wait for destroy to complete")
	// Create Flags
	groupCreateCmd.Flags().BoolP(WAIT_FLAG, "w", false, "Wait for group to become healthy")
	groupCreateCmd.Flags().BoolP(FORCE_FLAG, "f", false, "Force deployment (updates group if it already exists)")
	groupCreateCmd.Flags().BoolP(IGNORE_MISSING, "i", false, `Ignore missing ${PARAMS} that are declared in app config that could not be resolved
                        CAUTION: This can be dangerous if some params define versions or other required information.`)
	groupCreateCmd.Flags().StringSliceP(PARAMS_FLAG, "p", nil, `Adds a param(s) that can be used for substitution.
                  eg. -p MYVAR=value would replace ${MYVAR} with "value" in the application file.
                  These take precidence over env vars and params in file`)

}

func listGroups(cmd *cobra.Command, args []string) {
	v, e := client(cmd).ListGroups()
	cli.Output(AppGroups{v}, e)
}

func getGroup(cmd *cobra.Command, args []string) {
	if cli.EvalPrintUsage(cmd.Usage, args, 1) {
		return
	}

	v, e := client(cmd).GetGroup(args[0])
	cli.Output(AppGroup{v}, e)
}

func destroyGroup(cmd *cobra.Command, args []string) {
	if cli.EvalPrintUsage(cmd.Usage, args, 1) {
		return
	}
	v, e := client(cmd).DestroyGroup(args[0])
	cli.Output(DeploymentId{v}, e)
}

func createGroup(cmd *cobra.Command, args []string) {
	if cli.EvalPrintUsage(cmd.Usage, args, 1) {
		return
	}

	wait, _ := cmd.Flags().GetBool(WAIT_FLAG)
	force, _ := cmd.Flags().GetBool(FORCE_FLAG)
	params, _ := cmd.Flags().GetStringSlice(PARAMS_FLAG)
	ignore, _ := cmd.Flags().GetBool(IGNORE_MISSING)
	options := &marathon.CreateOptions{Wait: wait, Force: force, ErrorOnMissingParams: !ignore}

	if params != nil {
		envmap := make(map[string]string)
		for _, p := range params {
			if strings.Contains(p, "=") {
				v := strings.Split(p, "=")
				envmap[v[0]] = v[1]
			}
		}
		options.EnvParams = envmap
	}

	result, e := client(cmd).CreateGroupFromFile(args[0], options)
	if e != nil && e == marathon.ErrorGroupExists {
		cli.Output(nil, fmt.Errorf("%s, consider using the --force flag to update when group exists", e.Error()))
		return
	}
	cli.Output(AppGroup{result}, e)
}
