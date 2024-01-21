// Package iam provides services for managing Identity and Access Management (IAM) operations in Google Cloud.
package iam

import (
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/iam/v1"
)

// IamSvc represents the Identity and Access Management (IAM) service.
type IamSvc struct {
	svc       *iam.Service
	projectId string
}

// CrmSvc represents the Cloud Resource Manager (CRM) service.
type CrmSvc struct {
	svc       *cloudresourcemanager.Service
	projectId string
}

// AdminSvc represents the Admin service.
type AdminSvc struct {
	svc *admin.Service
}

// NewIam creates a new IAM service with the provided IAM service client and project ID.
func NewIam(svc *iam.Service, projectId string) *IamSvc {
	return &IamSvc{svc: svc, projectId: projectId}
}

// NewCrm creates a new CRM service with the provided CRM service client and project ID.
func NewCrm(svc *cloudresourcemanager.Service, projectId string) *CrmSvc {
	return &CrmSvc{svc: svc, projectId: projectId}
}

// NewAdminSvc creates a new Admin service with the provided Admin service client.
func NewAdminSvc(svc *admin.Service) *AdminSvc {
	return &AdminSvc{
		svc: svc,
	}
}
