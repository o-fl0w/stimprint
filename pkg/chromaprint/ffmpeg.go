package chromaprint

import (
	"context"
	"github.com/o-fl0w/stimprint/internal/bin"
)

func Calculate(ctx context.Context, ffmpeg string, audioFilePath string) ([]byte, error) {
	args := []string{
		"-v", "error",
		"-y",
		"-i", audioFilePath,
		"-f", "chromaprint",
		"-fp_format", "raw", "-"}
	bs, err := bin.Path(ffmpeg).Output(ctx, args...)
	if err != nil {
		return nil, err
	}
	return bs, nil
}
