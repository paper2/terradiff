package main

import (
	"context"
	"encoding/json"
)

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

type Runner interface {
	RunContext(ctx context.Context, name string, args ...string) error
	RunContextAndCaptureOutput(ctx context.Context, name string, args ...string) (string, error)
}

type Terraform struct {
	r Runner
}

func NewTerraform(r Runner) *Terraform {
	return &Terraform{r: r}
}

func (tf *Terraform) init(ctx context.Context) error {
	return tf.r.RunContext(ctx, "terraform", "init")
}

func (tf *Terraform) genPlanBinary(ctx context.Context, path string) error {
	return tf.r.RunContext(ctx, "terraform", "plan", "-out="+path)
}

func (tf *Terraform) unmarshalPlanBinary(ctx context.Context, path string) (*PlanResult, error) {
	out, err := tf.r.RunContextAndCaptureOutput(ctx, "terraform", "show", "-json", path)
	if err != nil {
		return nil, err
	}

	var pr PlanResult
	err = json.Unmarshal([]byte(out), &pr)
	if err != nil {
		return nil, err
	}

	return &pr, nil
}

func (tf Terraform) GenPlanResult(ctx context.Context) (*PlanResult, error) {
	err := tf.init(ctx)
	if err != nil {
		return nil, err
	}

	binaryPath := "plan.binary"
	err = tf.genPlanBinary(ctx, binaryPath)
	if err != nil {
		return nil, err
	}

	return tf.unmarshalPlanBinary(ctx, binaryPath)
}
