package acrc32

import (
	"context"
	"fmt"
	"regexp"
	"stimprint/internal/bin"
)

var r = regexp.MustCompile("CRC32=(\\w{8})")

func Calculate(ctx context.Context, ffmpeg string, audioFilePath string) (string, error) {
	args := []string{
		"-v", "info",
		"-y",
		"-i", audioFilePath,
		"-f", "hash",
		"-hash", "CRC32",
		"-",
	}

	out, err := bin.Path(ffmpeg).Output(ctx, args...)
	if err != nil {
		return "", err
	}

	ss := r.FindStringSubmatch(string(out))
	if len(ss) != 2 {
		return "", fmt.Errorf("error parsing ffmpeg output for CRC32")
	}
	return ss[1], nil
}
