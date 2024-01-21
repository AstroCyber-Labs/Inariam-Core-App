// Package iam provides structures and functionality related to AWS Identity and Access Management (IAM).
package iam

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// ActionOneServiceAccountRequest represents a request to perform an action on a service account.
type ActionOneServiceAccountRequest struct {
	Name string `param:"name" validate:"required"`
}

// Validate validates the ActionOneServiceAccountRequest structure using the go-playground/validator library.
func (actionOneServiceAccountRequest *ActionOneServiceAccountRequest) Validate() error {
	validate := validator.New()

	if err := validate.Struct(actionOneServiceAccountRequest); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	return nil
}

// CreateServiceAccountRequest represents a request to create a service account.
type CreateServiceAccountRequest struct {
	DisplayName string `json:"display_name" validate:"required"`
	Name        string `json:"name"         validate:"required"`
	Description string `json:"description"  validate:"required"`
}

// Validate validates the CreateServiceAccountRequest structure using the go-playground/validator library.
func (createServiceAccountRequest *CreateServiceAccountRequest) Validate() error {
	validate := validator.New()

	if err := validate.Struct(createServiceAccountRequest); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	return nil
}
