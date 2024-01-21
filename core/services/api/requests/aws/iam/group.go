// Package aws provides structures and functionality related to AWS IAM groups.
package aws

import "github.com/go-playground/validator/v10"

// CreateGroupRequest represents a request to create an IAM group.
type CreateGroupRequest struct {
	GroupName string `json:"name"`
}

// Validate validates the CreateGroupRequest structure using the go-playground/validator library.
func (createGroupRequest *CreateGroupRequest) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	return validate.Struct(createGroupRequest)
}

// GetGroupRequest represents a request to get information about an IAM group.
type GetGroupRequest struct {
	GroupName string `json:"group_name" param:"name"`
}

// Validate validates the GetGroupRequest structure using the go-playground/validator library.
func (awsIamGroupRequest *GetGroupRequest) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	return validate.Struct(awsIamGroupRequest)
}

// UpdateGroupRequest represents a request to update information about an IAM group.
type UpdateGroupRequest struct {
	NewGroupName string `json:"new_name"`
	Id           string `param:"id"`
}

// Validate validates the UpdateGroupRequest structure using the go-playground/validator library.
func (updateGroupRequest *UpdateGroupRequest) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	return validate.Struct(updateGroupRequest)
}
