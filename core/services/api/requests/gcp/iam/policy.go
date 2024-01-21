// Package iam provides structures and functionality related to AWS Identity and Access Management (IAM).
package iam

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// GetPolicyOptions represents options for retrieving an IAM policy.
type GetPolicyOptions struct {
	RequestedPolicyVersion int64 `json:"requestedPolicyVersion" validate:"required"`
}

// Binding represents a role binding in an IAM policy.
type Binding struct {
	Role    string   `json:"role"    validate:"required"`
	Members []string `json:"members" validate:"required"`
}

// Policy represents an IAM policy.
type Policy struct {
	Bindings []Binding `json:"bindings" validate:"required"`
	Etag     string    `json:"etag"     validate:"required"`
	Version  int64     `json:"version"  validate:"required"`
}

// SetIamPolicyRequest represents a request to set an IAM policy.
type SetIamPolicyRequest struct {
	Policy Policy `json:"policy" validate:"required"`
}

// Validate validates the SetIamPolicyRequest structure using the go-playground/validator library.
func (setIamPolicyRequest *SetIamPolicyRequest) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	for i, binding := range setIamPolicyRequest.Policy.Bindings {
		err := validate.Struct(binding)
		if err != nil {
			return fmt.Errorf("binding at index %d is invalid: %w", i, err)
		}
	}

	return validate.Struct(setIamPolicyRequest)
}

// GetIamPolicyRequest represents a request to get an IAM policy.
type GetIamPolicyRequest struct {
	Options *GetPolicyOptions `json:"options" validate:"required"`
}

// Validate validates the GetIamPolicyRequest structure using the go-playground/validator library.
func (getIamPolicyRequest *GetIamPolicyRequest) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	return validate.Struct(getIamPolicyRequest)
}

// TestIamPermissionsRequest represents a request to test IAM permissions.
type TestIamPermissionsRequest struct {
	Permissions []string `json:"permissions" validate:"required"`
}
