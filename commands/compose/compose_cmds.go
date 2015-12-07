package compose

import (
	"github.com/gondor/depcon/cliconfig"
	"github.com/gondor/depcon/compose"
	"github.com/spf13/cobra"
)

const (
	COMPOSE_FILE_FLAG string = "compose-file"
	PROJECT_NAME_FLAG string = "name"
	PARAMS_FLAG       string = "param"
	QUIET_FLAG        string = "quiet"
	INDEX_FLAG        string = "index"
	PROTO_FLAG        string = "protocol"
	IGNORE_MISSING    string = "ignore"
)

var (
	composeCmd = &cobra.Command{
		Use:   "compose",
		Short: "Docker Compose support with param substition",
		Long: `Using libcompose to manage docker compose files with param substition support

    See compose's subcommands for available choices`,
	}

	upCmd = &cobra.Command{
		Use:   "up [services ...]",
		Short: "Create and start containers",
		Run:   execAction(up),
	}

	killCmd = &cobra.Command{
		Use:   "kill [services ...]",
		Short: "Kill containers",
		Run:   execAction(kill),
	}

	logCmd = &cobra.Command{
		Use:   "logs [services ...]",
		Short: "View output from containers",
		Run:   execAction(logs),
	}

	rmCmd = &cobra.Command{
		Use:   "rm [services ...]",
		Short: "Remove stopped containers",
		Run:   execAction(delete),
	}

	buildCmd = &cobra.Command{
		Use:   "build [services ...]",
		Short: "Build or rebuild services",
		Run:   execAction(build),
	}

	psCmd = &cobra.Command{
		Use:   "ps",
		Short: "List containers",
		Run:   execAction(ps),
	}

	restartCmd = &cobra.Command{
		Use:   "restart [services ...]",
		Short: "Restart running containers",
		Run:   execAction(restart),
	}

	pullCmd = &cobra.Command{
		Use:   "pull [services ...]",
		Short: "Pulls service imagess",
		Run:   execAction(pull),
	}

	startCmd = &cobra.Command{
		Use:   "start [services ...]",
		Short: "Start services",
		Run:   execAction(start),
	}

	stopCmd = &cobra.Command{
		Use:   "stop [services ...]",
		Short: "Stops services",
		Run:   execAction(stop),
	}

	portCmd = &cobra.Command{
		Use:   "port [service] [private_port]",
		Short: "Stops services",
		Run:   execAction(port),
	}
)

// Associates the compose service to the given command
func AddComposeToCmd(rc *cobra.Command, c *cliconfig.ConfigFile) {
	rc.AddCommand(composeCmd)
}

func init() {
	composeCmd.PersistentFlags().BoolP(IGNORE_MISSING, "i", false, `Ignore missing ${PARAMS} that are declared in app config that could not be resolved
                        CAUTION: This can be dangerous if some params define versions or other required information.`)
	composeCmd.PersistentFlags().StringSliceP(PARAMS_FLAG, "p", nil, `Adds a param(s) that can be used for substitution.
                  eg. -p MYVAR=value would replace ${MYVAR} with "value" in the compose file.
                  These take precidence over env vars`)

	psCmd.Flags().BoolP(QUIET_FLAG, "q", false, "Only display IDs")
	portCmd.Flags().Int(INDEX_FLAG, 1, "index of the container if there are multiple instances of a service [default: 1]")
	portCmd.Flags().String(PROTO_FLAG, "tcp", "tcp or udp [default: tcp]")

	composeCmd.PersistentFlags().String(COMPOSE_FILE_FLAG, "docker-compose.yml", "Docker compose file")
	composeCmd.PersistentFlags().String(PROJECT_NAME_FLAG, compose.DEFAULT_PROJECT, "Project name for this composition")
	composeCmd.AddCommand(buildCmd, killCmd, logCmd, portCmd, psCmd, upCmd, pullCmd, restartCmd, rmCmd, startCmd, stopCmd, upCmd)
}
