package main

import (
	"context"
	"fmt"
	"os"

	"github.com/cockroachdb/errors"
	"github.com/urfave/cli/v2"
	"github.com/wI2L/jsondiff"
)

type Actions struct {
}

func NewActions() *Actions {
	return &Actions{}
}

func (a *Actions) Terradiff(cCtx *cli.Context) error {
	if err := teradiff(cCtx); err != nil {
		return err
	}
	return nil
}

func teradiff(cCtx *cli.Context) error {
	srcBranch := cCtx.String("branch")
	// TODO: 引数から受け取る
	dstBranch := "main"
	// TODO: 指定できるようにする。新規作成できるようにするのもありかもなあ。
	workDir := "/Users/yohei/Desktop/teradiff-work"
	srcDir := workDir + "/src"
	dstDir := workDir + "/dst"
	repoURL := "https://github.com/paper2/test-terradiff"

	err := gitClone(cCtx.Context, workDir, srcDir, srcBranch, repoURL)
	if err != nil {
		return err
	}
	srcResult, err := generatePlanResult(cCtx.Context, srcDir)
	if err != nil {
		return err
	}

	err = gitClone(cCtx.Context, workDir, dstDir, dstBranch, repoURL)
	if err != nil {
		return err
	}
	destResult, err := generatePlanResult(cCtx.Context, dstDir)
	if err != nil {
		return err
	}

	patch, err := jsondiff.Compare(srcResult, destResult)
	if err != nil {
		return errors.Wrap(err, "json diff")
	}
	for _, op := range patch {
		Logger().Debug(fmt.Sprintf("op: %+v", op))
	}

	return nil
}

func gitClone(ctx context.Context, workDir, cloneDir, branch, repoURL string) error {
	if _, err := os.Stat(cloneDir); !os.IsNotExist(err) {
		Logger().Info("skip clone repository.", "repository", repoURL, "branch", branch)
		return nil
	}

	err := NewCommandExecutor(workDir).RunContext(ctx, "git", "clone", "-b", branch, repoURL, cloneDir)
	if err != nil {
		return err
	}
	return nil
}
