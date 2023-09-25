package ff

import (
	"context"
	"fmt"
	"stimprint/internal/bin"
)

func GenerateSpectrum(ctx context.Context, ffmpeg string, audioFilePath string, outFilePath string, outputImageWidth int, outputImageHeight int, maxFrequency int) error {
	args := []string{
		"-v", "error",
		"-y",
		"-i", audioFilePath,
		"-filter_complex",
		fmt.Sprintf("showspectrumpic=s=%dx%d:legend=0:mode=combined:fscale=lin:start=0:stop=%d:scale=lin:drange=20:limit=0,colorkey=0x000000:0.1:0",
			outputImageWidth, outputImageHeight,
			maxFrequency),
		outFilePath}
	return bin.Path(ffmpeg).Run(ctx, args...)
}
