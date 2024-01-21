package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/accessanalyzer"

	"gitea/pcp-inariam/inariam/pkgs/log"
)

type AccessAnalyzerSvc struct {
	svc *accessanalyzer.AccessAnalyzer
}

func (awsSess *Session) CreateAccessAnalyzerSvc() {
	if awsSess.AccessAnalyzerSvc == nil {
		awsSess.AccessAnalyzerSvc = &AccessAnalyzerSvc{accessanalyzer.New(awsSess.ClientSession)}
	}
}

func (accessanalyzerSvc *AccessAnalyzerSvc) ListAccessAnalyzers() {
	resp, err := accessanalyzerSvc.svc.ListAnalyzers(&accessanalyzer.ListAnalyzersInput{})

	if err != nil {
		log.Logger.Fatalf("Failed to list analyzers: %s", err)
	}

	for _, analyzer := range resp.Analyzers {
		log.Logger.Infoln("Analyzer ARN:", *analyzer.Arn)
		log.Logger.Infoln("Analyzer Name:", *analyzer.Name)
		// ... add more fields as needed
	}
}

func (accessanalyzerSvc *AccessAnalyzerSvc) ListFindings(analyzerArn string) {

	input := &accessanalyzer.ListFindingsInput{
		AnalyzerArn: aws.String(analyzerArn),
	}

	resp, err := accessanalyzerSvc.svc.ListFindings(input)
	if err != nil {
		log.Logger.Fatalf("Failed to list findings: %s", err)
	}

	for _, findingId := range resp.Findings {
		log.Logger.Infoln("Finding ID:", *findingId.Id)
		// ... add more fields as needed
	}
}
