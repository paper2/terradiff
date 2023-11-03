package main

import (
	"github.com/cockroachdb/errors"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type Git struct {
	dir string
	wt  *git.Worktree
}

func NewGit(dir string) (*Git, error) {
	r, err := git.PlainOpen(dir)
	if err != nil {
		return nil, errors.Wrap(err, "open git repo")
	}
	w, err := r.Worktree()
	if err != nil {
		return nil, errors.Wrap(err, "get worktree")
	}
	return &Git{dir: dir, wt: w}, nil
}

func (g *Git) Checkout(branch string) error {
	status, err := g.wt.Status()
	if err != nil {
		return errors.Wrap(err, "git status")
	}
	if !status.IsClean() {
		return errors.New("working tree is not clean. commit or stash changes.")
	}

	err = g.wt.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branch),
	})
	if err != nil {
		return errors.Wrap(err, "checkout")
	}
	Logger().Debug("checkout: " + branch)

	return nil
}
