package chromaprint

import (
	"context"
	"github.com/o-fl0w/stimprint/internal/ff"
)

func Calculate(ctx context.Context, audioFilePath string) ([]byte, error) {
	bs, err := ff.Output(ctx, []string{audioFilePath}, []string{
		"-f", "chromaprint",
		"-fp_format", "raw", "-"})
	if err != nil {
		return nil, err
	}
	return bs, nil
}
