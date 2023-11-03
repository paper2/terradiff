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

func (a *Actions) Terradiff(cCtx *cli.Context) error {
	if err := teradiff(cCtx); err != nil {
		return err
	}
	return nil
}

func teradiff(cCtx *cli.Context) error {
	branch := cCtx.String("branch")
	// TODO: 指定できるようにする。新規作成できるようにするのもありかもなあ。
	testDir := "/Users/yohei/Desktop/test-git"

	srcResult, err := genPlanRsesultWithCheckout(cCtx.Context, branch, testDir)
	if err != nil {
		return err
	}

	destResult, err := genPlanRsesultWithCheckout(cCtx.Context, "main", testDir)
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

func genPlanRsesultWithCheckout(ctx context.Context, branch string, testDir string) (*PlanResult, error) {
	git, err := NewGit(testDir)
	if err != nil {
		return nil, err
	}

	err = git.Checkout(branch)
	if err != nil {
		return nil, err
	}

	return generatePlanResult(ctx, testDir)
}
