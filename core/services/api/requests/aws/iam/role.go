// Package aws provides structures and functionality related to AWS IAM roles.
package aws

import "github.com/go-playground/validator/v10"

// CreateRoleRequest represents a request to create an IAM role.
type CreateRoleRequest struct {
	Id          string                 `json:"id" validate:"required"`
	RoleName    string                 `json:"name" validate:"required"`
	TrustPolicy map[string]interface{} `json:"trust_policy" validate:"required"`
}

// Validate validates the CreateRoleRequest structure using the go-playground/validator library.
func (awsIamRoleRequest *CreateRoleRequest) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	return validate.Struct(awsIamRoleRequest)
}

// UpdateRoleRequest represents a request to update information about an IAM role.
type UpdateRoleRequest struct {
	Name        string                 `param:"id" validate:"required"`
	TrustPolicy map[string]interface{} `json:"trust_policy" validate:"required"`
}

// Validate validates the UpdateRoleRequest structure using the go-playground/validator library.
func (awsIamRoleRequest *UpdateRoleRequest) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	return validate.Struct(awsIamRoleRequest)
}
