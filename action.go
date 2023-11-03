package main

import (
	"context"
	"fmt"

	"github.com/cockroachdb/errors"
	"github.com/urfave/cli/v2"
	"github.com/wI2L/jsondiff"
)

type Actions struct {
}

func NewActions() *Actions {
	return &Actions{}
}

type gitCloner interface {
	gitClone(ctx context.Context) error
	getCloneDir() string
}

func (a *Actions) Terradiff(cCtx *cli.Context) error {
	srcBranch := cCtx.String("branch")
	// TODO: 引数から受け取る
	dstBranch := "main"
	// TODO: 指定できるようにする。新規作成できるようにするのもありかもなあ。
	repoURL := "https://github.com/paper2/test-terradiff"
	workDir := "/Users/yohei/Desktop/teradiff-work"
	srcDir := workDir + "/src"
	dstDir := workDir + "/dst"

	srcGit := NewGit(repoURL, srcDir, srcBranch)
	dstGit := NewGit(repoURL, dstDir, dstBranch)
	if err := teradiff(cCtx.Context, srcGit, dstGit); err != nil {
		return err
	}
	return nil
}

func teradiff(ctx context.Context, srcGit, dstGit gitCloner) error {
	srcResult, err := gitCloneAndgenPlanResult(ctx, srcGit)
	if err != nil {
		return err
	}

	dstResult, err := gitCloneAndgenPlanResult(ctx, dstGit)
	if err != nil {
		return err
	}

	err = compareResult(srcResult, dstResult)
	if err != nil {
		return err
	}

	return nil
}

func gitCloneAndgenPlanResult(ctx context.Context, git gitCloner) (*PlanResult, error) {
	err := git.gitClone(ctx)
	if err != nil {
		return nil, err
	}
	pr, err := generatePlanResult(ctx, git.getCloneDir())
	if err != nil {
		return nil, err
	}
	return pr, nil
}

func compareResult(src, dst *PlanResult) error {
	patch, err := jsondiff.Compare(src, dst)
	if err != nil {
		return errors.Wrap(err, "json diff")
	}
	for _, op := range patch {
		Logger().Debug(fmt.Sprintf("op: %+v", op))
	}
	return nil
}
