// Package gcp provides functionalities for managing sessions and services related to Google Cloud Platform (GCP).
package gcp

import (
	"golang.org/x/oauth2/google"

	IamGcp "gitea/pcp-inariam/inariam/pkgs/cloud/gcp/iam"
)

// Session represents a session for managing Google Cloud Platform (GCP) services and related components.
type Session struct {
	Credentials        *google.Credentials
	IamGCPService      *IamGcp.IamSvc
	CrmGCPService      *IamGcp.CrmSvc
	IamAdminGCPService *IamGcp.AdminSvc
	ProjectId          string
}
