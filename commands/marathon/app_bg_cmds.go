package marathon

import (
	"github.com/gondor/depcon/marathon/bluegreen"
	"github.com/gondor/depcon/pkg/cli"
	"github.com/spf13/cobra"
	"os"
	"time"
)

const (
	INSTANCES_FLAG  = "instances"
	STEP_DELAY_FLAG = "stepdel"
	RESUME_FLAG     = "resume"
	LB_FLAG         = "lb"
	LB_TIMEOUT_FLAG = "lb-timeout"
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
	opts := bluegreen.NewBlueGreenOptions()
	opts.Resume, _ = c.Flags().GetBool(RESUME_FLAG)
	opts.LoadBalancer, _ = c.Flags().GetString(LB_FLAG)

	opts.InitialInstances, _ = c.Flags().GetInt(INSTANCES_FLAG)
	ignore, _ := c.Flags().GetBool(IGNORE_MISSING)
	opts.ErrorOnMissingParams = !ignore

	sd, _ := c.Flags().GetInt(STEP_DELAY_FLAG)
	opts.StepDelay = time.Duration(sd) * time.Second

	lbtimeout, _ := c.Flags().GetInt(LB_TIMEOUT_FLAG)
	opts.ProxyWaitTimeout = time.Duration(lbtimeout) * time.Second

	return bluegreen.NewBlueGreenClient(client(c), opts)
}

/*

opts.InitialInstances = 1
	opts.ProxyWaitTimeout = time.Duration(300) * time.Second
	opts.StepDelay = time.Duration(6) * time.Second
// The max time to wait on HAProxy to drain connections (in seconds)
	ProxyWaitTimeout time.Duration
	// Initial number of app instances to create
	InitialInstances int
	// Delay (in seconds) to wait between each successive deployment step
	StepDelay time.Duration
	// Resume from previous deployment
	Resume bool
	// Marathon-LB stats endpoint - ex: http://host:9090
	LoadBalancer string
	// if true will attempt to wait until the NEW application or group is running
	Wait bool
	// If true an error will be returned on params defined in the configuration file that
	// could not resolve to user input and environment variables
	ErrorOnMissingParams bool
	// Additional environment params - looks at this map for token substitution which takes
	// priority over matching environment variables
	EnvParams map[string]string
*/
