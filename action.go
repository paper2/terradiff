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

	srcResult, err := gitCloneAndgenPlanResult(cCtx.Context, workDir, srcDir, srcBranch, repoURL)
	if err != nil {
		return err
	}

	dstResult, err := gitCloneAndgenPlanResult(cCtx.Context, workDir, dstDir, dstBranch, repoURL)
	if err != nil {
		return err
	}

	err = compareResult(srcResult, dstResult)
	if err != nil {
		return err
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

func gitCloneAndgenPlanResult(ctx context.Context, workDir, cloneDir, branch, repoURL string) (*PlanResult, error) {
	err := gitClone(ctx, workDir, cloneDir, branch, repoURL)
	if err != nil {
		return nil, err
	}
	pr, err := generatePlanResult(ctx, cloneDir)
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
