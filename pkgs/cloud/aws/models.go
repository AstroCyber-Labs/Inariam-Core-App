package aws

import (
	"gitea/pcp-inariam/inariam/pkgs/cloud/aws/cognito"
	inaIam "gitea/pcp-inariam/inariam/pkgs/cloud/aws/iam"
	"github.com/aws/aws-sdk-go/aws/session"
)

type Session struct {
	ClientSession     *session.Session
	AccessAnalyzerSvc *AccessAnalyzerSvc
	CloudTrailSvc     *CloudTrailSvc
	GuarddutySvc      *GuardDutySvc
	IamSvc            *inaIam.Svc
	SecurityHubSvc    *SecurityHubSvc
	CognitoSvc        *cognito.Svc
}

type Credentials struct {
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
	Region          string
}
