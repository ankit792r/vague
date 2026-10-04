package command

import (
	"context"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "vague [files...]",
	Short: "Vim Emacs inspired text editor, WITH HATE OF BOTH WORLD",
	Args:  cobra.ArbitraryArgs,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Help()
		// if err := process.EnsureServer(cmd.Context()); err != nil {
		// 	return err
		// }
		//
		// return webview.Run(cmd.Context(), args...)
		return nil
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
