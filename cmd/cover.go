package cmd

import (
	"context"
	"fmt"
	"github.com/o-fl0w/stimprint/internal/slogger"
	"github.com/o-fl0w/stimprint/pkg/cover"
	"github.com/o-fl0w/stimprint/pkg/cover/param"
	"github.com/o-fl0w/stimprint/pkg/metadata"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
	"strconv"
	"strings"
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
	coverParams            = param.DefaultCoverParams()
	overwriteExistingFiles = false
)

func init() {

	coverCmd.Flags().BoolVar(&overwriteExistingFiles, "overwrite", overwriteExistingFiles, "Overwrite existing files")

	coverCmd.Flags().IntVar(&coverParams.OutputImageWidth, "width", coverParams.OutputImageWidth, "Width of output image")
	coverCmd.Flags().IntVar(&coverParams.OutputImageHeight, "height", coverParams.OutputImageHeight, "Height of output image")

	coverCmd.Flags().StringVar(&coverParams.WaveColorLeft, "colorLeft", coverParams.WaveColorLeft, "Color used for the waveform representing the left output")
	coverCmd.Flags().StringVar(&coverParams.WaveColorRight, "colorRight", coverParams.WaveColorRight, "Color used for the waveform representing the right output")
	coverCmd.Flags().StringVar(&coverParams.WaveColorOverlap, "colorOverlap", coverParams.WaveColorOverlap, "Color used for the overlap of left and right output")
	coverCmd.Flags().StringVar(&coverParams.WaveColorTriphase, "colorTriphase", coverParams.WaveColorTriphase, "Color used for the waveform representing the triphase output")
	coverCmd.Flags().Float64Var(&coverParams.WaveColorTriphaseAlpha, "colorTriphaseAlpha", coverParams.WaveColorTriphaseAlpha, "Opacity of the waveform representing the triphase output")
	coverCmd.Flags().StringVar(&coverParams.WaveColorMono, "colorMono", coverParams.WaveColorMono, "Color used for the waveform representing mono output")

	coverCmd.Flags().IntVar(&coverParams.FrequencyLimit, "freqLimit", coverParams.FrequencyLimit, "Frequency ceiling of spectrum")

	fhv := frequencyHintsValue{dest: &coverParams}
	coverCmd.Flags().Var(&fhv, "freqHints", "Draw lines at given frequencies, hz:color;hz:color...")

	_ = coverCmd.MarkFlagRequired("out")
	hideHelp(coverCmd)
	rootCmd.AddCommand(coverCmd)
}

func coverCmdRun(_ *cobra.Command, args []string) {

	ctx := slogger.WithContext(context.TODO(), slog.Default())

	inputFilePath := args[0]
	outputFilePath := args[1]

	md, err := metadata.GetMetadata(ctx, inputFilePath)
	if err != nil {
		slogger.Ctx(ctx).Error("Failed to get audio file metadata", "error", err)
		return
	}

	slogger.Ctx(ctx).Info("Generating...", "input", inputFilePath, "channels", md.Channels, "output", outputFilePath)

	start := time.Now()
	err = cover.Generate(ctx, coverParams, inputFilePath, md.Channels, outputFilePath, overwriteExistingFiles)

	if err != nil {
		slogger.Ctx(ctx).Error("Error generating", "error", err)
		os.Exit(1)
	}
	duration := time.Since(start)

	slogger.Ctx(ctx).Info("Done", "duration", duration)
}

type frequencyHintsValue struct {
	dest *param.Cover
}

func (f *frequencyHintsValue) String() string {
	kvs := make([]string, len(f.dest.FrequencyHints))
	for i, fh := range f.dest.FrequencyHints {
		kvs[i] = fmt.Sprintf("%d:%s", fh.Hz, fh.Color)
	}
	return strings.Join(kvs, ";")
}

func (f *frequencyHintsValue) Set(s string) error {
	kvs := strings.Split(s, ";")
	f.dest.FrequencyHints = make([]param.FrequencyHint, len(kvs))
	for i, kv := range kvs {
		ss := strings.Split(kv, ":")
		hz, err := strconv.Atoi(ss[0])
		if err != nil {
			return err
		}
		f.dest.FrequencyHints[i] = param.FrequencyHint{
			Hz:    hz,
			Color: ss[1],
		}
	}
	return nil
}

func (f *frequencyHintsValue) Type() string {
	return "string"
}
