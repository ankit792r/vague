package command

import (
	"fmt"
	"os"
	"vagues/backbone/process"
	server "vagues/backbone/servers"

	"github.com/spf13/cobra"
)

var serverStartCmd = &cobra.Command{
	Use:   "server-start",
	Short: "Start vague deamon server",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		listener, err := process.Listen()
		if err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "vague server listening on %s\n", listener.Path())

		if err := server.NewServer().Serve(ctx, listener); err != nil {
			return fmt.Errorf("server: %w", err)
		}

		return nil

	},
}

func init() {
	rootCmd.AddCommand(serverStartCmd)
}
