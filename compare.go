package main

import (
	"encoding/json"
	"os"

	"github.com/cockroachdb/errors"
	"github.com/google/go-cmp/cmp"
)

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
