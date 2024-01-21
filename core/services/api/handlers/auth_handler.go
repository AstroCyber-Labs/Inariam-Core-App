package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"gitea/pcp-inariam/inariam/core/services/api"
	requests "gitea/pcp-inariam/inariam/core/services/api/requests/auth"
	"gitea/pcp-inariam/inariam/core/services/api/responses"
	inaAws "gitea/pcp-inariam/inariam/pkgs/cloud/aws"
	"gitea/pcp-inariam/inariam/pkgs/log"
)

type AuthHandler struct {
	api *api.API
}

func NewAuthHandler(api *api.API) *AuthHandler {
	return &AuthHandler{api}
}

// @Summary		Authenticate a user
// @Description	Perform user login
// @ID				user-login
// @Tags			User Actions
// @Accept			json
// @Produce		json
// @Param		params body	requests.LoginProcessRequest true	"User's credentials"
// @Success		200		{object}	responses.LoginResponse
// @Failure		401		{object}	responses.Error
// @Router			/auth/login [post]
func (authHandler *AuthHandler) StartSignInProcess(ctx echo.Context) error {

	loginReq := requests.LoginProcessRequest{}

	if err := ctx.Bind(&loginReq); err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, responses.HttpErrMissingLoginField)
	}

	err := loginReq.Validate()

	if err != nil {
		fmt.Println("Error")
		return responses.ErrorResponse(ctx, http.StatusBadRequest, responses.HttpErrMissingLoginField)
	}

	awsSession, err := authHandler.retrieveAwsSession()
	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusInternalServerError, responses.HttpErrServerFailed)
	}

	res, err := awsSession.CognitoSvc.StartSignInProcess(loginReq.Email, loginReq.Password)
	if err != nil {
		log.Logger.Infoln(err.Error())
		return responses.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	switch res.ChallengeName {

	case "MFA_SETUP":
		{
			res, err := awsSession.CognitoSvc.GenerateMFAActivationCode(res.SessionKey, loginReq.Email)

			if err != nil {
				return responses.ErrorResponse(ctx, http.StatusInternalServerError, responses.HttpErrInvalidCode)
			}

			return responses.Response(ctx, http.StatusOK, responses.GenerateMFAResponse{
				AuthSession:   res.Session,
				ChallengeName: "MFA_SETUP_RESPONSE",
				ImageQrCode:   res.QrCode,
				PlainCode:     res.Code,
			})
		}
	case "SOFTWARE_TOKEN_MFA":
		{
			return responses.Response(ctx, http.StatusOK, responses.LoginResponse{
				AuthSession:   res.SessionKey,
				ChallengeName: res.ChallengeName,
			})
		}

	}

	return responses.ErrorResponse(ctx, http.StatusInternalServerError, responses.HttpErrServerFailed)
}

// @Summary		Verify a user email
// @Description	Verifying a user email with the code sent to him in the email
// @ID				verify-user
// @Tags			User Actions
// @Accept				json
// @Produce			 	json
// @Param		params body	requests.VerifyUserRequest	true	"User's credentials"
// @Success		200		{boolean}	bool
// @Failure		401		{object}	responses.Error
// @Router			/auth/verify-user [post]
func (authHandler *AuthHandler) VerifyUser(ctx echo.Context) error {
	verifyUserReq := requests.VerifyUserRequest{}

	if err := ctx.Bind(&verifyUserReq); err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, responses.HttpErrInvalidCode)
	}

	awsSession, err := authHandler.retrieveAwsSession()
	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusInternalServerError, responses.HttpErrServerFailed)
	}

	err = awsSession.CognitoSvc.ConfirmSignUp(verifyUserReq.Email, verifyUserReq.Code)

	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, responses.HttpErrInvalidCode)
	}

	return responses.MessageResponse(ctx, http.StatusOK, "Your account is verified proceed to login.")
}

// @Summary		Activate MFA For user
// @Description	Perform user login
// @ID				activate-code
// @Tags			User Actions
// @Accept		json
// @Produce		json
// @Param		params body	requests.GetMFADeviceCode	true "User's credentials"
// @Success		200		{string} false "String response here"
// @Failure		401		{object}	responses.Error
// @Router			/auth/activate-mfa [post]
func (authHandler *AuthHandler) GetMFADeviceCode(ctx echo.Context) error {
	attachMfaDeviceReq := requests.GetMFADeviceCode{}

	if err := ctx.Bind(&attachMfaDeviceReq); err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, responses.HttpErrBadRequest)
	}

	awsSession, err := authHandler.retrieveAwsSession()
	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusInternalServerError, responses.HttpErrServerFailed)
	}

	res, err := awsSession.CognitoSvc.GenerateMFAActivationCode(
		attachMfaDeviceReq.AuthSession,
		attachMfaDeviceReq.Email,
	)

	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusInternalServerError, responses.HttpErrInvalidCode)
	}

	return responses.Response(ctx, http.StatusOK, responses.GenerateMFAResponse{
		AuthSession: res.Session,
		PlainCode:   res.Code,
		ImageQrCode: res.QrCode,
	})
}

