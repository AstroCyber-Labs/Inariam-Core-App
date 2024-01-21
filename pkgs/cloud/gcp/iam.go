// Package gcp provides functions for managing Google Cloud Platform (GCP) services and sessions.
package gcp

import (
	"context"
	"errors"
	"fmt"

	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/iam/v1"
	"google.golang.org/api/option"

	IamGcp "gitea/pcp-inariam/inariam/pkgs/cloud/gcp/iam"
)

var (
	// ErrorFailedToCreateCrmService is returned when there is an error creating the Cloud Resource Manager (CRM) service client.
	ErrorFailedToCreateCrmService = errors.New("error failure to create CRM service client")
	// ErrorFailedToCreateIAMService is returned when there is an error creating the Identity and Access Management (IAM) service client.
	ErrorFailedToCreateIAMService = errors.New("error failure to create IAM service client")
	// ErrorFailedToCreateIAMAdminService is returned when there is an error creating the IAM admin service client.
	ErrorFailedToCreateIAMAdminService = errors.New("error failure to create IAM admin service client")
)

// OpenIamService opens a session for the Identity and Access Management (IAM) service.
func (gSession *Session) OpenIamService() error {
	ctx := context.Background()

	if gSession.IamGCPService != nil {
		return nil
	} else {
		newService, err := iam.NewService(ctx, option.WithCredentials(gSession.Credentials))
		if err != nil {
			return fmt.Errorf("%w : %w", ErrorFailedToCreateIAMService, err)
		}
		gSession.IamGCPService = IamGcp.NewIam(newService, gSession.ProjectId)
	}
	return nil
}

// OpenCrmService opens a session for the Cloud Resource Manager (CRM) service.
func (gSession *Session) OpenCrmService() error {
	ctx := context.Background()

	if gSession.CrmGCPService != nil {
		return nil
	} else {
		newService, err := cloudresourcemanager.NewService(ctx, option.WithCredentials(gSession.Credentials))
		if err != nil {
			return fmt.Errorf("%w : %w", ErrorFailedToCreateCrmService, err)
		}
		gSession.CrmGCPService = IamGcp.NewCrm(newService, gSession.ProjectId)
	}
	return nil
}

// OpenAdminService opens a session for the IAM admin service.
func (gSession *Session) OpenAdminService() error {
	ctx := context.Background()
	client, err := admin.NewService(ctx, option.WithCredentials(gSession.Credentials))

	if err != nil {
		return fmt.Errorf("%w : %w", ErrorFailedToCreateIAMAdminService, err)
	}

	gSession.IamAdminGCPService = IamGcp.NewAdminSvc(client)

	return nil
}
