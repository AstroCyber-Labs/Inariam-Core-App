package aws

import (
	cognitoSvc "gitea/pcp-inariam/inariam/pkgs/cloud/aws/cognito"

	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

func (awsSession *Session) CreateCognitoSvc(cognitoAppClientID, appSecret string) {
	if awsSession.CognitoSvc == nil {

		awsSession.CognitoSvc = &cognitoSvc.Svc{
			Svc:         cognito.New(awsSession.ClientSession),
			AppClientId: cognitoAppClientID,
			AppSecret:   appSecret,
		}
	}
}
