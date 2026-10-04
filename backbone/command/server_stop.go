package command

import "github.com/spf13/cobra"

var serverStopCmd = &cobra.Command{
	Use:   "server-stop",
	Short: "Stop vague daemon server",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	rootCmd.AddCommand(serverStopCmd)
}
