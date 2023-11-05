package main

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTerradiff_Compare(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		srcBranch string
		dstBranch string
		repoURL   string
		isEqual   bool
	}{
		{
			name:      "compare branches having diffelent changes",
			srcBranch: "branch1",
			dstBranch: "main",
			repoURL:   "https://github.com/paper2/test-terradiff",
			isEqual:   false,
		},
		{
			name:      "compare same branches",
			srcBranch: "main",
			dstBranch: "main",
			repoURL:   "https://github.com/paper2/test-terradiff",
			isEqual:   true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			workDir := t.TempDir()

			srcDir := filepath.Join(workDir, "src")
			dstDir := filepath.Join(workDir, "dst")
			srcGit := NewGit(NewCommandExecutor("."), tc.repoURL, srcDir, tc.srcBranch)
			dstGit := NewGit(NewCommandExecutor("."), tc.repoURL, dstDir, tc.dstBranch)

			srcTerraform := NewTerraform(NewCommandExecutor(srcDir))
			dstTerrafom := NewTerraform(NewCommandExecutor(dstDir))

			terradiff := NewTerradiff(
				Operators{tf: srcTerraform, git: srcGit},
				Operators{tf: dstTerrafom, git: dstGit},
			)

			cr, err := terradiff.Compare(context.Background())
			assert.NoError(t, err)
			assert.Equal(t, tc.isEqual, cr.IsEqual, "IsEqual should match expected value for test case: "+tc.name)
		})
	}
}
