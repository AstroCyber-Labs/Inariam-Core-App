package cognito

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/skip2/go-qrcode"

	"gitea/pcp-inariam/inariam/pkgs/log"
)

const (
	ErrAcquiringMFACode           = "error acquiring MFA Code failed"
	ErrGeneratingQrCode           = "error generating QRcode string"
	ErrGeneratingQrCodePNG        = "error generating QRcode png"
	ErrActivatingMFA              = "error activating MFA"
	ErrInvalidMFACode             = "error verifying MFA code"
	ErrInvalidVerificationCode    = "error verifying verification code"
	ErrResendingSignUpConfirmCode = "error resending sign-up confirmation code"
	ErrHookingDeviceTOTP          = "error hooking device TOTP"
	ErrCompletingMFAFlow          = "error completing MFA flow"
	ErrAcquiringSecretCode        = "error acquiring secret code"
	ErrStartingSignInProcess      = "error starting sign in process"
	ErrConfirmingSignUp           = "error confirming sign-up"
	ErrSignUp                     = "error signing up"
	ErrAppSecretEmpty             = "COGNITO_CLIENT_SECRET not set"
)

type AuthClient interface {
	SignUp(email, password string) (string, error)
	ConfirmSignUp(email, code string) (string, error)
	SimpleSignIn(email, password string) (string, error)
	ResendConfirmationCode(email string) error
	StartSignInProcess(email, password string) (error, string)
	CompleteMFAAuthFlow(username, code, session string) *MFAAuthSuccessResult
}

// SignUp a plain signup, needs both email and password. Even when passing email, you pass it under the username field
func (cognitoSvc *Svc) SignUp(email, password string) (string, error) {
	user := &cognito.SignUpInput{
		Username:   aws.String(email),
		Password:   aws.String(password),
		ClientId:   aws.String(cognitoSvc.AppClientId),
		SecretHash: aws.String(computeSecretHash(cognitoSvc.AppSecret, email, cognitoSvc.AppClientId)),
	}

	result, err := cognitoSvc.Svc.SignUp(user)
	if err != nil {
		return "", fmt.Errorf("SignUp: %s, %w", ErrSignUp, err)
	}
	return result.String(), nil
}

// ConfirmSignUp confirming sign-up with the code sent to the email ( or SMS depending on the console configuration )
// returns error
func (cognitoSvc *Svc) ConfirmSignUp(email, code string) error {
	if cognitoSvc.AppSecret == "" {
		return fmt.Errorf("Cognito.ConfirmSignUp: %s", ErrAppSecretEmpty)
	}
	confirmSignUpInput := &cognito.ConfirmSignUpInput{
		Username:         aws.String(email),
		ConfirmationCode: aws.String(code),
		ClientId:         aws.String(cognitoSvc.AppClientId),
		SecretHash:       aws.String(computeSecretHash(cognitoSvc.AppSecret, email, cognitoSvc.AppClientId)),
	}
	_, err := cognitoSvc.Svc.ConfirmSignUp(confirmSignUpInput)
	if err != nil {
		return fmt.Errorf("Cognito.ConfirmSignUp: %s, %w", ErrConfirmingSignUp, err)
	}

	return nil
}

// ResendConfirmSignUp resend the confirmation code to the user email
func (cognitoSvc *Svc) ResendConfirmSignUp(email string) error {
	secretHash := computeSecretHash(cognitoSvc.AppSecret, email, cognitoSvc.AppClientId)

	input := cognito.ResendConfirmationCodeInput{
		ClientId:   &cognitoSvc.AppClientId,
		Username:   &email,
		SecretHash: &secretHash,
	}

	_, err := cognitoSvc.Svc.ResendConfirmationCode(&input)

	if err != nil {
		return fmt.Errorf("Cognito.ResendConfirmSignUp: %s %w", ErrResendingSignUpConfirmCode, err)
	}

	return nil
}

