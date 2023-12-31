package cmd

import (
	"context"
	"github.com/o-fl0w/stimprint/internal/slogger"
	"github.com/o-fl0w/stimprint/pkg/cover"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

var coverCmd = &cobra.Command{
	Use:   "cover INPUT_FILE OUTPUT_FILE",
	Short: "Generate a cover image.",
	Long:  "Processes an audio file and generates an image visualizing the audio signals stim-properties. Properties represented are spectrum of peak frequencies with their amplitudes and distribution as left, right and triphase output. Image can additionally contain hints of safe, preferred or other user defined frequency ranges to help determine suitability of an audio file as input to stim devices.",
	Run:   coverCmdRun,
	Args:  cobra.ExactArgs(2),
}

var (
	coverParams = cover.DefaultParams()
)

func init() {

	coverCmd.Flags().BoolVar(&coverParams.OverwriteExistingFiles, "overwrite", coverParams.OverwriteExistingFiles, "Overwrite existing files")

	coverCmd.Flags().IntVar(&coverParams.OutputImageWidth, "width", coverParams.OutputImageWidth, "Width of output image")
	coverCmd.Flags().IntVar(&coverParams.OutputImageHeight, "height", coverParams.OutputImageHeight, "Height of output image")

	coverCmd.Flags().StringVar(&coverParams.WaveColorLeft, "colorLeft", coverParams.WaveColorLeft, "Color used for the waveform representing the left output")
	coverCmd.Flags().StringVar(&coverParams.WaveColorRight, "colorRight", coverParams.WaveColorRight, "Color used for the waveform representing the right output")
	coverCmd.Flags().StringVar(&coverParams.WaveColorOverlap, "colorOverlap", coverParams.WaveColorOverlap, "Color used for the overlap of left and right output")
	coverCmd.Flags().StringVar(&coverParams.WaveColorTriphase, "colorTriphase", coverParams.WaveColorTriphase, "Color used for the waveform representing the triphase output")
	coverCmd.Flags().Float64Var(&coverParams.WaveColorTriphaseAlpha, "colorTriphaseAlpha", coverParams.WaveColorTriphaseAlpha, "Opacity of the waveform representing the triphase output")
	coverCmd.Flags().StringVar(&coverParams.WaveColorMono, "colorMono", coverParams.WaveColorMono, "Color used for the waveform representing mono output")

	coverCmd.Flags().IntVar(&coverParams.FrequencyLimit, "freqLimit", coverParams.FrequencyLimit, "Frequency ceiling of spectrum")

	//var freqHints string
	coverCmd.Flags().Var(&coverParams.FrequencyHints, "freqHints", "Draw lines at given frequencies, hz:color;hz:color...")

	_ = coverCmd.MarkFlagRequired("out")
	hideHelp(coverCmd)
	rootCmd.AddCommand(coverCmd)
}

func coverCmdRun(_ *cobra.Command, args []string) {

	ctx := slogger.WithContext(context.TODO(), slog.Default())

	inputFilePath := args[0]
	outputFilePath := args[1]

	//err := g.FrequencyHints.Set(freqHints)
	//if err != nil {
	//	flag.Usage()
	//	log.Fatalf("error parsing frequency hints: %v", err)
	//}

	coverParams.Ffprobe = filepath.Join(ffmpegRoot, "ffprobe")
	coverParams.Ffmpeg = filepath.Join(ffmpegRoot, "ffmpeg")

	//outputFilename := outputFilePath
	//if outputFilename == "" {
	//	outputFilename = fileutil.MkRandomTmpFilePath()
	//}

	slogger.Ctx(ctx).Info("Generating...", "input", inputFilePath, "output", outputFilePath)

	start := time.Now()
	err := cover.Generate(ctx, coverParams, inputFilePath, 0, outputFilePath)

	if err != nil {
		slogger.Ctx(ctx).Error("Error generating", "error", err)
		os.Exit(1)
	}
	duration := time.Since(start)

	slogger.Ctx(ctx).Info("Done", "duration", duration)
}
