package marathon

import (
	"bytes"
	"fmt"
	"github.com/ContainX/depcon/marathon"
	"github.com/ContainX/depcon/pkg/cli"
	"github.com/ContainX/depcon/pkg/encoding"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"time"
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
	groupCreateCmd.Flags().String(TEMPLATE_CTX_FLAG, DEFAULT_CTX, "Provides data per environment in JSON form to do a first pass parse of descriptor as template")
	groupCreateCmd.Flags().BoolP(WAIT_FLAG, "w", false, "Wait for group to become healthy")
	groupCreateCmd.Flags().Bool(STOP_DEPLOYS_FLAG, false, "Stop an existing deployment for this group (if exists) and use this revision")
	groupCreateCmd.Flags().BoolP(FORCE_FLAG, "f", false, "Force deployment (updates group if it already exists)")
	groupCreateCmd.Flags().DurationP(TIMEOUT_FLAG, "t", time.Duration(0), "Max duration to wait for application health (ex. 90s | 2m). See docs for ordering")

	groupCreateCmd.Flags().BoolP(IGNORE_MISSING, "i", false, `Ignore missing ${PARAMS} that are declared in app config that could not be resolved
                        CAUTION: This can be dangerous if some params define versions or other required information.`)
	groupCreateCmd.Flags().StringSliceP(PARAMS_FLAG, "p", nil, `Adds a param(s) that can be used for substitution.
                  eg. -p MYVAR=value would replace ${MYVAR} with "value" in the application file.
                  These take precidence over env vars and params in file`)

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
	if cli.EvalPrintUsage(cmd.Usage, args, 1) {
		return
	}

	v, e := client(cmd).GetGroup(args[0])
	arr := flattenGroup(v, []*marathon.Group{})
	cli.Output(templateFor(T_GROUPS, arr), e)
}

func destroyGroup(cmd *cobra.Command, args []string) {
	if cli.EvalPrintUsage(cmd.Usage, args, 1) {
		return
	}
	v, e := client(cmd).DestroyGroup(args[0])
	cli.Output(templateFor(T_DEPLOYMENT_ID, v), e)
}

func createGroup(cmd *cobra.Command, args []string) {
	if cli.EvalPrintUsage(cmd.Usage, args, 1) {
		return
	}

	wait, _ := cmd.Flags().GetBool(WAIT_FLAG)
	force, _ := cmd.Flags().GetBool(FORCE_FLAG)
	params, _ := cmd.Flags().GetStringSlice(PARAMS_FLAG)
	ignore, _ := cmd.Flags().GetBool(IGNORE_MISSING)
	stop_deploy, _ := cmd.Flags().GetBool(STOP_DEPLOYS_FLAG)

	tempctx, _ := cmd.Flags().GetString(TEMPLATE_CTX_FLAG)
	options := &marathon.CreateOptions{Wait: wait, Force: force, ErrorOnMissingParams: !ignore, StopDeploy: stop_deploy}

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

	var result *marathon.Group = nil
	var e error

	if TemplateExists(tempctx) {
		b := &bytes.Buffer{}

		r, err := LoadTemplateContext(tempctx)
		if err != nil {
			exitWithError(err)
		}

		if err := r.Transform(b, args[0]); err != nil {
			exitWithError(err)
		}
		result, e = client(cmd).CreateGroupFromString(args[0], b.String(), options)
	} else {
		result, e = client(cmd).CreateGroupFromFile(args[0], options)
	}

	if e != nil {
		if e == marathon.ErrorGroupExists {
			cli.Output(nil, fmt.Errorf("%s, consider using the --force flag to update when group exists", e.Error()))
		} else {
			cli.Output(nil, e)
		}
		return

	}
	arr := flattenGroup(result, []*marathon.Group{})
	cli.Output(templateFor(T_GROUPS, arr), e)
}

func flattenGroup(g *marathon.Group, arr []*marathon.Group) []*marathon.Group {
	arr = append(arr, g)
	for _, cg := range g.Groups {
		arr = flattenGroup(cg, arr)
	}
	return arr
}

func convertGroupFile(cmd *cobra.Command, args []string) {
	if cli.EvalPrintUsage(cmd.Usage, args, 2) {
		os.Exit(1)
	}
	if err := encoding.ConvertFile(args[0], args[1], &marathon.Groups{}); err != nil {
		cli.Output(nil, err)
		os.Exit(1)
	}
	fmt.Printf("Source file %s has been re-written into new format in %s\n\n", args[0], args[1])
}
