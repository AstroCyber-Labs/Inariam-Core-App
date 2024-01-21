package cli

import (
	"gitea/pcp-inariam/inariam/core/config"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration related commands",
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check the configuration",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		err := config.Check(args[0]) // Assuming you have a Check function in the config package

		if err != nil {
			color.Red("Configuration Error: %w", err)
		} else {
			color.Green("Configuration looks good!")
		}
	},
}
