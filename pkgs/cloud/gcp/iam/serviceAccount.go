// Package iam provides functionality for managing Identity and Access Management (IAM) service accounts in Google Cloud.
package iam

import (
	"fmt"

	"google.golang.org/api/iam/v1"

	"gitea/pcp-inariam/inariam/pkgs/log"
)

// Error messages for IAM service account-related failures.
const (
	ErrorFailedToCreateServiceAccount  = "failed to create service account"
	ErrorFailedToEnableServiceAccount  = "failed to enable service account"
	ErrorFailedToDisableServiceAccount = "failed to disable service account"
	ErrorFailedToDeleteServiceAccount  = "failed to delete service account"
	ErrorFailedToListServiceAccounts   = "failed to list service accounts"
)

// GetServiceAccount retrieves information about an IAM service account by its name.
func (iamService *IamSvc) GetServiceAccount(
	name string,
) (*iam.ServiceAccount, error) {
	svcAccount, err := iamService.svc.Projects.ServiceAccounts.Get("projects/" + iamService.projectId + "/serviceAccounts/" + name).
		Do()
	if err != nil {
		return nil, err
	}
	return svcAccount, nil
}

// CreateIamServiceAccount creates a new IAM service account.
func (iamService *IamSvc) CreateIamServiceAccount(
	displayName, name, description string,
) (*iam.ServiceAccount, error) {
	request := &iam.CreateServiceAccountRequest{
		AccountId: name,
		ServiceAccount: &iam.ServiceAccount{
			DisplayName: displayName,
			Description: description,
		},
	}
	account, err := iamService.svc.Projects.ServiceAccounts.Create("projects/"+iamService.projectId, request).
		Do()
	if err != nil {
		return nil, fmt.Errorf("%s %w", ErrorFailedToCreateServiceAccount, err)
	}
	log.Logger.Info("Created service account: %s", account.Email)
	return account, nil
}

// DeleteIamServiceAccount deletes an IAM service account by its email.
func (iamService *IamSvc) DeleteIamServiceAccount(email string) error {
	_, err := iamService.svc.Projects.ServiceAccounts.Delete("projects/" + iamService.projectId + "/serviceAccounts/" + email).
		Do()
	if err != nil {
		return fmt.Errorf("%s %w", ErrorFailedToDeleteServiceAccount, err)
	}
	log.Logger.Info("Deleted service account: %s", email)
	return nil
}

// ListIamServiceAccounts retrieves a list of all IAM service accounts in the project.
func (iamService *IamSvc) ListIamServiceAccounts() ([]*iam.ServiceAccount, error) {
	response, err := iamService.svc.Projects.ServiceAccounts.List("projects/" + iamService.projectId).
		Do()
	if err != nil {
		return nil, fmt.Errorf("%s %w", ErrorFailedToListServiceAccounts, err)
	}

	for _, account := range response.Accounts {
		log.Logger.Infof("Listing service account: %v\n", account.Name)
	}

	return response.Accounts, nil
}

// EnableIamServiceAccount enables an IAM service account by its name.
func (iamService *IamSvc) EnableIamServiceAccount(name string) error {
	request := &iam.EnableServiceAccountRequest{}
	_, err := iamService.svc.Projects.ServiceAccounts.Enable("projects/"+iamService.projectId+"/serviceAccounts/"+name, request).
		Do()

	if err != nil {
		return fmt.Errorf("%s %w", ErrorFailedToEnableServiceAccount, err)
	}
	log.Logger.Info("Enabled service account: %s", name)
	return nil
}

// DisableIamServiceAccount disables an IAM service account by its name.
func (iamService *IamSvc) DisableIamServiceAccount(name string) error {

	request := &iam.DisableServiceAccountRequest{}
	_, err := iamService.svc.Projects.ServiceAccounts.Disable("projects/"+iamService.projectId+"/serviceAccounts/"+name, request).
		Do()
	if err != nil {
		return fmt.Errorf("%s %w", ErrorFailedToDisableServiceAccount, err)
	}
	log.Logger.Info("Disabled service account: %s", name)
	return nil
}
