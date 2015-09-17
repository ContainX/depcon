package marathon

import (
	"github.com/spf13/cobra"

)

var eventCmd = &cobra.Command{
	Use:   "event",
	Short: "Marathon event streaming and subscription management",
	Long: `Manage subscriptions and listen to live streaming

    See events's subcommands for available choices`,
}

func init() {
}
