package iam

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ActionRoleRequest struct {
	Name string `param:"name" json:"name"`
}

// CreateRoleRequest create role request
type CreateRoleRequest struct {
	Id          string   `json:"id"          validate:"required"`
	Name        string   `json:"name"        validate:"required"`
	Title       string   `json:"title"       validate:"required"`
	Description string   `json:"description" validate:"required"`
	Stage       *string  `json:"stage"`
	Permissions []string `json:"permissions" validate:"required,dive,required"`
}

type UpdateRoleRequest struct {
	Id          string   `param:"id"`
	Title       string   `json:"title"       validate:"required"`
	Description string   `json:"description" validate:"required"`
	Permissions []string `json:"permissions" validate:"required"`
}

// Validate validates the CreateRoleRequest structure
func (createRoleRequest *CreateRoleRequest) Validate() error {
	validate := validator.New()
	if err := validate.Struct(createRoleRequest); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	return nil
}

// Validate validates the UpdateRoleRequest structure
func (updateRoleRequest *UpdateRoleRequest) Validate() error {
	validate := validator.New()

	if err := validate.Struct(updateRoleRequest); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	return nil
}

// Validate validates the UpdateRoleRequest structure
func (actionRoleRequest *ActionRoleRequest) Validate() error {
	validate := validator.New()

	if err := validate.Struct(actionRoleRequest); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	return nil
}
