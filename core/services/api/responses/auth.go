// Package responses provides structures for common HTTP responses.
package responses

// HTTP error messages.
const (
	HttpErrMissingLoginField = "Email and password are required"
	HttpErrInvalidCode       = "Invalid code"
)

// LoginResponse represents a response for a login operation.
type LoginResponse struct {
	AuthSession   string `json:"auth_session"`
	ChallengeName string `json:"challenge_name"`
}

// GenerateMFAResponse represents a response for generating Multi-Factor Authentication (MFA) codes.
type GenerateMFAResponse struct {
	AuthSession   string `json:"auth_session"`
	ImageQrCode   string `json:"qrcode_image"`
	PlainCode     string `json:"plain_code"`
	ChallengeName string `json:"challenge_name"`
}

// CompleteSignInResponse represents a response for completing a sign-in operation.
type CompleteSignInResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
