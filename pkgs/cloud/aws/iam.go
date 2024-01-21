package aws

import (
	inaIam "gitea/pcp-inariam/inariam/pkgs/cloud/aws/iam"

	"github.com/aws/aws-sdk-go/service/iam"
)

func (awsSess *Session) OpenIamService() {
	if awsSess.IamSvc == nil {
		awsSess.IamSvc = inaIam.New(iam.New(awsSess.ClientSession))
	}
}