// SimpleSignIn login with AWS Cognito returns an AccessToken
// Use for users that doesn't have MFA enabled.
func (cognitoSvc *Svc) SimpleSignIn(email string, password string) (string, error) {

	initialAuthInput := &cognito.InitiateAuthInput{
		AuthFlow: aws.String("USER_PASSWORD_AUTH"),
		AuthParameters: aws.StringMap(map[string]string{
			"USERNAME": email,
			"PASSWORD": password,
			"SECRET_HASH": computeSecretHash(
				cognitoSvc.AppSecret,
				email,
				cognitoSvc.AppClientId,
			),
		}),
		ClientId: aws.String(cognitoSvc.AppClientId),
	}
	result, err := cognitoSvc.Svc.InitiateAuth(initialAuthInput)

	if err != nil {
		return "", fmt.Errorf("Cognito.SimpleSignIn: %s, %w", ErrInvalidMFACode, err)
	}

	return *result.Session, nil
}

// StartSignInProcessResult
type StartSignInProcessResult struct {
	ChallengeName string
	SessionKey    string
}

// StartSignInProcess it starts the sign in process, it returns a session that will be used to complete the MFA process.
// returns Session,error ( can return accessToken however that is not needed for now )
func (cognitoSvc *Svc) StartSignInProcess(email, password string) (*StartSignInProcessResult, error) {
	authParamMap := make(map[string]string)

	authParamMap["USERNAME"] = email
	authParamMap["PASSWORD"] = password
	authParamMap["SECRET_HASH"] = computeSecretHash(
		cognitoSvc.AppSecret,
		email,
		cognitoSvc.AppClientId,
	)

	input := cognito.InitiateAuthInput{
		AuthFlow:       aws.String(cognito.AuthFlowTypeUserPasswordAuth),
		ClientId:       &cognitoSvc.AppClientId,
		AuthParameters: aws.StringMap(authParamMap),
	}

	res, err := cognitoSvc.Svc.InitiateAuth(&input)

	if err != nil {
		// print("Error logging in", err.Error())
		return nil, fmt.Errorf("StartSignInProcess: %s, %w", ErrStartingSignInProcess, err)
	}

	return &StartSignInProcessResult{
		ChallengeName: *res.ChallengeName,
		SessionKey:    *res.Session,
	}, nil

}

type MFASetupResult struct {
	Session string
	QrCode  string
	Code    string
}

// GenerateMFAActivationCode it gets the secret code of the user and generates a QR code that will be sent back to the user, so he scans it with the authenticator app.
// returns Session,QRCode,error
func (cognitoSvc *Svc) GenerateMFAActivationCode(session, email string) (*MFASetupResult, error) {
	associateSoftwareTokenInput := cognito.AssociateSoftwareTokenInput{
		AccessToken: aws.String(session),
		Session:     aws.String(session),
	}
	output, err := cognitoSvc.Svc.AssociateSoftwareToken(&associateSoftwareTokenInput)
	if err != nil {
		return nil, fmt.Errorf("Cognito.GenerateMFAActivationCode: %s, %w", ErrAcquiringSecretCode, err)
	}

	qrCodeString := fmt.Sprintf("otpauth://totp/Inariam:%s?secret=%s", email, *output.SecretCode)

	qrCode, err := generateQRCode(qrCodeString, 256)

	if err != nil {
		return nil, fmt.Errorf("Cognito.GenerateMFAActivationCode: %w", err)
	}

	return &MFASetupResult{
		Session: *output.Session,
		QrCode:  qrCode,
		Code:    *output.SecretCode,
	}, nil
}

