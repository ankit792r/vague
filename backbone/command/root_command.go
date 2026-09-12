package command

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "vague",
	Short: "A Vim-inspired text editor",

	// Usage text is only helpful for malformed invocations, not for runtime
	// failures. Errors are printed once, by main.
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Run the server and client")
		return nil
	},
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
		// DisableNoDescFlag: true,
		// DisableDescriptions: true,
		// HiddenDefaultCmd: true,
	} ,
}

func Execute(ctx context.Context) error {
	return rootCmd.ExecuteContext(ctx)
}
