package main

import (
	"context"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

type gitCloner interface {
	gitClone(ctx context.Context) error
	getCloneDir() string
}

func terradiff(cCtx *cli.Context) error {
	SetLogger(cCtx)

	srcBranch := cCtx.String(srcBranchFlag)
	dstBranch := cCtx.String(dstBranchFlag)
	repoURL := cCtx.String(repoURLFlag)
	workDir := cCtx.String(workDirFlag)

	srcDir := filepath.Join(workDir, "src")
	dstDir := filepath.Join(workDir, "dst")
	srcGit := NewGit(repoURL, srcDir, srcBranch)
	dstGit := NewGit(repoURL, dstDir, dstBranch)

	cr, err := terradiffCompare(cCtx.Context, srcGit, dstGit)
	if err != nil {
		return err
	}

	err = saveToFile(workDir+"/result.json", *cr)
	if err != nil {
		return err
	}

	return nil
}

func terradiffCompare(ctx context.Context, srcGit, dstGit gitCloner) (*CompareResult, error) {
	srcResult, err := gitCloneAndgenPlanResult(ctx, srcGit)
	if err != nil {
		return nil, err
	}

	dstResult, err := gitCloneAndgenPlanResult(ctx, dstGit)
	if err != nil {
		return nil, err
	}

	return compare(srcResult, dstResult), nil
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
