package marathon

import (
	"github.com/gondor/depcon/pkg/cli"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Marathon server information",
	Long: `View Marathon server information and current leader

    See server's subcommands for available choices`,
}

var serverInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get info about the Marathon Instance",
	Run: func(cmd *cobra.Command, args []string) {
		v, e := client(cmd).GetMarathonInfo()
		cli.Output(MarInfo{v}, e)
	},
}

var serverLeaderCmd = &cobra.Command{
	Use:   "leader",
	Short: "Marathon leader management",
}

var serverLeaderGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Show the current leader",
	Run: func(cmd *cobra.Command, args []string) {
		v, e := client(cmd).GetCurrentLeader()
		cli.Output(MarLeaderInfo{v}, e)
	},
}

var serverPingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Ping the current marathon host",
	Run: func(cmd *cobra.Command, args []string) {
		v, e := client(cmd).Ping()
		cli.Output(MarPing{v}, e)
	},
}

var serverLeaderAbdicateCmd = &cobra.Command{
	Use:   "abdicate",
	Short: "Force the current leader to relinquish control (elect a new leader)",
	Run: func(cmd *cobra.Command, args []string) {
		v, e := client(cmd).AbdicateLeader()
		cli.Output(Message{v}, e)
	},
}

func init() {
	serverLeaderCmd.AddCommand(serverLeaderGetCmd, serverLeaderAbdicateCmd)
	serverCmd.AddCommand(serverInfoCmd, serverLeaderCmd, serverPingCmd)
}
