package cmd

import (
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   "chart-viewer",
		Short: "http server to interact with helm",
		Run: func(c *cobra.Command, args []string) {
			c.HelpFunc()(c, args)
		},
	}

	command.AddCommand(
		NewServeCommand(),
		NewSeedCommand(),
	)

	return command
}
