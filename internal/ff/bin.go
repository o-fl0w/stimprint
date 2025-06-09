package ff

import (
	"context"
	"github.com/o-fl0w/stimprint/internal/slogger"
	"os/exec"
	"strings"

	"slices"
)

func ffmpegCmd(ctx context.Context, inputFilePaths []string, args []string) *exec.Cmd {
	if !slices.Contains(args, "-hide_banner") {
		args = append(args, "-hide_banner")
	}
	if !slices.Contains(args, "-loglevel") {
		args = append(args, "-loglevel", "level+error")
	}

	for _, f := range inputFilePaths {
		args = append(args, "-i", f)
	}
	slogger.Ctx(ctx).Debug("Cmd", "args", strings.Join(args, " "))
	return exec.CommandContext(ctx, "ffmpeg", args...)
}

func Run(ctx context.Context, inputFilePaths []string, args []string) error {
	cmd := ffmpegCmd(ctx, inputFilePaths, args)
	return cmd.Run()
}

func Output(ctx context.Context, inputFilePaths []string, args []string) ([]byte, error) {
	cmd := ffmpegCmd(ctx, inputFilePaths, args)
	return cmd.CombinedOutput()
}

func OutputString(ctx context.Context, inputFilePaths []string, args []string) (string, error) {
	cmd := ffmpegCmd(ctx, inputFilePaths, args)
	bs, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

func ffprobeCmd(ctx context.Context, inputFilePath string, args []string) *exec.Cmd {
	if !slices.Contains(args, "-hide_banner") {
		args = append(args, "-hide_banner")
	}
	if !slices.Contains(args, "-loglevel") {
		args = append(args, "-loglevel", "level+error")
	}

	args = append(args, "-i", inputFilePath, "-print_format", "json")

	return exec.CommandContext(ctx, "ffprobe", args...)
}

func ProbeOutput(ctx context.Context, inputFilePath string, args []string) ([]byte, error) {
	cmd := ffprobeCmd(ctx, inputFilePath, args)
	return cmd.CombinedOutput()
}
