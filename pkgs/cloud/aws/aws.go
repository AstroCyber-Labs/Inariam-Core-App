package aws

import (
	"fmt"
	"gitea/pcp-inariam/inariam/pkgs/log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
)

func OpenSession(creds *Credentials) (*Session, error) {
	var sess *session.Session

	var err error

	fmt.Println(creds.Region)

	if creds == nil {
		sess, err = session.NewSession(&aws.Config{
			Region: aws.String(creds.Region),
		})
	} else {
		sess, err = session.NewSession(&aws.Config{
			Region:      aws.String(creds.Region),
			Credentials: credentials.NewStaticCredentials(creds.AccessKeyID, creds.SecretAccessKey, creds.SessionToken),
		})
	}

	if err != nil {
		return nil, err
	}

	sess.Handlers.Send.PushFront(func(r *request.Request) {
		// Log every request made and its payload
		log.Logger.Infof("Request: %s/%v, Payload: %s",
			r.ClientInfo.ServiceName, r.Operation, r.Params)
	})
	return &Session{
		ClientSession: sess,
	}, nil
}
