package main

import (
	"context"
	"os"
)

type Git struct {
	repoURL  string
	cloneDir string
	branch   string
}

func NewGit(repoURL, cloneDir, branch string) *Git {
	return &Git{repoURL: repoURL, cloneDir: cloneDir, branch: branch}
}

func (g *Git) gitClone(ctx context.Context) error {
	if _, err := os.Stat(g.cloneDir); !os.IsNotExist(err) {
		Logger().Info("skip clone repository.", "repository", g.repoURL, "branch", g.branch)
		return nil
	}

	err := NewCommandExecutor(".").RunContext(ctx, "git", "clone", "-b", g.branch, g.repoURL, g.cloneDir)
	if err != nil {
		return err
	}
	return nil
}

func (g *Git) getCloneDir() string {
	return g.cloneDir
}
