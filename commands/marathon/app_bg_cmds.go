package marathon

import (
	"github.com/gondor/depcon/marathon/bluegreen"
	"github.com/gondor/depcon/pkg/cli"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"time"
)

const (
	INSTANCES_FLAG  = "instances"
	STEP_DELAY_FLAG = "stepdel"
	RESUME_FLAG     = "resume"
	LB_FLAG         = "lb"
	LB_TIMEOUT_FLAG = "lb-timeout"
	BG_DRYRUN_FLAG  = "dry"
)

var bgCmd = &cobra.Command{
	Use:   "bluegreen [file(.json | .yaml)]",
	Short: "Marathon blue/green deployments",
	Long: `Blue/Green deployments handled through HAProxy or Marathon-LB

    See bluegreen's subcommands for available choices`,
	Run: deployBlueGreenCmd,
}

func init() {
	bgCmd.Flags().String(LB_FLAG, "http://localhost:9090", "HAProxy URL and Stats Port")
	bgCmd.Flags().Int(LB_TIMEOUT_FLAG, 300, "HAProxy timeout - default 300 seconds")
	bgCmd.Flags().Int(INSTANCES_FLAG, 1, "Initial intances of the app to create")
	bgCmd.Flags().Int(STEP_DELAY_FLAG, 6, "Delay (in seconds) to wait between successive deployment steps. ")
	bgCmd.Flags().Bool(RESUME_FLAG, true, "Resume from a previous deployment")
	bgCmd.Flags().BoolP(IGNORE_MISSING, "i", false, `Ignore missing ${PARAMS} that are declared in app config that could not be resolved
                        CAUTION: This can be dangerous if some params define versions or other required information.`)
	bgCmd.Flags().StringP(ENV_FILE_FLAG, "c", "", `Adds a file with a param(s) that can be used for substitution.
		   These take precidence over env vars`)
	bgCmd.Flags().StringSliceP(PARAMS_FLAG, "p", nil, `Adds a param(s) that can be used for substitution.
                  eg. -p MYVAR=value would replace ${MYVAR} with "value" in the application file.
                  These take precidence over env vars`)
	bgCmd.Flags().Bool(BG_DRYRUN_FLAG, false, "Dry run (no deployment or scaling)")

}

func deployBlueGreenCmd(cmd *cobra.Command, args []string) {

	a, err := bgc(cmd).DeployBlueGreenFromFile(args[0])
	if err != nil {
		cli.Output(nil, err)
		os.Exit(1)
	}
	cli.Output(Application{a}, err)

}

func bgc(c *cobra.Command) bluegreen.BlueGreen {

	paramsFile, _ := c.Flags().GetString(ENV_FILE_FLAG)
	params, _ := c.Flags().GetStringSlice(PARAMS_FLAG)
	ignore, _ := c.Flags().GetBool(IGNORE_MISSING)
	sd, _ := c.Flags().GetInt(STEP_DELAY_FLAG)
	lbtimeout, _ := c.Flags().GetInt(LB_TIMEOUT_FLAG)

	// Create Options
	opts := bluegreen.NewBlueGreenOptions()
	opts.Resume, _ = c.Flags().GetBool(RESUME_FLAG)
	opts.LoadBalancer, _ = c.Flags().GetString(LB_FLAG)
	opts.InitialInstances, _ = c.Flags().GetInt(INSTANCES_FLAG)
	opts.ErrorOnMissingParams = !ignore
	opts.StepDelay = time.Duration(sd) * time.Second
	opts.ProxyWaitTimeout = time.Duration(lbtimeout) * time.Second
	opts.DryRun, _ = c.Flags().GetBool(BG_DRYRUN_FLAG)

	if paramsFile != "" {
		envParams, _ := parseParamsFile(paramsFile)
		opts.EnvParams = envParams
	} else {
		opts.EnvParams = make(map[string]string)
	}

	if params != nil {
		for _, p := range params {
			if strings.Contains(p, "=") {
				v := strings.Split(p, "=")
				opts.EnvParams[v[0]] = v[1]
			}
		}
	}

	return bluegreen.NewBlueGreenClient(client(c), opts)
}
