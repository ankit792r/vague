package command

import (
	"vague/backbone/client"

	"github.com/spf13/cobra"
)

var openUi = &cobra.Command{
	Use:   "open-ui",
	Short: "open vague ui",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		return client.OpenUI(ctx)
	},
}

func init() {
	rootCmd.AddCommand(openUi)
}
