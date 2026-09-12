package command

import (
	client "vagues/backbone/clients"
	"vagues/backbone/process"

	"github.com/spf13/cobra"
)

var clientConnectCmd = &cobra.Command{
	Use:   "client-connect",
	Short: "Connect vague deamon client",
	RunE: func(cmd *cobra.Command, args []string) error {
		conn, err := client.ClientConnect()

		ctx := cmd.Context()

		if err != nil {
			return err
		}
		defer conn.Close()

		conn.Execute(ctx, process.ExecuteParams{
			Name:  "test",
			Args:  []string{"some", "args"},
			Bang:  false,
			Count: 0,
		}, nil)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(clientConnectCmd)
}
