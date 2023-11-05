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
		name            string
		repoURL         string
		srcBranch       string
		terraformSrcDir string
		dstBranch       string
		terraformDstDir string
		isEqual         bool
	}{
		{
			name:            "compare branches having diffelent changes",
			repoURL:         "https://github.com/paper2/terradiff",
			srcBranch:       "main",
			terraformSrcDir: "resource/test/terraform/src",
			dstBranch:       "main",
			terraformDstDir: "resource/test/terraform/dst",
			isEqual:         false,
		},
		{
			name:            "compare branches having same changes",
			repoURL:         "https://github.com/paper2/terradiff",
			srcBranch:       "main",
			terraformSrcDir: "resource/test/terraform/dst",
			dstBranch:       "main",
			terraformDstDir: "resource/test/terraform/dst",
			isEqual:         true,
		},
	}

	for _, tc := range tests {
		tc := tc
		// TODO: terradiff関数を直接使って作成されたjsonの内容を確認するテストにする
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			workDir := t.TempDir()

			srcDir := filepath.Join(workDir, "src")
			dstDir := filepath.Join(workDir, "dst")
			srcGit := NewGit(NewCommandExecutor("."), tc.repoURL, srcDir, tc.srcBranch)
			dstGit := NewGit(NewCommandExecutor("."), tc.repoURL, dstDir, tc.dstBranch)

			srcTerraformDir := filepath.Join(srcDir, tc.terraformSrcDir)
			dstTerraformDir := filepath.Join(dstDir, tc.terraformDstDir)
			srcTerraform := NewTerraform(NewCommandExecutor(srcTerraformDir))
			dstTerrafom := NewTerraform(NewCommandExecutor(dstTerraformDir))

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