// @Summary		Confirm MFA Code
// @Description	Perform user login
// @ID				confirm-code
// @Tags			User Actions
// @Accept		json
// @Produce		json
// @Param		params body	requests.AttachMfaDeviceRequest	true "User's credentials"
// @Success		200		{string} false "String response here"
// @Failure		401		{object}	responses.Error
// @Router			/auth/confirm-mfa [post]
func (authHandler *AuthHandler) ConfirmMFACode(ctx echo.Context) error {
	attachMfaDeviceReq := requests.AttachMfaDeviceRequest{}

	if err := ctx.Bind(&attachMfaDeviceReq); err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, responses.HttpErrBadRequest)
	}

	awsSession, err := authHandler.retrieveAwsSession()
	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusInternalServerError, responses.HttpErrServerFailed)
	}

	err = awsSession.CognitoSvc.ConfirmMFAActivation(
		attachMfaDeviceReq.AuthSession,
		attachMfaDeviceReq.Email,
		attachMfaDeviceReq.TOTPCode,
	)

	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusInternalServerError, responses.HttpErrInvalidCode)
	}

	return responses.MessageResponse(ctx, http.StatusOK, "Successfully added device, please proceed to sign-in.")
}

// @Summary Resender user's email
// @Description Resend confirmation code to user
// @ID resend-confirmation-code
// @Tags User Actions
// @Accept json
// @Produce json
// @Param params body requests.ResendConfirmationCodeRequest true "User's credentials"
// @Success 200 {object} responses.Data
// @Failure 401 {object} responses.Error
// @Router /auth/confirmation-resend [post]
func (authHandler *AuthHandler) ResendConfirmatioNEmail(ctx echo.Context) error {

	resendConfirmationEmail := requests.ResendConfirmationCodeRequest{}

	if err := ctx.Bind(&resendConfirmationEmail); err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, responses.HttpErrBadRequest)
	}

	err := resendConfirmationEmail.Validate()

	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, "Please provide email")
	}

	awsSession, err := authHandler.retrieveAwsSession()
	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusInternalServerError, responses.HttpErrServerFailed)
	}

	err = awsSession.CognitoSvc.ResendConfirmSignUp(resendConfirmationEmail.Email)

	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, "Error resending email")
	}

	return responses.MessageResponse(ctx, http.StatusOK, "A new activation code was sent to your email")
}

// @Summary Complete signin flow
// @Description Complete signin flow and receive accessToken and
// @ID complete-signin-flow
// @Tags User Actions
// @Accept json
// @Produce json
// @Param params body requests.CompleteMfaSignInRequest true "User's credentials"
// @Success 200 {object} responses.CompleteSignInResponse
// @Failure 401 {object} responses.Error
// @Router /auth/signin [post]
func (authHandler *AuthHandler) CompleteSignIn(ctx echo.Context) error {
	completeMFASignInReq := requests.CompleteMfaSignInRequest{}

	if err := ctx.Bind(&completeMFASignInReq); err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, responses.HttpErrBadRequest)
	}

	awsSession, err := authHandler.retrieveAwsSession()
	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusInternalServerError, responses.HttpErrServerFailed)
	}

	res, err := awsSession.CognitoSvc.CompleteMFAAuthFlow(completeMFASignInReq.Email, completeMFASignInReq.AuthSession, completeMFASignInReq.TOTPCode)
	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusUnauthorized, responses.HttpErrInvalidCode)
	}

	return responses.Response(ctx, http.StatusOK, responses.CompleteSignInResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	})
}

// retrieveAwsSession retrieve AwsSession with the config passed to the handler
func (authHandler *AuthHandler) retrieveAwsSession() (*inaAws.Session, error) {
	creds := inaAws.Credentials{
		AccessKeyID:     authHandler.api.Config.AWS.AccessKeyID,
		SecretAccessKey: authHandler.api.Config.AWS.SecretAccessKey,
		Region:          authHandler.api.Config.AWS.Region,
	}
	awsSession, err := inaAws.OpenSession(&creds)

	if err != nil {
		return nil, err
	}

	awsSession.CreateCognitoSvc(authHandler.api.Config.CognitoConfig.ClientId, authHandler.api.Config.CognitoConfig.AppSecret)

	return awsSession, nil
}
