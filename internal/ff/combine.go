package ff

import (
	"context"
	"fmt"
	"regexp"
	"stimprint/internal/bin"
	"strconv"
	"strings"
)

const lineHeight = 2

func Combine(ctx context.Context, ffmpeg string, wavesImageFilePath string, spectrumImageFilePath string, outFilePath string, outputImageHeight int, freqHints []FrequencyHint, maxFrequency int) error {
	crop, err := cropDetect(ctx, ffmpeg, wavesImageFilePath)
	if err != nil {
		return fmt.Errorf("error detecting crop dimensions: %w", err)
	}

	fc := fmt.Sprintf("[0]crop=h=%d:x=0:y=%d,scale=h=%d:flags=neighbor[w];"+
		"[w][1]overlay=format=rgb", crop.height, crop.y, outputImageHeight)

	if len(freqHints) > 0 {
		lineArgs := make([]string, len(freqHints))
		for i, l := range freqHints {
			lineArgs[i] = fmt.Sprintf("drawbox=y=%d:h=%d:c=%s:t=fill:replace=1", l.Y(maxFrequency, outputImageHeight)-lineHeight/2, lineHeight, l.Color)
		}
		fc += "[ws];[ws]" + strings.Join(lineArgs, ",")
	}

	args := []string{
		"-v", "error",
		"-y",
		"-i", wavesImageFilePath,
		"-i", spectrumImageFilePath,
		"-filter_complex", fc,
		outFilePath}
	return bin.Path(ffmpeg).Run(ctx, args...)
}

type cropDim struct {
	height int
	y      int
}

var rCrop = regexp.MustCompile(`crop=\d+:(\d+):\d+:(\d+)`)

func cropDetect(ctx context.Context, ffmpeg string, imageFilePath string) (cropDim, error) {
	//ffmpeg -i image -filter_complex "color=black,format=rgba[b];[b][0]scale2ref[b][i];[b][i]overlay[bi];[bi]cropdetect=round=0:limit=0" -frames:v 2 -f null -
	args := []string{
		"-v", "info",
		"-i", imageFilePath,
		"-filter_complex", "color=black,format=rgba[b];[b][0]scale2ref[b][i];[b][i]overlay[bi];[bi]cropdetect=round=0:limit=0",
		"-frames:v", "2",
		"-f", "null", "-",
	}
	out, err := bin.Path(ffmpeg).CombinedOutput(ctx, args...)

	if err != nil {
		return cropDim{}, err
	}

	cd := cropDim{}
	ms := rCrop.FindAllStringSubmatch(string(out), -1)
	cd.height, _ = strconv.Atoi(ms[0][1])
	cd.y, _ = strconv.Atoi(ms[0][2])

	return cd, nil
}
