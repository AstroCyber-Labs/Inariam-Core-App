package aws_test

import (
	awsCloud "gitea/pcp-inariam/inariam/pkgs/cloud/aws"
	"testing"
)

func TestListAccessAnalyzers(t *testing.T) {
	// Initialize a session in us-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.

	credentials := awsCloud.GetTestCredentials()

	awsSess, err := awsCloud.OpenSession(&credentials)
	if err != nil {
		panic("Error opening session")
	}
	awsSess.CreateAccessAnalyzerSvc()

	t.Run("TestListFindings", func(t *testing.T) {
		// Similar to above, start by creating a session.

		// For this test, you need an analyzer ARN. Here, you might want to
		// provide a placeholder or use ListAccessAnalyzers to get an ARN
		// analyzerArn := "arn:aws:access-analyzer:us-west-2:123456789012:analyzer/example-analyzer"
		awsSess.AccessAnalyzerSvc.ListAccessAnalyzers()
		//awsSess.accessAnalyzerSvc.listFindings(analyzerArn)
	})

}
