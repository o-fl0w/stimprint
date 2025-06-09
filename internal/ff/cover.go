package ff

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
)

var overlayWaveAndSpectrumFilter = "[wave][spectrum]overlay=format=rgb"

type Params struct {
	FrequencyLimit         int
	FrequencyHints         FrequencyHints
	OutputImageWidth       int
	OutputImageHeight      int
	WaveColorLeft          string
	WaveColorRight         string
	WaveColorOverlap       string
	WaveColorTriphase      string
	WaveColorTriphaseAlpha float64
	WaveColorMono          string
}

func GenerateCover(ctx context.Context, params Params, audioFilePath string, numChannels int, outputFilePath string) error {
	switch numChannels {
	case 1:
		return generateMonoCover(ctx, params, audioFilePath, outputFilePath)
	case 2:
		return generateTriCover(ctx, params, audioFilePath, outputFilePath)
	default:
		return fmt.Errorf("illegal argument numChannels=%d: more than 2 channels not supported", numChannels)
	}
}

func generateTriCover(ctx context.Context, params Params, inputFilePath string, outputFilePath string) error {
	colorLeft := hex2rgb(params.WaveColorLeft)
	colorRight := hex2rgb(params.WaveColorLeft)
	colorOverlap := hex2rgb(params.WaveColorOverlap)

	filterParts := []string{
		spectrumFilter(params),
		"[0:a]asplit=2[A][B];",
		fmt.Sprintf(
			"[A]aformat=channel_layouts=mono,showwavespic=s=%dx%d:split_channels=0:colors=%s,format=rgba,lut=a='if(val,255*%f,0)'[T];",
			params.OutputImageWidth,
			params.OutputImageHeight,
			params.WaveColorTriphase,
			params.WaveColorTriphaseAlpha),
		fmt.Sprintf("[B]showwavespic=s=%dx%d:split_channels=1:colors=%s|%s,format=rgba,split[B1][B2];",
			params.OutputImageWidth,
			params.OutputImageHeight*2,
			params.WaveColorLeft,
			params.WaveColorRight),
		fmt.Sprintf("[B1]crop=w=%d:h=%d:x=0:y=0[L];", params.OutputImageWidth, params.OutputImageHeight),
		fmt.Sprintf("[B2]crop=w=%d:h=%d:x=0:y=%d[R];",
			params.OutputImageWidth,
			params.OutputImageHeight,
			params.OutputImageHeight),
		"[L][R]blend=all_mode=addition[LR];",
		fmt.Sprintf(
			"[LR]geq=r='if( eq(r(X,Y),0x%x) * eq(g(X,Y),0x%x) * eq(b(X,Y),0x%x) + eq(r(X,Y),0x%x) * eq(g(X,Y),0x%x) * eq(b(X,Y),0x%x) + eq(alpha(X,Y),0), p(X,Y), 0x%x)'",
			colorLeft.r,
			colorLeft.g,
			colorLeft.b,
			colorRight.r,
			colorRight.g,
			colorRight.b,
			colorOverlap.r),
		fmt.Sprintf(
			":g='if(eq(r(X,Y),0x%x) * eq(g(X,Y),0x%x) * eq(b(X,Y),0x%x) + eq(r(X,Y),0x%x) * eq(g(X,Y),0x%x) * eq(b(X,Y),0x%x) + eq(alpha(X,Y),0), p(X,Y),0x%x)'",
			colorLeft.r,
			colorLeft.g,
			colorLeft.b,
			colorRight.r,
			colorRight.g,
			colorRight.b,
			colorOverlap.g),
		fmt.Sprintf(
			":b='if(eq(r(X,Y),0x%x) * eq(g(X,Y),0x%x) * eq(b(X,Y),0x%x) + eq(r(X,Y),0x%x) * eq(g(X,Y),0x%x) * eq(b(X,Y),0x%x) + eq(alpha(X,Y),0), p(X,Y), 0x%x)'",
			colorLeft.r,
			colorLeft.g,
			colorLeft.b,
			colorRight.r,
			colorRight.g,
			colorRight.b,
			colorOverlap.b),
		":a='if(p(X,Y),255,0)'[LR];",
		"[LR][T]overlay=format=rgb[wave];",
		overlayWaveAndSpectrumFilter,
		frequencyLineFilters(params),
	}
	filter := strings.Join(filterParts, "")
	return ffExec(ctx, inputFilePath, filter, outputFilePath)
}

func generateMonoCover(ctx context.Context, params Params, inputFilePath string, outputFilePath string) error {
	filterParts := []string{
		spectrumFilter(params),
		fmt.Sprintf("[0:a]showwavespic=s=%dx%d:split_channels=0:colors=%s,format=rgba[wave];",
			params.OutputImageWidth,
			params.OutputImageHeight,
			params.WaveColorMono),
		overlayWaveAndSpectrumFilter,
		frequencyLineFilters(params),
	}
	filter := strings.Join(filterParts, "")
	return ffExec(ctx, inputFilePath, filter, outputFilePath)
}

func frequencyLineFilters(params Params) string {
	if len(params.FrequencyHints) == 0 {
		return ""
	}
	var lines []string
	for _, fh := range params.FrequencyHints {
		lines = append(lines, frequencyLine(params, fh))
	}
	return "[ws];[ws]" + strings.Join(lines, ",")
}

func frequencyLine(params Params, fh FrequencyHint) string {
	y := float64(params.OutputImageHeight) - float64(fh.Hz)/float64(params.FrequencyLimit)*float64(params.OutputImageHeight) - 1.0
	return fmt.Sprintf("drawbox=y=%f:h=2:c=%s:t=fill:replace=1", y, fh.Color)
}

func spectrumFilter(params Params) string {
	return fmt.Sprintf(
		"[0:a]showspectrumpic=s=%dx%d:legend=0:mode=combined:fscale=lin:start=0:stop=%d:scale=lin:drange=20:limit=0,colorkey=0x000000:0.1:0[spectrum];",
		params.OutputImageWidth,
		params.OutputImageHeight,
		params.FrequencyLimit)
}

func ffExec(ctx context.Context, inputFilePath string, filter string, outputFilePath string) error {
	var args = []string{
		"-y", //overwrite
		"-metadata", "Title=" + filepath.Base(inputFilePath),
		"-filter_complex", filter,
		outputFilePath}
	return Run(ctx, []string{inputFilePath}, args)
}
