// Package gcp provides structures and functionality related to Google Cloud Platform (GCP).
package gcp

// HTTP error messages related to GCP IAM.
const (
	HttpErrOpenedSession = "Failed to open IAM session"
	HttpErrRoleNotFound  = "Role not found"
)

// GroupDetails represents details of a GCP IAM group.
type GroupDetails struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Etag        string `json:"etag"`
	Members     []struct {
		Email string `json:"email"`
	} `json:"members"`
	Id string `json:"id"`
}
