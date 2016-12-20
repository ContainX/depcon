package marathon

import (
	"fmt"
	"github.com/ContainX/depcon/marathon"
	"github.com/ContainX/depcon/pkg/cli"
	"github.com/ContainX/depcon/pkg/encoding"
	"github.com/spf13/cobra"
	"os"
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
	Long:  "Creates a new Group in the cluster.  This is an alias for the 'deploy create' command",
	Run:   deployAppOrGroup,
}

var groupConvertFileCmd = &cobra.Command{
	Use:   "convert [from.(json | yaml)] [to.(json | yaml)]",
	Short: "Utilty to convert an group file from json to yaml or yaml to json.",
	Run:   convertGroupFile,
}

func init() {
	groupCmd.AddCommand(groupListCmd, groupGetCmd, groupCreateCmd, groupDestroyCmd, groupConvertFileCmd)

	// Destroy Flags
	groupDestroyCmd.Flags().BoolP(WAIT_FLAG, "w", false, "Wait for destroy to complete")
	// Create Flags
	addDeployCreateFlags(groupCreateCmd)
}

func listGroups(cmd *cobra.Command, args []string) {
	v, e := client(cmd).ListGroups()
	arr := []*marathon.Group{}

	for _, group := range v.Groups {
		arr = flattenGroup(group, arr)
	}

	cli.Output(templateFor(T_GROUPS, arr), e)
}

func getGroup(cmd *cobra.Command, args []string) {
	if cli.EvalPrintUsage(Usage(cmd), args, 1) {
		return
	}

	v, e := client(cmd).GetGroup(args[0])
	arr := flattenGroup(v, []*marathon.Group{})
	cli.Output(templateFor(T_GROUPS, arr), e)
}

func destroyGroup(cmd *cobra.Command, args []string) {
	if cli.EvalPrintUsage(Usage(cmd), args, 1) {
		return
	}
	v, e := client(cmd).DestroyGroup(args[0])
	cli.Output(templateFor(T_DEPLOYMENT_ID, v), e)
}

func flattenGroup(g *marathon.Group, arr []*marathon.Group) []*marathon.Group {
	arr = append(arr, g)
	for _, cg := range g.Groups {
		arr = flattenGroup(cg, arr)
	}
	return arr
}

func convertGroupFile(cmd *cobra.Command, args []string) {
	if cli.EvalPrintUsage(Usage(cmd), args, 2) {
		os.Exit(1)
	}
	if err := encoding.ConvertFile(args[0], args[1], &marathon.Groups{}); err != nil {
		cli.Output(nil, err)
		os.Exit(1)
	}
	fmt.Printf("Source file %s has been re-written into new format in %s\n\n", args[0], args[1])
}
