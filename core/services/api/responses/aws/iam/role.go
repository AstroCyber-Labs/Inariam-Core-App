// Package iam provides structures and functionality related to AWS Identity and Access Management (IAM) roles.
package iam

// HTTP error messages related to IAM roles.
const (
	HttpErrMissingRoleName           = "role name is required"
	HttpErrMissingTrustPolicy        = "trust policy is required"
	HttpErrRoleNotFound              = "role not found"
	ErrorConvertingTrustPolicyToJSON = "error converting trust policy to JSON"
)

// RoleDetailResponse represents a response detailing an IAM role.
type RoleDetailResponse struct {
	RoleName string `json:"role_name"`
	RoleID   string `json:"role_id"`
	RoleArn  string `json:"role_arn"`
}
