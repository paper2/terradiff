package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/cockroachdb/errors"
	"github.com/google/go-cmp/cmp"
	"github.com/urfave/cli/v2"
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
	cr, err := teradiff(cCtx.Context, srcGit, dstGit)
	if err != nil {
		return err
	}

	err = saveToFile(workDir+"/result.json", cr)
	if err != nil {
		return err
	}

	return nil
}

func teradiff(ctx context.Context, srcGit, dstGit gitCloner) (*CompareResult, error) {
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

type CompareResult struct {
	IsEqual bool
	Diff    string
}

func compare(src, dst *PlanResult) *CompareResult {
	var cr CompareResult
	if cmp.Equal(src, dst) {
		Logger().Info("resources are the same.")
		cr.IsEqual = true
	} else {
		Logger().Info("resources are different.")
		cr.IsEqual = false
		cr.Diff = cmp.Diff(src, dst)
	}
	return &cr
}

func saveToFile(filename string, data any) error {
	file, err := os.Create(filename)
	if err != nil {
		return errors.Wrap(err, "failed to create file")
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(data)
	if err != nil {
		return errors.Wrap(err, "failed to encode json")
	}

	return nil
}
