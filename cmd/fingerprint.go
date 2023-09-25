package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"log/slog"
	"path/filepath"
	"stimprint/internal/slogger"
	"stimprint/pkg/acrc32"
	"stimprint/pkg/oshash"
	"sync"
)

var (
	isRequestedCrc32  bool
	isRequestedOsHash bool
)

var fingerprintCmd = &cobra.Command{
	Use:     "fingerprint FILE",
	Aliases: []string{"fp"},
	Short:   "Calculate audio file fingerprints",
	Run:     fingerprintCmdRun,
	Args:    cobra.ExactArgs(1),
}

func init() {
	fingerprintCmd.Flags().BoolVar(&isRequestedCrc32, "acrc32", true, "Calculate CRC32 of audio signal")
	fingerprintCmd.Flags().BoolVar(&isRequestedOsHash, "oshash", true, "Calculate oshash of file")

	hideHelp(fingerprintCmd)
	rootCmd.AddCommand(fingerprintCmd)
}

func fingerprintCmdRun(_ *cobra.Command, args []string) {
	ctx := slogger.WithContext(context.TODO(), slog.Default())

	file := args[0]
	ffmpeg := filepath.Join(ffmpegRoot, "ffmpeg")

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		aCrc32, err := acrc32.Calculate(ctx, ffmpeg, file)
		if err != nil {
			slogger.Ctx(ctx).Error("Error generating audio CRC32", "error", err)
			return
		}
		slogger.Ctx(ctx).Info("Generated audio CRC32", "aCrc32", aCrc32)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		hash, err := oshash.FromFilePath(file)
		if err != nil {
			slogger.Ctx(ctx).Error("Error generating oshash", "error", err)
			return
		}
		slogger.Ctx(ctx).Info("Generated oshash", "oshash", hash)
	}()

	wg.Wait()
}
