// Package aws provides structures and functionality related to AWS policies.
package aws

import "github.com/go-playground/validator/v10"

// PolicyDocument represents the structure of an AWS policy document.
type PolicyDocument struct {
	Version    string           `json:"version"`
	Statements []StatementEntry `json:"statements"`
}

// StatementEntry represents an entry in an AWS policy statement.
type StatementEntry struct {
	Effect   string `json:"effect" validate:"required"`
	Action   string `json:"action" validate:"required"`
	Resource string `json:"resource" validate:"required"`
}

// CreatePolicyRequest represents a request to create an AWS policy.
type CreatePolicyRequest struct {
	PolicyName     string         `json:"name" validate:"required"`
	Description    string         `json:"description" validate:"required"`
	PolicyDocument PolicyDocument `json:"document" validate:"required"`
}

// Validate validates the CreatePolicyRequest structure using the go-playground/validator library.
func (createPolicyRequest *CreatePolicyRequest) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	return validate.Struct(createPolicyRequest)
}

// UpdatePolicyRequest represents a request to update an AWS policy.
type UpdatePolicyRequest struct {
	PolicyARN      string         `param:"arn" validate:"required"`
	PolicyDocument PolicyDocument `json:"document"`
}

// Validate validates the UpdatePolicyRequest structure using the go-playground/validator library.
func (updatePolicyRequest *UpdatePolicyRequest) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	return validate.Struct(updatePolicyRequest)
}
