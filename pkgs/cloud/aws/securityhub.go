package aws

import (
	"fmt"

	"gitea/pcp-inariam/inariam/pkgs/log"

	"github.com/aws/aws-sdk-go/service/securityhub"
)

type SecurityHubSvc struct {
	svc *securityhub.SecurityHub
}

func (awsSess *Session) CreateSecurityHubSvc() {

	if awsSess.SecurityHubSvc == nil {
		return
	}

	awsSess.SecurityHubSvc = &SecurityHubSvc{
		securityhub.New(awsSess.ClientSession),
	}

}

func (securityHubSvc *SecurityHubSvc) CheckIfEnabled() (bool, error) {

	resp, err := securityHubSvc.svc.DescribeHub(&securityhub.DescribeHubInput{})
	if err != nil {
		log.Logger.Infoln(err)
		return false, err
	}

	if resp.HubArn == nil {
		return false, err
	}

	fmt.Printf("ARN: %s,  Activation Date: %s, ", *resp.HubArn, *resp.SubscribedAt)
	return true, nil
}
