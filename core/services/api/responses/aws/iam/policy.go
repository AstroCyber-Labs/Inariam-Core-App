// Package iam provides structures and functionality related to AWS Identity and Access Management (IAM) policies.
package iam

// Error messages related to IAM policies.
const (
	ErrorMissingPolicyARN         = "Policy ARN is required"
	ErrorMissingPolicyDocument    = "Policy document is required"
	ErrorMarshalingPolicyDocument = "Error marshaling policy document"
	ErrorParsingPolicyDocument    = "Error parsing policy document"
)

// PolicyDetailResponse represents a response detailing an IAM policy.
type PolicyDetailResponse struct {
	PolicyName string `json:"policy_name"`
	PolicyID   string `json:"policy_id"`
}
