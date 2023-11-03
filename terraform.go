package main

import (
	"context"
	"encoding/json"
)

// TODO: jsonを返すようにしよう
func generatePlanResult(ctx context.Context, testDir string) (*PlanResult, error) {
	ce := NewCommandExecutor(testDir)

	err := ce.RunContext(ctx, "terraform", "init")
	if err != nil {
		return nil, err
		// return errors.Wrap(err, "terraform init")
	}

	binaryName := "plan-result-binary"
	err = ce.RunContext(ctx, "terraform", "plan", "-out="+binaryName)
	if err != nil {
		return nil, err
	}

	out, err := ce.RunContextAndCaptureOutput(ctx, "terraform", "show", "-json", binaryName)
	if err != nil {
		return nil, err
	}
	Logger().Debug(out)

	var pr PlanResult
	err = json.Unmarshal([]byte(out), &pr)
	if err != nil {
		return nil, err
	}

	return &pr, nil
}

type PlanResult struct {
	// FormatVersion      string                  `json:"format_version"`
	// PriorState         json.RawMessage         `json:"prior_state"`
	// Configuration      json.RawMessage         `json:"configuration"`
	// PlannedValues      json.RawMessage         `json:"planned_values"`
	// ProposedUnknown    json.RawMessage         `json:"proposed_unknown"`
	// Variables          map[string]Variable     `json:"variables"`
	ResourceChanges []ResourceChange `json:"resource_changes"`
	// ResourceDrift      []ResourceChange        `json:"resource_drift"`
	// RelevantAttributes []RelevantAttribute     `json:"relevant_attributes"`
	// OutputChanges map[string]OutputChange `json:"output_changes"`
	// Checks             json.RawMessage         `json:"checks"`
	// Errored            bool                    `json:"errored"`
}

// type Variable struct {
// 	Value string `json:"value"`
// }

type ResourceChange struct {
	Address         string          `json:"address"`
	PreviousAddress string          `json:"previous_address"`
	ModuleAddress   string          `json:"module_address"`
	Mode            string          `json:"mode"`
	Type            string          `json:"type"`
	Name            string          `json:"name"`
	Index           int             `json:"index"`
	Deposed         string          `json:"deposed"`
	Change          json.RawMessage `json:"change"`
	ActionReason    string          `json:"action_reason"`
}

// type RelevantAttribute struct {
// 	Resource  string `json:"resource"`
// 	Attribute string `json:"attribute"`
// }

// type OutputChange struct {
// 	Change json.RawMessage `json:"change"`
// }
