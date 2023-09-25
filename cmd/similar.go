package cmd

import (
	"context"
	"fmt"
	"github.com/o-fl0w/stimprint/internal/slogger"
	"github.com/o-fl0w/stimprint/pkg/chromaprint"
	"github.com/o-fl0w/stimprint/pkg/metadata"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"os"
	"path/filepath"
)

var similarCmd = &cobra.Command{
	Use:   "similar FILE_A FILE_B",
	Short: "Calculate similarity between two audio files using chromaprint",
	Long:  "Compares two files using their chromaprints. Returns a score from 0 to 1, higher meaning more similar.",
	Run:   similarCmdRun,
	Args:  cobra.ExactArgs(2),
}

func init() {
	hideHelp(similarCmd)
	rootCmd.AddCommand(similarCmd)
}

func similarCmdRun(_ *cobra.Command, args []string) {
	ctx := slogger.WithContext(context.TODO(), slog.Default())

	fileA := args[0]
	fileB := args[1]
	ffmpeg := filepath.Join(ffmpegRoot, "ffmpeg")
	ffprobe := filepath.Join(ffmpegRoot, "ffprobe")

	var mdA metadata.Metadata
	var mdB metadata.Metadata

	gMd, _ := errgroup.WithContext(ctx)
	gMd.Go(func() error {
		var err error
		mdA, err = metadata.GetMetadata(ctx, ffprobe, fileA)
		if err != nil {
			return fmt.Errorf("error getting metadata for file '%s': %v", fileA, err)
		}
		return nil
	})

	gMd.Go(func() error {
		var err error
		mdB, err = metadata.GetMetadata(ctx, ffprobe, fileB)
		if err != nil {
			return fmt.Errorf("error getting metadata for file '%s': %v", fileB, err)
		}
		return nil
	})

	err := gMd.Wait()
	if err != nil {
		slogger.Ctx(ctx).Error("Error getting metadata", "error", err)
		os.Exit(1)
	}

	if mdA.Duration != mdB.Duration {
		slogger.Ctx(ctx).Error("Audio must be of same duration")
		os.Exit(0)
	}

	var cpA []int32
	var cpB []int32

	gCp, _ := errgroup.WithContext(ctx)
	gCp.Go(func() error {
		bsA, err := chromaprint.Calculate(ctx, ffmpeg, fileA)
		if err != nil {
			return fmt.Errorf("error generatinc chromaprint for file '%s': %v", fileA, err)
		}
		cpA, _ = chromaprint.FromBytes(bsA)
		return nil
	})

	gCp.Go(func() error {
		bsB, err := chromaprint.Calculate(ctx, ffmpeg, fileB)
		if err != nil {
			return fmt.Errorf("error generatinc chromaprint for file '%s': %v", fileB, err)
		}
		cpB, _ = chromaprint.FromBytes(bsB)
		return nil
	})

	err = gCp.Wait()
	if err != nil {
		slogger.Ctx(ctx).Error("Error generating chromaprint", "error", err)
		os.Exit(1)
	}

	similarity, err := chromaprint.Compare(cpA, cpB)
	if err != nil {
		slogger.Ctx(ctx).Error("Error calculating similarity", "error", err)
		os.Exit(1)
	}

	fmt.Printf("%f", similarity)
}