// ConfirmMFAActivation it verifies a code from the user to a generated one and activates MFA.
// returns error
func (cognitoSvc *Svc) ConfirmMFAActivation(token, email, code string) error {
	output, err := cognitoSvc.Svc.VerifySoftwareToken(&cognito.VerifySoftwareTokenInput{
		AccessToken:        &token,
		FriendlyDeviceName: &email,
		UserCode:           &code,
	})
	if err != nil {
		return fmt.Errorf("Cognito.ConfirmMFAActivation: %s, %w", ErrActivatingMFA, err)
	}

	if *output.Status == "SUCCESS" {
		return nil
	}

	return fmt.Errorf("Cognito.ConfirmMFAActivation: %s", ErrActivatingMFA)
}

// ConfirmTOTPDevice it completes the MFA setup process.
func (cognitoSvc *Svc) ConfirmTOTPDevice(username, code, session string) (string, error) {
	challengeResponses := make(map[string]string)

	challengeResponses["USERNAME"] = username
	challengeResponses["SECRET_HASH"] = computeSecretHash(cognitoSvc.AppSecret, username, cognitoSvc.AppClientId)
	challengeResponses["SOFTWARE_MFA_CODE"] = code
	challengeResponses["SESSION"] = session

	input := &cognito.RespondToAuthChallengeInput{
		ChallengeName:      aws.String("MFA_SETUP"),
		Session:            &session,
		ClientId:           &cognitoSvc.AppClientId,
		ChallengeResponses: aws.StringMap(challengeResponses),
	}

	res, err := cognitoSvc.Svc.RespondToAuthChallenge(input)
	if err != nil {
		log.Logger.Errorln("error occurred ", err)
		return "", fmt.Errorf("Cognito.CompleteMFASetup: %s, %w", ErrGeneratingQrCodePNG, err)
	}

	log.Logger.Infoln("Completing MFA SETUP %v\n", res)

	return *res.Session, nil
}

type MFAAuthSuccessResult struct {
	RefreshToken string
	AccessToken  string
	IdToken      string
}

func (cognitoSvc *Svc) CompleteMFAAuthFlow(username, code, session string) (*MFAAuthSuccessResult, error) {
	challengeResponses := make(map[string]string)

	challengeResponses["USERNAME"] = username
	challengeResponses["SECRET_HASH"] = computeSecretHash(cognitoSvc.AppSecret, username, cognitoSvc.AppClientId)
	challengeResponses["SOFTWARE_TOKEN_MFA_CODE"] = code

	input := &cognito.RespondToAuthChallengeInput{
		ChallengeName:      aws.String("SOFTWARE_TOKEN_MFA"),
		Session:            &session,
		ClientId:           &cognitoSvc.AppClientId,
		ChallengeResponses: aws.StringMap(challengeResponses),
	}

	res, err := cognitoSvc.Svc.RespondToAuthChallenge(input)

	if err != nil {
		return nil, fmt.Errorf("Cognito.CompleteMFAAuthFlow: %s, %w", ErrCompletingMFAFlow, err)
	}

	return &MFAAuthSuccessResult{
		RefreshToken: *res.AuthenticationResult.RefreshToken,
		AccessToken:  *res.AuthenticationResult.AccessToken,
		IdToken:      *res.AuthenticationResult.IdToken,
	}, nil
}

// generateQRCode generate QR Code base64 png to pass to an <img> tag.
func generateQRCode(input string, size int) (string, error) {

	qrCode, err := qrcode.New(input, qrcode.Medium)

	if err != nil {
		return "", fmt.Errorf("Cognito.generateQRCode: %s, %w", ErrGeneratingQrCode, err)
	}

	pngBytes, err := qrCode.PNG(size)

	if err != nil {
		return "", fmt.Errorf("Cognito.generateQRCode: %s, %w", ErrGeneratingQrCodePNG, err)
	}

	encodedQr := base64.StdEncoding.EncodeToString(pngBytes)

	return encodedQr, nil
}

// computeSecretHash It creates a secret hash that will help identify the user from the users pool side.
func computeSecretHash(clientSecret, username, clientId string) string {
	mac := hmac.New(sha256.New, []byte(clientSecret))
	mac.Write([]byte(username + clientId))

	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
