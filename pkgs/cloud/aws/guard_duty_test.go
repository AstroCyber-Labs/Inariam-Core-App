package aws_test

import (
	"gitea/pcp-inariam/inariam/pkgs/cloud/aws"
	"testing"
)

func TestGuardDuty(t *testing.T) {
	// Initialize AWS session

	credentials := aws.GetTestCredentials()

	awsSess, err := aws.OpenSession(&credentials)
	if err != nil {
		panic("guardduty_test failed, error opening session")
	}

	awsSess.CreateGuarddutySvc()

	// Ensure there's a detector before testing, or you can create one if none exists.
	awsSess.GuarddutySvc.ListFindings()
}
