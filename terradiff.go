package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cockroachdb/errors"
	"github.com/google/go-cmp/cmp"
	"github.com/urfave/cli/v2"
)

type gitCloner interface {
	GitClone(ctx context.Context) error
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

	return td.compare(srcResult, dstResult), nil
}

func (op *Operators) gitCloneAndgenPlanResult(ctx context.Context) (*PlanResult, error) {
	err := op.git.GitClone(ctx)
	if err != nil {
		return nil, err
	}
	pr, err := op.tf.GenPlanResult(ctx)
	if err != nil {
		return nil, err
	}
	return pr, nil
}

type CompareResult struct {
	IsEqual bool
	Diff    string
}

func (td *Terradiff) compare(src, dst *PlanResult) *CompareResult {
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

func TerradiffAction(cCtx *cli.Context) error {
	SetLogger(cCtx)

	srcBranch := cCtx.String(srcBranchFlag)
	dstBranch := cCtx.String(dstBranchFlag)
	repoURL := cCtx.String(repoURLFlag)
	workDir := cCtx.String(workDirFlag)

	srcDir := filepath.Join(workDir, "src")
	dstDir := filepath.Join(workDir, "dst")
	srcGit := NewGit(NewCommandExecutor("."), repoURL, srcDir, srcBranch)
	dstGit := NewGit(NewCommandExecutor("."), repoURL, dstDir, dstBranch)

	srcTerraformDir := filepath.Join(srcDir, cCtx.String(terraformDirFlag))
	dstTerraformDir := filepath.Join(dstDir, cCtx.String(terraformDirFlag))
	srcTerraform := NewTerraform(NewCommandExecutor(srcTerraformDir))
	dstTerrafom := NewTerraform(NewCommandExecutor(dstTerraformDir))

	terradiff := NewTerradiff(
		Operators{tf: srcTerraform, git: srcGit},
		Operators{tf: dstTerrafom, git: dstGit},
	)

	cr, err := terradiff.Compare(cCtx.Context)
	if err != nil {
		return err
	}
	Logger().Info(fmt.Sprintf("%+v", cr))

	// TODO: resultのファイル名変更できるようにする
	err = saveToFile(workDir+"/result.json", *cr)
	if err != nil {
		return err
	}

	return nil
}

func saveToFile(filename string, cr CompareResult) error {
	file, err := os.Create(filename)
	if err != nil {
		return errors.Wrap(err, "failed to create file")
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(cr)
	if err != nil {
		return errors.Wrap(err, "failed to encode json")
	}

	return nil
}
