package command

import (
	"context"
	"vague/backbone/process"
	"vague/bonefire/frame"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "vague [files...]",
	Short: "A Vim-inspired text editor",
	Args:  cobra.ArbitraryArgs,

	// Usage text is only helpful for malformed invocations, not for runtime
	// failures. Errors are printed once, by main.
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := process.EnsureServer(cmd.Context()); err != nil {
			return err
		}

		return frame.NewFrame(cmd.Context(), args...)
	},
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
		// DisableNoDescFlag: true,
		// DisableDescriptions: true,
		// HiddenDefaultCmd: true,
	},
}

func Execute(ctx context.Context) error {
	return rootCmd.ExecuteContext(ctx)
}
