package gcp

import (
	"context"
	"fmt"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/iam/v1"
)

type Credentials google.Credentials

const (
	ErrorFailedToLoadCredentials    = "failed to load credentials"
	ErrorFailedToConvertCredentials = "failed to convert credentials"
	ErrorInvalidCredentials         = "invalid credentials"
)

// LoadCredentialsFromFile loads GCP credentials from the specified file path and associates them with the given project ID.
func LoadCredentialsFromFile(credentialsPath string, projectId string) (*Session, error) {
	ctx := context.Background()
	credentialsJSON, err := os.ReadFile(credentialsPath)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", ErrorFailedToLoadCredentials, err)
	}

	creds, err := google.CredentialsFromJSON(ctx, credentialsJSON, iam.CloudPlatformScope)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", ErrorFailedToConvertCredentials, err)
	}

	return &Session{
		Credentials: creds,
		ProjectId:   projectId,
	}, nil
}
