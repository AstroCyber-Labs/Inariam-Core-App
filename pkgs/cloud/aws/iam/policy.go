// Package iam provides functionality for managing AWS Identity and Access Management (IAM) policies.
package iam

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/iam"

	"gitea/pcp-inariam/inariam/pkgs/log"
)

// Constants for IAM errors.
const (
	// ErrIamPolicyExists is an error indicating that the IAM policy already exists.
	ErrIamPolicyExists = "error IAM policy already exists"
	// ErrIamPolicyNotExists is an error indicating that the IAM policy does not exist.
	ErrIamPolicyNotExists = "error IAM policy does not exist"
	// ErrMarshallingPolicy is an error indicating that marshalling IAM policy failed.
	ErrMarshallingPolicy = "error marshalling IAM policy failed"
)

// PolicyDocument represents the structure of an IAM policy document.
type PolicyDocument struct {
	Version   string
	Statement []StatementEntry
}

// StatementEntry represents an entry in an IAM policy statement.
type StatementEntry struct {
	Effect   string
	Action   string
	Resource string
}

// GetIamPolicy retrieves an IAM policy using its ARN.
// Returns the IAM policy or nil if it doesn't exist.
// TODO build the complete policyARN path using the account ID
func (IamSvc *Svc) GetIamPolicy(policyARN string) (*iam.Policy, error) {
	getPolicyInput := &iam.GetPolicyInput{
		PolicyArn: aws.String(policyARN),
	}

	policyOutput, err := IamSvc.svc.GetPolicy(getPolicyInput)
	if err != nil {
		var iamErr awserr.Error
		if errors.As(err, &iamErr) {
			if iamErr.Code() == iam.ErrCodeNoSuchEntityException {
				return nil, nil
			}
		}
		return nil, fmt.Errorf("CheckIfIamPolicyExists: %w", err)
	}

	return policyOutput.Policy, nil
}

// CreateIamPolicy creates a new IAM policy with the given name, description, and policy document.
// Returns the created IAM policy or an error if creation fails.
func (IamSvc *Svc) CreateIamPolicy(policyName string, description string, policy PolicyDocument) (*iam.Policy, error) {
	// We check with the ARB of the policy (Change the AccountID).
	policyDetails, err := IamSvc.GetIamPolicy("arn:aws:iam::" + "782390994097:policy/" + policyName)

	if policyDetails != nil {
		return nil, fmt.Errorf("CreateIamPolicy: %s %w", ErrIamPolicyExists, err)
	}

	result, err := json.Marshal(policy)
	if err != nil {
		return nil, fmt.Errorf("CreateIamPolicy: %s %w", ErrMarshallingPolicy, err)
	}

	log.Logger.Infoln(string(result))

	createPolicyInput := &iam.CreatePolicyInput{
		PolicyDocument: aws.String(string(result)),
		PolicyName:     aws.String(policyName),
		Description:    aws.String(description),
	}

	outputPolicy, err := IamSvc.svc.CreatePolicy(createPolicyInput)
	if err != nil {
		return nil, fmt.Errorf("CreateIamPolicy: %w", err)
	}

	log.Logger.Infof("IAM policy '%s' created successfully\n", policyName)
	return outputPolicy.Policy, nil
}

// DeleteIamPolicy deletes an IAM policy using its ARN.
// Returns an error if deletion fails or if the policy doesn't exist.
func (IamSvc *Svc) DeleteIamPolicy(policyARN string) error {
	policyDetails, err := IamSvc.GetIamPolicy(policyARN)
	if policyDetails == nil {
		return fmt.Errorf("DeletePolicy: %s %w", ErrIamPolicyNotExists, err)
	}

	// Delete all versions of the policy
	output, err := IamSvc.svc.ListPolicyVersions(&iam.ListPolicyVersionsInput{
		PolicyArn: &policyARN,
	})
	if err != nil {
		return fmt.Errorf("DeletePolicy: %w", err)
	}

	for _, version := range output.Versions {
		if *version.IsDefaultVersion {
			continue // Skip the default version
		}

		deletePolicyVersionInput := &iam.DeletePolicyVersionInput{
			PolicyArn: aws.String(policyARN),
			VersionId: version.VersionId,
		}

		_, err = IamSvc.svc.DeletePolicyVersion(deletePolicyVersionInput)
		if err != nil {
			return fmt.Errorf("DeletePolicy: %w", err)
		}

		log.Logger.Infof("IAM policy version '%s' deleted successfully\n", *version.VersionId)
	}

	// Delete the policy
	deletePolicyInput := &iam.DeletePolicyInput{
		PolicyArn: aws.String(policyARN),
	}

	_, err = IamSvc.svc.DeletePolicy(deletePolicyInput)
	if err != nil {
		return fmt.Errorf("DeletePolicy: %w", err)
	}

	log.Logger.Infof("IAM policy '%s' deleted successfully\n", policyARN)
	return nil
}

// ListIamPolicies returns a list of IAM policies.
func (IamSvc *Svc) ListIamPolicies() ([]*iam.Policy, error) {
	listPoliciesInput := &iam.ListPoliciesInput{
		MaxItems: aws.Int64(10),
	}

	listPoliciesOutput, err := IamSvc.svc.ListPolicies(listPoliciesInput)
	if err != nil {
		return nil, fmt.Errorf("ListIamPolicies: %w", err)
	}

	return listPoliciesOutput.Policies, nil
}

// UpdateIamPolicy updates the content of an existing IAM policy.
// Returns the updated IAM policy or an error if the update fails.
func (IamSvc *Svc) UpdateIamPolicy(policyARN string, newPolicy PolicyDocument) (*iam.Policy, error) {
	policyDetails, err := IamSvc.GetIamPolicy(policyARN)
	if err != nil {
		return nil, fmt.Errorf("UpdateIamPolicy: %w", err)
	}
	if policyDetails == nil {
		return nil, errors.New(ErrIamPolicyNotExists)
	}

	result, err := json.Marshal(newPolicy)
	if err != nil {
		return nil, fmt.Errorf("UpdateIamPolicy: %s", ErrMarshallingPolicy)
	}

	// Create a new version of the policy
	createPolicyVersionInput := &iam.CreatePolicyVersionInput{
		PolicyArn:      aws.String(policyARN),
		PolicyDocument: aws.String(string(result)),
		SetAsDefault:   aws.Bool(true), // Set the new version as the default version
	}

	_, err = IamSvc.svc.CreatePolicyVersion(createPolicyVersionInput)
	if err != nil {
		return nil, fmt.Errorf("UpdateIamPolicy: %w", err)
	}

	// Fetch and return updated policy details
	getPolicyInput := &iam.GetPolicyInput{
		PolicyArn: aws.String(policyARN),
	}

	getPolicyOutput, err := IamSvc.svc.GetPolicy(getPolicyInput)
	if err != nil {
		return nil, fmt.Errorf("UpdateIamPolicy: %w", err)
	}

	return getPolicyOutput.Policy, nil
}
