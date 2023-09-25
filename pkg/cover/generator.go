package cover

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"math/rand"
	"os"
	"stimprint/internal/ff"
	"stimprint/internal/fileutil"
	"stimprint/internal/slogger"
	"stimprint/pkg/metadata"
)

type Params struct {
	Ffprobe                string
	Ffmpeg                 string
	OverwriteExistingFiles bool
	FrequencyLimit         int
	FrequencyHints         ff.FrequencyHints
	OutputImageWidth       int
	OutputImageHeight      int
	WaveColorLeft          string
	WaveColorRight         string
	WaveColorOverlap       string
	WaveColorTriphase      string
	WaveColorTriphaseAlpha float64
	WaveColorMono          string
}

func DefaultParams() Params {
	return defaultParams
}

func Generate(ctx context.Context, params Params, audioFilePath string, channels int, outputFilePath string) error {
	if !params.OverwriteExistingFiles && fileutil.Exists(outputFilePath) {
		slogger.Ctx(ctx).Debug("image already exists, overwrite not requested")
		return nil
	}

	tmpId := rand.Uint64()
	spectrumFile := fileutil.MkTmpFilePath(tmpId) + "_s.png"
	waveFile := fileutil.MkTmpFilePath(tmpId) + "_w.png"

	slogger.Ctx(ctx).Debug("Temporary files", "spectrum", spectrumFile, "wave", waveFile)

	defer func() {
		os.Remove(spectrumFile)
		os.Remove(waveFile)
	}()

	g, gCtx := errgroup.WithContext(ctx)

	//generate spectrum
	if params.OverwriteExistingFiles || !fileutil.Exists(spectrumFile) {
		g.Go(func() error {
			err := ff.GenerateSpectrum(gCtx, params.Ffmpeg, audioFilePath, spectrumFile, params.OutputImageWidth, params.OutputImageHeight, params.FrequencyLimit)
			if err != nil {
				return fmt.Errorf("generate spectrum: %v", err)
			}
			return nil
		})
	}

	//generate waves
	if params.OverwriteExistingFiles || !fileutil.Exists(waveFile) {
		g.Go(func() error {
			if channels == 0 {
				md, err := metadata.GetMetadata(ctx, params.Ffprobe, audioFilePath)
				if err != nil {
					return fmt.Errorf("get metadata: %v", err)
				}
				channels = md.Channels
			}

			return generateWaves(ctx, params, audioFilePath, channels, waveFile)
		})
	}

	err := g.Wait()
	if err != nil {
		return err
	}

	//crop, scale, overlay and draw hints
	if params.OverwriteExistingFiles || !fileutil.Exists(outputFilePath) {
		err = ff.Combine(ctx, params.Ffmpeg, waveFile, spectrumFile, outputFilePath, params.OutputImageHeight, params.FrequencyHints, params.FrequencyLimit)
		if err != nil {
			return fmt.Errorf("combine: %v", err)
		}
		slogger.Ctx(ctx).Info("Image generated")
	}
	return nil
}

func generateWaves(ctx context.Context, params Params, audioFilePath string, channels int, outputFilePath string) error {
	if channels == 2 {
		err := ff.GenerateTriWaves(ctx, params.Ffmpeg, audioFilePath, outputFilePath, params.OutputImageWidth, params.OutputImageHeight, params.WaveColorLeft, params.WaveColorRight, params.WaveColorOverlap, params.WaveColorTriphase, float32(params.WaveColorTriphaseAlpha))
		if err != nil {
			return fmt.Errorf("generate tri-waves: %v", err)
		}
	} else {
		err := ff.GenerateMonoWaves(ctx, params.Ffmpeg, audioFilePath, outputFilePath, params.OutputImageWidth, params.OutputImageHeight, params.WaveColorMono)
		if err != nil {
			return fmt.Errorf("generate mono wave: %v", err)
		}
	}
	return nil
}
