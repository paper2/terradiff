package main

import (
	"os/exec"

	"github.com/cockroachdb/errors"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/urfave/cli/v2"
	"golang.org/x/exp/slog"
)

func teradiff(cCtx *cli.Context) error {
	name := cCtx.String("branch")
	testDir := "/Users/yohei/Desktop/test-git"
	r, err := git.PlainOpen(testDir)
	if err != nil {
		return errors.Wrap(err, "open git repo")
	}
	w, err := r.Worktree()
	if err != nil {
		return errors.Wrap(err, "get worktree")
	}
	// TODO: checkout するとworking directoryが消えるので消えないオプションを使った方がよさそう
	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(name),
	})
	if err != nil {
		return errors.Wrap(err, "checkout")
	}
	slog.Info("checkout: " + name)

	generateJson(testDir)

	return nil
}

// TODO: jsonを返すようにしよう
func generateJson(testDir string) error {
	cmd := exec.Command("terraform", "init")
	cmd.Dir = testDir
	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, "run terraform")
	}

	cmd = exec.Command("terraform", "plan", "-out=plan")
	cmd.Dir = testDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrap(err, string(out))
	}

	cmd = exec.Command("terraform", "show", "-json", "plan")
	cmd.Dir = testDir
	out, err = cmd.CombinedOutput()
	slog.Info("output", "json", string(out))
	if err != nil {
		return errors.Wrap(err, "run terraform")
	}

	return nil

}
