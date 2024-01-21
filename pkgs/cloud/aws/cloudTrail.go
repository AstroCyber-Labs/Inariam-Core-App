package aws

import (
	"gitea/pcp-inariam/inariam/pkgs/log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudtrail"
)

type CloudTrailSvc struct {
	svc *cloudtrail.CloudTrail
}

func (awsSess *Session) CreateCloudTrailSvc() {
	if awsSess.CloudTrailSvc == nil {
		awsSess.CloudTrailSvc = &CloudTrailSvc{cloudtrail.New(awsSess.ClientSession)}
	}

}

func (CloudTrailSvc *CloudTrailSvc) DescribeTrails() {

	resp, err := CloudTrailSvc.svc.DescribeTrails(&cloudtrail.DescribeTrailsInput{})
	if err != nil {
		log.Logger.Fatalf("Failed to describe trails: %s", err)
	}

	for _, trail := range resp.TrailList {
		log.Logger.Infoln("Trail ARN:", *trail.TrailARN)
		log.Logger.Infoln("Home Region:", *trail.HomeRegion)
		log.Logger.Infoln("Log File Validation Enabled:", *trail.LogFileValidationEnabled)

		// ... add more fields as needed
	}
}

func (CloudTrailSvc *CloudTrailSvc) LookupEvents() {

	input := &cloudtrail.LookupEventsInput{
		LookupAttributes: []*cloudtrail.LookupAttribute{
			{
				AttributeKey:   aws.String("EventName"),
				AttributeValue: aws.String("DescribeInstances"),
			},
		},
	}

	resp, err := CloudTrailSvc.svc.LookupEvents(input)
	if err != nil {
		log.Logger.Fatalf("Failed to lookup events: %s", err)
	}

	for _, event := range resp.Events {
		log.Logger.Infoln("Event ID:", *event.EventId)
		log.Logger.Infoln("Event Name:", *event.EventName)
		// ... add more fields as needed
	}
}
