// Package gcp provides structures and functionality related to Google Cloud Platform (GCP).
package gcp

// ServiceAccountDetails represents details of a GCP IAM service account.
type ServiceAccountDetails struct {
	Description string `json:"description"`
	Disabled    bool   `json:"disabled"`
	Etag        string `json:"etag"`
	Name        string `json:"name"`
	ProjectId   string `json:"project_id"`
	DisplayName string `json:"displayName"`
	UniqueId    string `json:"uniqueId"`
}
