package main

import (
	"context"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

type gitCloner interface {
	// TODO: interafaceに入れるのは公開関数にする
	gitClone(ctx context.Context) error
	getCloneDir() string
}

type Terraformer interface {
	GenPlanResult(ctx context.Context) (*PlanResult, error)
}

type Operators struct {
	tf  Terraformer
	git gitCloner
}

type Terradiff struct {
	src Operators
	dst Operators
}

func NewTerradiff(src, dst Operators) *Terradiff {
	return &Terradiff{src: src, dst: dst}
}

func (td *Terradiff) Compare(ctx context.Context) (*CompareResult, error) {
	srcResult, err := td.src.gitCloneAndgenPlanResult(ctx)
	if err != nil {
		return nil, err
	}

	dstResult, err := td.dst.gitCloneAndgenPlanResult(ctx)
	if err != nil {
		return nil, err
	}

	return compare(srcResult, dstResult), nil
}

func (op *Operators) gitCloneAndgenPlanResult(ctx context.Context) (*PlanResult, error) {
	err := op.git.gitClone(ctx)
	if err != nil {
		return nil, err
	}
	pr, err := op.tf.GenPlanResult(ctx)
	if err != nil {
		return nil, err
	}
	return pr, nil
}

func TerradiffAction(cCtx *cli.Context) error {
	SetLogger(cCtx)

	srcBranch := cCtx.String(srcBranchFlag)
	dstBranch := cCtx.String(dstBranchFlag)
	repoURL := cCtx.String(repoURLFlag)
	workDir := cCtx.String(workDirFlag)

	srcDir := filepath.Join(workDir, "src")
	dstDir := filepath.Join(workDir, "dst")
	srcGit := NewGit(repoURL, srcDir, srcBranch)
	dstGit := NewGit(repoURL, dstDir, dstBranch)

	// TODO: Terraformの実行パスを指定出来るようにする
	srcTerraform := NewTerraform(NewCommandExecutor(srcDir))
	dstTerrafom := NewTerraform(NewCommandExecutor(dstDir))

	terradiff := NewTerradiff(
		Operators{tf: srcTerraform, git: srcGit},
		Operators{tf: dstTerrafom, git: dstGit},
	)

	cr, err := terradiff.Compare(cCtx.Context)
	if err != nil {
		return err
	}

	err = saveToFile(workDir+"/result.json", *cr)
	if err != nil {
		return err
	}

	return nil
}
