package cmd

import (
	"fmt"
	"github.com/lmittmann/tint"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
	"time"
)

var rootCmd = &cobra.Command{
	Use:   "stimprint",
	Short: "Audio stim file fingerprinting and visualizing tool",
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

var (
	ffmpegRoot string
)

func init() {
	rootCmd.PersistentFlags().StringVar(&ffmpegRoot, "ffmpegRoot", ``, "Path to ffmpeg bin dir (if none provided ffmpeg must be available through $PATH)")
}

func Execute() {
	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.Kitchen,
		}),
	))

	hideHelp(rootCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func hideHelp(cmd *cobra.Command) {
	cmd.Flags().Bool("help", false, "")
	_ = cmd.Flags().MarkHidden("help")
	cmd.SetHelpCommand(&cobra.Command{Hidden: true})
}
