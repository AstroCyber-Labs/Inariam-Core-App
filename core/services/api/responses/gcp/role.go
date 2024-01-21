// Package gcp provides structures and functionality related to Google Cloud Platform (GCP).
package gcp

// RoleResponse represents details of a GCP IAM role.
type RoleResponse struct {
	Name        string   `json:"name"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Stage       string   `json:"stage"`
	Permissions []string `json:"permissions"`
}

// ListRolesResponse represents a response containing a list of GCP IAM roles.
type ListRolesResponse struct {
	Roles []*RoleResponse `json:"roles"`
}
