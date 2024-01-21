package iam

import (
	"github.com/go-playground/validator/v10"
)

// ActionGroupRequest a request that represents a unique identifier for groups
type ActionGroupRequest struct {
	Name string `param:"name" validate:"required"`
}

// Validate this method is to validate the actionGroupReq
func (actionGroupReq *ActionGroupRequest) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	return validate.Struct(actionGroupReq)
}

// CreateGroupRequest the request body to create a  group
type CreateGroupRequest struct {
	Name        string `json:"name"        validate:"required"`
	Description string `json:"description" validate:"required"`
}

func (createGroupReq *CreateGroupRequest) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	return validate.Struct(createGroupReq)
}

// UpdateGroupRequest
type UpdateGroupRequest struct {
	Name        string `param:"name" validate:"required"`
	Description string `             validate:"required" json:"description"`
}

func (updatedGroupReq *UpdateGroupRequest) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	return validate.Struct(updatedGroupReq)
}
