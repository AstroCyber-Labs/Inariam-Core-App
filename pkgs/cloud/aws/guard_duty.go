package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/service/guardduty"
)

type GuardDutySvc struct {
	svc *guardduty.GuardDuty
}

func (awsSess *Session) CreateGuarddutySvc() {
	if awsSess.GuarddutySvc == nil {
		awsSess.GuarddutySvc = &GuardDutySvc{guardduty.New(awsSess.ClientSession)}
	}
}

func (GuarddutySvc *GuardDutySvc) ListFindings() {

	// Initialize a session using Amazon SDK with a region closest to you

	// Create a new instance of the GuardDuty service's client with a Session.

	// Get the DetectorId (needed for other GuardDuty operations)
	detectorResp, err := GuarddutySvc.svc.ListDetectors(&guardduty.ListDetectorsInput{})
	if err != nil {
		log.Fatalf("Failed to list detectors: %s", err)
	}
	if len(detectorResp.DetectorIds) == 0 {
		log.Fatalf("No GuardDuty detectors found")
	}
	detectorID := detectorResp.DetectorIds[0] // Taking the first detector

	fmt.Printf("Listing findings for detector %s", *detectorID)

	// List findings
	findingsResp, err := GuarddutySvc.svc.ListFindings(&guardduty.ListFindingsInput{
		DetectorId: detectorID,
	})
	if err != nil {
		log.Fatalf("Failed to list findings: %s", err)
	}

	// Get details for each finding
	for _, findingId := range findingsResp.FindingIds {
		getFindingsResp, err := GuarddutySvc.svc.GetFindings(&guardduty.GetFindingsInput{
			DetectorId: detectorID,
			FindingIds: []*string{findingId},
		})
		if err != nil {
			log.Printf("Failed to get finding details for %s: %s", *findingId, err)
			continue
		}

		fmt.Printf("%v", getFindingsResp.Findings)

		for _, finding := range getFindingsResp.Findings {
			fmt.Println("Title:", *finding.Title)
			fmt.Println("Description:", *finding.Description)
			// Add more fields as required
			fmt.Println("---------------------")
		}
	}
}
