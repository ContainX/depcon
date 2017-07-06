package marathon

import (
	"github.com/ContainX/depcon/cliconfig"
	"github.com/ContainX/depcon/marathon"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	WAIT_FLAG      string = "wait"
	TIMEOUT_FLAG   string = "wait-timeout"
	FORCE_FLAG     string = "force"
	DETAIL_FLAG    string = "detail"
	PARAMS_FLAG    string = "param"
	ENV_FILE_FLAG  string = "env-file"
	IGNORE_MISSING string = "ignore"
	INSECURE_FLAG  string = "insecure"
	ENV_NAME       string = "env_name"
	DRYRUN_FLAG    string = "dry-run"
)

var (
	marathonCmd = &cobra.Command{
		Use:   "mar",
		Short: "Manage apache marathon services",
		Long: `Manage apache marathon services (eg. apps, deployments, tasks)

    See marathon's subcommands for available choices`,
	}
	marathonClient marathon.Marathon
	configFile     *cliconfig.ConfigFile
)

// Associates the marathon service to the given command
func AddMarathonToCmd(rc *cobra.Command, c *cliconfig.ConfigFile) {
	configFile = c
	associateServiceCommands(marathonCmd)
	rc.AddCommand(marathonCmd)
}

func AddToMarathonCommand(child *cobra.Command) {
	marathonCmd.AddCommand(child)
}

// Jails (chroots) marathon by including only it's sub commands
// when we only have a single environment declared and already know the cluster type
func AddJailedMarathonToCmd(rc *cobra.Command, c *cliconfig.ConfigFile) {
	configFile = c
	associateServiceCommands(rc)
}

// Associates all marathon service commands to specified parent
func associateServiceCommands(parent *cobra.Command) {
	parent.PersistentFlags().Bool(INSECURE_FLAG, false, "Skips Insecure TLS/HTTPS Certificate checks")
	viper.BindPFlag(INSECURE_FLAG, parent.PersistentFlags().Lookup(INSECURE_FLAG))

	parent.AddCommand(appCmd, groupCmd, deployCmd, taskCmd, eventCmd, serverCmd)
}

func client(c *cobra.Command) marathon.Marathon {
	if marathonClient == nil {
		envName := viper.GetString(ENV_NAME)
		insecure := viper.GetBool(INSECURE_FLAG)
		mc := *configFile.Environments[envName].Marathon
		opts := &marathon.MarathonOptions{}
		if timeout, err := c.Flags().GetDuration(TIMEOUT_FLAG); err == nil {
			opts.WaitTimeout = timeout
		}
		opts.TLSAllowInsecure = insecure

		marathonClient = marathon.NewMarathonClientWithOpts(mc.HostUrl, mc.Username, mc.Password, mc.Token, opts)

	}
	return marathonClient
}

func Usage(c *cobra.Command) func() error {

	return func() error {
		return c.UsageFunc()(c)
	}
}
