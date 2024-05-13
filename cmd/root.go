package cmd

import "github.com/spf13/cobra"

func New() *cobra.Command {
	root := &cobra.Command{
		Use:   "main.go <subcommand>",
		Short: "Root Daemon",
	}

	addBotCommand(root)
	return root
}
