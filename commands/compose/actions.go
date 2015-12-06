package compose

import (
	"errors"
	"github.com/gondor/depcon/compose"
	"github.com/gondor/depcon/pkg/cli"
	"github.com/spf13/cobra"
)

var (
	PortInvalidArgs error = errors.New("Arguments must be in the form of: SERVICE PRIVATE_PORT")
)

type ComposeAction func(c compose.Compose, cmd *cobra.Command, args []string) error

func execAction(action ComposeAction) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		composeFile, _ := cmd.Flags().GetString(COMPOSE_FILE_FLAG)
		projName, _ := cmd.Flags().GetString(PROJECT_NAME_FLAG)

		context := &compose.Context{
			ComposeFile: composeFile,
			ProjectName: projName,
		}
		compose := compose.NewCompose(context)
		err := action(compose, cmd, args)

		if err != nil {
			cli.Output(nil, err)
		}
	}
}

func logs(c compose.Compose, cmd *cobra.Command, args []string) error {
	return c.Logs(args...)
}

func build(c compose.Compose, cmd *cobra.Command, args []string) error {
	return c.Build(args...)
}

func delete(c compose.Compose, cmd *cobra.Command, args []string) error {
	return c.Delete(args...)
}

func ps(c compose.Compose, cmd *cobra.Command, args []string) error {
	q, _ := cmd.Flags().GetBool(QUIET_FLAG)
	return c.PS(q)
}

func restart(c compose.Compose, cmd *cobra.Command, args []string) error {
	return c.Restart(args...)
}

func pull(c compose.Compose, cmd *cobra.Command, args []string) error {
	return c.Pull(args...)
}

func start(c compose.Compose, cmd *cobra.Command, args []string) error {
	return c.Start(args...)
}

func stop(c compose.Compose, cmd *cobra.Command, args []string) error {
	return c.Stop(args...)
}

func kill(c compose.Compose, cmd *cobra.Command, args []string) error {
	return c.Kill(args...)
}

func port(c compose.Compose, cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return PortInvalidArgs
	}
	index, _ := cmd.Flags().GetInt(INDEX_FLAG)
	proto, _ := cmd.Flags().GetString(PROTO_FLAG)

	return c.Port(index, proto, args[0], args[1])
}

func up(c compose.Compose, cmd *cobra.Command, args []string) error {
	params, _ := cmd.Flags().GetStringSlice(PARAMS_FLAG)
	ignore, _ := cmd.Flags().GetBool(IGNORE_MISSING)

	if len(args) > 0 {
		return c.Up(cli.NameValueSliceToMap(params), !ignore, args...)
	}
	return c.Up(cli.NameValueSliceToMap(params), !ignore)
}
