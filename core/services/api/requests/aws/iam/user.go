// Package aws provides structures and functionality related to AWS IAM users.
package aws

import "github.com/go-playground/validator/v10"

// CreateUserRequest represents a request to create an IAM user.
type CreateUserRequest struct {
	Username string `json:"username"`
}

// Validate validates the CreateUserRequest structure using the go-playground/validator library.
func (awsIamUserRequest *CreateUserRequest) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	return validate.Struct(awsIamUserRequest)
}

// DeleteUserRequest represents a request to delete an IAM user.
type DeleteUserRequest struct {
	Username string `param:"id"`
}

// Validate validates the DeleteUserRequest structure using the go-playground/validator library.
func (awsIamUserRequest *DeleteUserRequest) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	return validate.Struct(awsIamUserRequest)
}

// UpdateUserRequest represents a request to update information about an IAM user.
type UpdateUserRequest struct {
	Username    string `param:"id"`
	NewUsername string `json:"new_username"`
}

// Validate validates the UpdateUserRequest structure using the go-playground/validator library.
func (awsIamUserUpdateRequest *UpdateUserRequest) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	return validate.Struct(awsIamUserUpdateRequest)
}
