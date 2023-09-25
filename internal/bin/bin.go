package bin

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
)

type Path string

func (p Path) Run(ctx context.Context, args ...string) error {
	_, err := run(ctx, false, string(p), args...)
	return err
}

func (p Path) CombinedOutput(ctx context.Context, args ...string) ([]byte, error) {
	out, err := run(ctx, true, string(p), args...)
	return out, err
}

func (p Path) Output(ctx context.Context, args ...string) ([]byte, error) {
	out, err := run(ctx, false, string(p), args...)
	return out, err
}

func run(ctx context.Context, combinedOutput bool, name string, args ...string) ([]byte, error) {
	var output bytes.Buffer
	var stdErr *bytes.Buffer
	if combinedOutput {
		stdErr = &output
	} else {
		var b bytes.Buffer
		stdErr = &b
	}

	//slog.Debug("exec: " + name + " " + strings.Join(args, " "))

	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdout = &output
	cmd.Stderr = stdErr

	err := cmd.Run()

	if err != nil {
		return nil, fmt.Errorf("exec %s: %w: %s", name, err, stdErr.String())
	}

	return output.Bytes(), nil
}
