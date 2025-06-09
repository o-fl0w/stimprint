package hash

import (
	"context"
	"fmt"
	"github.com/o-fl0w/stimprint/internal/ff"
	"regexp"
)

var r = regexp.MustCompile("=([a-fA-F0-9]+)")

func CalculateAudioHash(ctx context.Context, audioFilePath string, hashType string) (string, error) {
	out, err := ff.OutputString(ctx, []string{audioFilePath}, []string{"-f", "hash", "-hash", hashType, "-"})
	if err != nil {
		return "", err
	}

	ss := r.FindStringSubmatch(out)
	if len(ss) != 2 {
		return "", fmt.Errorf("error parsing ffmpeg output for hash: %s", out)
	}
	return ss[1], nil
}
