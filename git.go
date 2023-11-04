package main

import (
	"context"
	"os"
)

type Git struct {
	r        Runner
	repoURL  string
	cloneDir string
	branch   string
}

func NewGit(runner Runner, repoURL, cloneDir, branch string) *Git {
	return &Git{r: runner, repoURL: repoURL, cloneDir: cloneDir, branch: branch}
}

func (g *Git) GitClone(ctx context.Context) error {
	if _, err := os.Stat(g.cloneDir); !os.IsNotExist(err) {
		Logger().Info("skip clone repository. there already exists.", "repository", g.repoURL, "branch", g.branch)
		return nil
	}

	err := g.r.RunContext(ctx, "git", "clone", "-b", g.branch, g.repoURL, g.cloneDir)
	if err != nil {
		return err
	}
	return nil
}
