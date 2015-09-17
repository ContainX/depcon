package marathon

import (
	"github.com/gondor/depcon/marathon"
	"github.com/gondor/depcon/cliconfig"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	WAIT_FLAG      string = "wait"
	FORCE_FLAG     string = "force"
	DETAIL_FLAG    string = "detail"
	PARAMS_FLAG    string = "param"
	IGNORE_MISSING string = "ignore"
)

var (
	marathonCmd = &cobra.Command{
		Use:   "mar",
		Short: "Manage apache marathon services",
		Long: `Manage apache marathon services (eg. apps, deployments, tasks)

    See marathon's subcommands for available choices`,
	}
	marathonClient marathon.Marathon
	configFile *cliconfig.ConfigFile
)

// Associates the marathon service to the given command
func AddMarathonToCmd(rc *cobra.Command, c *cliconfig.ConfigFile) {
	configFile = c
	associateServiceCommands(marathonCmd)
	rc.AddCommand(marathonCmd)
}

// Jails (chroots) marathon by including only it's sub commands
// when we only have a single environment declared and already know the cluster type
func AddJailedMarathonToCmd(rc *cobra.Command, c *cliconfig.ConfigFile) {
	configFile = c
	associateServiceCommands(rc)
}

// Associates all marathon service commands to specified parent
func associateServiceCommands(parent *cobra.Command) { parent.AddCommand(appCmd, groupCmd, deployCmd, taskCmd, eventCmd, serverCmd) }

func client(c *cobra.Command) marathon.Marathon {
	if marathonClient == nil {
		envName := viper.GetString("env_name")
		mc := *configFile.Environments[envName].Marathon
		marathonClient = marathon.NewMarathonClient(mc.HostUrl, mc.Username, mc.Password)

	}
	return marathonClient
}
