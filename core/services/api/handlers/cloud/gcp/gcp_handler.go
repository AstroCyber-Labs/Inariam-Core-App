package gcp

import (
	"fmt"

	"gitea/pcp-inariam/inariam/core/services/api"
	inaGCP "gitea/pcp-inariam/inariam/pkgs/cloud/gcp"
)

type Handler struct {
	api *api.API
	// TODO: add gcpSession variable to not open the session on each request.
}

func NewGCPHandler(api *api.API) *Handler {
	return &Handler{api}
}

// RetrieveAwsSession retrieve an IAM AwsSession with the config passed to the handler
func (handler *Handler) retrieveGCPSession() (*inaGCP.Session, error) {

	gcpSession, err := inaGCP.LoadCredentialsFromFile(
		handler.api.Config.GCP.CredentialsPath,
		handler.api.Config.GCP.ProjectID,
	)

	if err != nil {
		return nil, fmt.Errorf("retrieveGCPSession: %w", err)
	}

	return gcpSession, nil
}

func (handler *Handler) openAdminIamSession() (*inaGCP.Session, error) {
	gcpSession, err := handler.retrieveGCPSession()
	if err != nil {
		return nil, fmt.Errorf("openAdminIamSession: %w", err)
	}

	err = gcpSession.OpenAdminService()
	if err != nil {
		return nil, fmt.Errorf("openAdminIamSession: %w", err)
	}

	return gcpSession, nil
}

func (handler *Handler) openIamSession() (*inaGCP.Session, error) {
	gcpSession, err := handler.retrieveGCPSession()
	if err != nil {
		return nil, fmt.Errorf("openAdminIamSession: %w", err)
	}

	err = gcpSession.OpenIamService()
	if err != nil {
		return nil, fmt.Errorf("openAdminIamSession: %w", err)
	}

	return gcpSession, nil
}

func (handler *Handler) openCrmSession() (*inaGCP.Session, error) {
	gcpSession, err := handler.retrieveGCPSession()
	if err != nil {
		return nil, fmt.Errorf("openAdminIamSession: %w", err)
	}

	err = gcpSession.OpenCrmService()
	if err != nil {
		return nil, fmt.Errorf("openAdminIamSession: %w", err)
	}

	return gcpSession, nil
}
