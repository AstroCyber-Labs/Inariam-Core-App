package iam

import (
	"github.com/aws/aws-sdk-go/service/iam"
)

type Svc struct {
	svc *iam.IAM
}

func New(svc *iam.IAM) *Svc {
	return &Svc{svc: svc}
}
