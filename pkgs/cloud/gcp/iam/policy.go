// Package iam provides functionality for interacting with Google Cloud Identity and Access Management (IAM) policies.
package iam

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/api/cloudresourcemanager/v1"
)

// Error messages for policy-related failures.
const (
	ErrorFailedToGetPolicy    = "failed to get policy"
	ErrorFailedToSetPolicy    = "failed to set policy"
	ErrorFailedToDeletePolicy = "failed to delete policy"
)

// SetPolicy sets the IAM policy for the specified project.
func (crmService *CrmSvc) SetPolicy(policy *cloudresourcemanager.Policy) error {
	_, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	request := new(cloudresourcemanager.SetIamPolicyRequest)
	request.Policy = policy
	_, err := crmService.svc.Projects.SetIamPolicy(crmService.projectId, request).Do()
	if err != nil {
		return fmt.Errorf("%s : %w", ErrorFailedToSetPolicy, err)
	}
	return nil
}

// GetIamPolicy retrieves the IAM policy for the specified project.
func (crmService *CrmSvc) GetIamPolicy() (*cloudresourcemanager.Policy, error) {
	policy, err := crmService.svc.Projects.GetIamPolicy(crmService.projectId, &cloudresourcemanager.GetIamPolicyRequest{}).
		Do()

	if err != nil {
		return nil, fmt.Errorf("%s : %w", ErrorFailedToGetPolicy, err)
	}
	return policy, nil
}

// DeletePolicy deletes the IAM policy for the specified project.
func (crmService *CrmSvc) DeletePolicy() error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	_, err := crmService.svc.Projects.SetIamPolicy(crmService.projectId, &cloudresourcemanager.SetIamPolicyRequest{}).
		Do()
	if err != nil {
		return fmt.Errorf("%s : %w", ErrorFailedToDeletePolicy, err)
	}
	return nil
}
