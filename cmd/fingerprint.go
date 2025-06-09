package cmd

import (
	"context"
	"github.com/o-fl0w/stimprint/internal/slogger"
	"github.com/o-fl0w/stimprint/pkg/hash"
	"github.com/spf13/cobra"
	"log/slog"
	"sync"
	"time"
)

var (
	isRequestedCrc32   bool
	isRequestedMurmur3 bool
	isRequestedOsHash  bool
)

var fingerprintCmd = &cobra.Command{
	Use:     "fingerprint FILE",
	Aliases: []string{"fp"},
	Short:   "Calculate audio file fingerprints",
	Run:     fingerprintCmdRun,
	Args:    cobra.ExactArgs(1),
}

func init() {
	fingerprintCmd.Flags().BoolVar(&isRequestedCrc32, "crc32", true, "Calculate CRC32 hash of audio signal")
	fingerprintCmd.Flags().BoolVar(&isRequestedMurmur3, "murmur3", true, "Calculate murmur3 hash of audio signal")
	fingerprintCmd.Flags().BoolVar(&isRequestedOsHash, "oshash", true, "Calculate oshash of file")

	hideHelp(fingerprintCmd)
	rootCmd.AddCommand(fingerprintCmd)
}

func fingerprintCmdRun(_ *cobra.Command, args []string) {
	ctx := slogger.WithContext(context.TODO(), slog.Default())

	file := args[0]

	wg := sync.WaitGroup{}

	if isRequestedCrc32 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			start := time.Now()
			h, err := hash.CalculateAudioHash(ctx, file, "CRC32")
			duration := time.Since(start)
			if err != nil {
				slogger.Ctx(ctx).Error("Error calculating audio CRC32 hash", "error", err)
				return
			}
			slogger.Ctx(ctx).Info("Calculated audio CRC32", "CRC32", h, "duration", duration.Truncate(time.Millisecond))
		}()
	}

	if isRequestedOsHash {
		wg.Add(1)
		go func() {
			defer wg.Done()
			start := time.Now()
			h, err := hash.OsHashFromFilePath(file)
			duration := time.Since(start)
			if err != nil {
				slogger.Ctx(ctx).Error("Error calculating oshash", "error", err)
				return
			}
			slogger.Ctx(ctx).Info("Calculated file oshash", "oshash", h, "duration", duration.Truncate(time.Millisecond))
		}()
	}

	if isRequestedMurmur3 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			start := time.Now()
			h, err := hash.CalculateAudioHash(ctx, file, "murmur3")
			duration := time.Since(start)
			if err != nil {
				slogger.Ctx(ctx).Error("Error calculating audio murmur3", "error", err)
				return
			}
			slogger.Ctx(ctx).Info("Calculated audio murmur3 hash", "murmur3", h, "duration", duration.Truncate(time.Millisecond))
		}()
	}

	wg.Wait()
}
