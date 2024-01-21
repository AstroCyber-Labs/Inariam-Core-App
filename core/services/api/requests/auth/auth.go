// Package auth provides structures and functionality related to user authentication processes.
package auth

import "github.com/go-playground/validator/v10"

// LoginProcessRequest represents a request to initiate the login process.
type LoginProcessRequest struct {
	Email    string `json:"email" validate:"required,email" example:"john.doe@example.com"`
	Password string `json:"password" validate:"required" example:"11111111"`
}

// Validate validates the LoginProcessRequest structure using the go-playground/validator library.
func (loginRequest *LoginProcessRequest) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	return validate.Struct(loginRequest)
}

// VerifyUserRequest represents a request to verify a user's identity.
type VerifyUserRequest struct {
	Code  string `json:"verification_code" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

// Validate validates the VerifyUserRequest structure using the go-playground/validator library.
func (verifyUserRequest *VerifyUserRequest) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	return validate.Struct(verifyUserRequest)
}

// GenerateMFARequest represents a request to generate a Multi-Factor Authentication (MFA) code.
type GenerateMFARequest struct {
	AuthSession string `json:"auth_session"`
}

// Validate validates the GenerateMFARequest structure using the go-playground/validator library.
func (genMFARequest *GenerateMFARequest) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	return validate.Struct(genMFARequest)
}

// GetMFADeviceCode represents a request to get the Multi-Factor Authentication (MFA) device code.
type GetMFADeviceCode struct {
	AuthSession string `json:"auth_session"`
	Email       string `json:"email" validate:"email"`
}

// AttachMfaDeviceRequest represents a request to attach an MFA device.
type AttachMfaDeviceRequest struct {
	AuthSession string `json:"auth_session"`
	Email       string `json:"email" validate:"email"`
	TOTPCode    string `json:"totp_code" validate:"min:6"`
}

// Validate validates the AttachMfaDeviceRequest structure using the go-playground/validator library.
func (attachMfaDeviceReq *AttachMfaDeviceRequest) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	return validate.Struct(attachMfaDeviceReq)
}

// ResendConfirmationCodeRequest represents a request to resend a confirmation code.
type ResendConfirmationCodeRequest struct {
	Email string `json:"email"`
}

// Validate validates the ResendConfirmationCodeRequest structure using the go-playground/validator library.
func (resendConfirmationEmail *ResendConfirmationCodeRequest) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	return validate.Struct(resendConfirmationEmail)
}

// CompleteMfaSignInRequest represents a request to complete the sign-in process with Multi-Factor Authentication (MFA).
type CompleteMfaSignInRequest struct {
	AuthSession string `json:"auth_session"`
	TOTPCode    string `json:"totp_code" validate:"min:6"`
	Email       string `json:"email" validate:"email"`
}

// Validate validates the CompleteMfaSignInRequest structure using the go-playground/validator library.
func (completeMfaSignInReq *CompleteMfaSignInRequest) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	return validate.Struct(completeMfaSignInReq)
}
