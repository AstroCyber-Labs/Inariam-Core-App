// Package cognito provides a service structure for interacting with Amazon Cognito Identity Provider.
/*
 Package cognito provides functionality for working with Amazon Cognito Identity Provider.
 It includes features such as managing user pools, user authentication, and multi-factor authentication (MFA).
 For more information, refer to the official Amazon Cognito documentation:
 https://docs.aws.amazon.com/cognito/latest/developerguide/user-pool-settings-mfa-totp.html#user-pool-settings-mfa-totp-associate-token
*/
package cognito

import (
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

// Svc represents the Cognito service with essential configuration parameters.
type Svc struct {
	Svc         *cognito.CognitoIdentityProvider
	AppClientId string
	AppSecret   string
}
