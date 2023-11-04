package main

import (
	"context"
	"encoding/json"
)

// TODO: 関数もう少し小さく分ける。structとinteraceも用意する？
func generatePlanResult(ctx context.Context, planDir string) (*PlanResult, error) {
	ce := NewCommandExecutor(planDir)

	err := ce.RunContext(ctx, "terraform", "init")
	if err != nil {
		return nil, err
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
	ResourceChanges []ResourceChange `json:"resource_changes"`
}

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
