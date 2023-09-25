package ff

import (
	"context"
	"fmt"
	"github.com/o-fl0w/stimprint/internal/bin"
	"strconv"
)

type rgb struct {
	r uint8
	g uint8
	b uint8
}

func hex2rgb(hex string) (rgb, error) {
	if hex[0] == '#' {
		hex = hex[1:]
	}
	var c rgb
	values, err := strconv.ParseUint(string(hex), 16, 32)

	if err != nil {
		return rgb{}, err
	}

	c = rgb{
		r: uint8(values >> 16),
		g: uint8((values >> 8) & 0xFF),
		b: uint8(values & 0xFF),
	}

	return c, nil
}

func GenerateTriWaves(ctx context.Context, ffmpeg string, audioFilePath string, outFilePath string, outputImageWidth int, outputImageHeight int,
	waveColorLeft string, waveColorRight string, waveColorOverlap string, waveColorTriphase string, waveColorTriphaseAlpha float32) error {
	rgbLeft, err := hex2rgb(waveColorLeft)
	if err != nil {
		return fmt.Errorf("error parsing color (%s): %w", waveColorLeft, err)
	}
	rgbRight, err := hex2rgb(waveColorRight)
	if err != nil {
		return fmt.Errorf("error parsing color (%s): %w", waveColorRight, err)
	}
	rgbOverlap, err := hex2rgb(waveColorOverlap)
	if err != nil {
		return fmt.Errorf("error parsing color (%s): %w", waveColorOverlap, err)
	}

	args := []string{
		"-v", "error",
		"-y",
		"-i", audioFilePath,
		"-frames:v", "1",
		"-filter_complex",
		fmt.Sprintf("[0]asplit=2[A][B];"+
			"[A]aformat=channel_layouts=mono,showwavespic=s=%dx%d:split_channels=0:colors=%s,format=rgba,lut=a='if(val,255*%f,0)'[T];"+
			"[B]showwavespic=s=%dx%d:split_channels=1:colors=%s|%s,format=rgba,split[B1][B2];"+
			"[B1]crop=w=%d:h=%d:x=0:y=0[L];"+
			"[B2]crop=w=%d:h=%d:x=0:y=%d[R];"+
			"[L][R]blend=all_mode=addition[LR];"+
			"[LR]geq=r='if( eq(r(X,Y),0x%x) * eq(g(X,Y),0x%x) * eq(b(X,Y),0x%x) + eq(r(X,Y),0x%x) * eq(g(X,Y),0x%x) * eq(b(X,Y),0x%x) + eq(alpha(X,Y),0), p(X,Y), 0x%x)'"+
			":g='if(eq(r(X,Y),0x%x) * eq(g(X,Y),0x%x) * eq(b(X,Y),0x%x) + eq(r(X,Y),0x%x) * eq(g(X,Y),0x%x) * eq(b(X,Y),0x%x) + eq(alpha(X,Y),0), p(X,Y),0x%x)'"+
			":b='if(eq(r(X,Y),0x%x) * eq(g(X,Y),0x%x) * eq(b(X,Y),0x%x) + eq(r(X,Y),0x%x) * eq(g(X,Y),0x%x) * eq(b(X,Y),0x%x) + eq(alpha(X,Y),0), p(X,Y), 0x%x)'"+
			":a='if(p(X,Y),255,0)'[LR];"+
			"[LR][T]overlay=format=rgb",
			outputImageWidth, outputImageHeight, waveColorTriphase, waveColorTriphaseAlpha,
			outputImageWidth, outputImageHeight*2, waveColorLeft, waveColorRight,
			outputImageWidth, outputImageHeight,
			outputImageWidth, outputImageHeight, outputImageHeight,
			rgbLeft.r, rgbLeft.g, rgbLeft.b, rgbRight.r, rgbRight.g, rgbRight.b, rgbOverlap.r,
			rgbLeft.r, rgbLeft.g, rgbLeft.b, rgbRight.r, rgbRight.g, rgbRight.b, rgbOverlap.g,
			rgbLeft.r, rgbLeft.g, rgbLeft.b, rgbRight.r, rgbRight.g, rgbRight.b, rgbOverlap.b,
		),
		outFilePath}
	return bin.Path(ffmpeg).Run(ctx, args...)
}

func GenerateMonoWaves(ctx context.Context, ffmpeg string, audioFilePath string, outFilePath string, outputImageWidth int, outputImageHeight int, waveColorMono string) error {
	args := []string{
		"-v", "quiet",
		"-y",
		"-i", audioFilePath,
		"-filter_complex",
		fmt.Sprintf("showwavespic=s=%dx%d:split_channels=0:colors=%s,format=rgba",
			outputImageWidth, outputImageHeight, waveColorMono),
		outFilePath}
	return bin.Path(ffmpeg).Run(ctx, args...)
}
