package main

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"

	"github.com/cockroachdb/errors"
)

type CommandExecutor struct {
	dir string
}

func NewCommandExecutor(dir string) *CommandExecutor {
	return &CommandExecutor{dir: dir}
}

func (ce *CommandExecutor) RunContext(ctx context.Context, name string, args ...string) error {
	_, err := ce.RunContextAndCaptureOutput(ctx, name, args...)
	if err != nil {
		return err
	}
	return nil
}

func (ce *CommandExecutor) RunContextAndCaptureOutput(ctx context.Context, name string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = ce.dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	Logger().Debug(fmt.Sprintf("exec: %s", cmd.String()))
	Logger().Debug(fmt.Sprintf("stdout: %s", stdout.String()))
	if err != nil {
		return "", errors.Wrap(err, "command execution stdout: "+stderr.String())
	}
	return stdout.String(), nil
}
