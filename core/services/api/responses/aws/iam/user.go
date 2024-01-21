// Package iam provides structures and functionality related to AWS Identity and Access Management (IAM) users.
package iam

// HTTP error messages related to IAM users.
const (
	HttpErrMissingUserName    = "Username is required"
	HttpErrMissingNewUserName = "Group name is required"
)

// UserDetailResponse represents a response detailing an IAM user.
type UserDetailResponse struct {
	Username string `json:"username"`
	UserID   string `json:"id"`
	UserArn  string `json:"arn"`
}
