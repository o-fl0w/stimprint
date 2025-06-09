package cover

import (
	"context"
	"github.com/o-fl0w/stimprint/internal/ff"
	"github.com/o-fl0w/stimprint/internal/fileutil"
	"github.com/o-fl0w/stimprint/internal/slogger"
	"github.com/o-fl0w/stimprint/pkg/cover/param"
)

func Generate(ctx context.Context, params param.Cover, audioFilePath string, channels int, outputFilePath string, overwrite bool) error {
	if !overwrite && fileutil.Exists(outputFilePath) {
		slogger.Ctx(ctx).Debug("output file already exists, overwrite not requested")
		return nil
	}
	return ff.GenerateCover(ctx, params, audioFilePath, channels, outputFilePath)
}
