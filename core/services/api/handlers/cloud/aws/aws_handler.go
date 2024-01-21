package aws

import (
	"gitea/pcp-inariam/inariam/core/services/api"
	inaAws "gitea/pcp-inariam/inariam/pkgs/cloud/aws"
)

type Handler struct {
	api *api.API
}

func NewAwsHandler(api *api.API) *Handler {
	return &Handler{api}
}

// RetrieveAwsSession retrieve an IAM AwsSession with the config passed to the handler
func (awsHandler *Handler) RetrieveAwsIamSession() (*inaAws.Session, error) {
	creds := inaAws.Credentials{
		AccessKeyID:     awsHandler.api.Config.AWS.AccessKeyID,
		SecretAccessKey: awsHandler.api.Config.AWS.SecretAccessKey,
	}
	awsSession, err := inaAws.OpenSession(&creds)

	if err != nil {
		return nil, err
	}

	awsSession.OpenIamService()
	return awsSession, nil
}
