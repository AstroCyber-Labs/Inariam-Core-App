package gcp

import (
	"net/http"

	"github.com/labstack/echo/v4"
	iam "google.golang.org/api/iam/v1"

	req "gitea/pcp-inariam/inariam/core/services/api/requests/gcp/iam"
	"gitea/pcp-inariam/inariam/core/services/api/responses"
	resp "gitea/pcp-inariam/inariam/core/services/api/responses/gcp"
	"gitea/pcp-inariam/inariam/pkgs/log"
)

// CreateServiceAccount
// @Summary Create Service Account
// @Description Create a new IAM service account
// @ID create-service-account
// @Accept json
// @Produce json
// @Param body body req.CreateServiceAccountRequest true "Service account details"
// @Success 200 {object} iam.ServiceAccount
// @Router /service-accounts [post]
func (handler *Handler) CreateServiceAccount(ctx echo.Context) error {

	createServiceAccountReq := &req.CreateServiceAccountRequest{}
	if err := ctx.Bind(createServiceAccountReq); err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, responses.HttpErrBadRequest)
	}
	log.Logger.Infoln(createServiceAccountReq)
	// Validate the request
	if err := createServiceAccountReq.Validate(); err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	iamSession, err := handler.openIamSession()
	if err != nil {
		return responses.ErrorResponse(
			ctx,
			http.StatusInternalServerError,
			responses.HttpErrOpenedSession,
		)
	}

	var createdAccount *iam.ServiceAccount
	createdAccount, err = iamSession.IamGCPService.CreateIamServiceAccount(
		createServiceAccountReq.DisplayName,
		createServiceAccountReq.Name,
		createServiceAccountReq.Description)

	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	// Return the response
	return responses.Response(ctx, http.StatusOK, createdAccount)
}

// ListServiceAccounts
// @Summary List Service Accounts
// @Description Get a list of IAM service accounts
// @ID list-service-accounts
// @Produce json
// @Param project_id path string true "Project ID"
// @Success 200 {array} resp.ServiceAccountDetails
// @Router /service-accounts/{project_id} [get]
func (handler *Handler) ListServiceAccounts(ctx echo.Context) error {

	iamSession, err := handler.openIamSession()
	if err != nil {
		return responses.ErrorResponse(
			ctx,
			http.StatusInternalServerError,
			responses.HttpErrOpenedSession,
		)
	}

	serviceAccs, err := iamSession.IamGCPService.ListIamServiceAccounts()

	// Map the response to the desired format
	var details []resp.ServiceAccountDetails
	for _, account := range serviceAccs {
		detail := resp.ServiceAccountDetails{
			Name:        account.Name,
			Description: account.Description,
			Etag:        account.Etag,
			Disabled:    account.Disabled,
			DisplayName: account.DisplayName,
			UniqueId:    account.UniqueId,
			ProjectId:   account.ProjectId,
		}
		details = append(details, detail)
	}

	// Return the response
	return responses.Response(ctx, http.StatusOK, details)
}

// EnableServiceAccount
// @Summary List Service Accounts
// @Description Get a list of IAM service accounts
// @ID list-service-accounts
// @Produce json
// @Param project_id path string true "Project ID"
// @Success 200 {array} resp.ServiceAccountDetails
// @Router /service-accounts/{project_id} [get]
func (handler *Handler) EnableServiceAccount(ctx echo.Context) error {
	serviceAccountAction := &req.ActionOneServiceAccountRequest{}
	if err := ctx.Bind(serviceAccountAction); err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, responses.HttpErrBadRequest)
	}
	log.Logger.Infoln(serviceAccountAction)

	iamSession, err := handler.openIamSession()
	if err != nil {
		return responses.ErrorResponse(
			ctx,
			http.StatusInternalServerError,
			responses.HttpErrOpenedSession,
		)
	}

	err = iamSession.IamGCPService.EnableIamServiceAccount(serviceAccountAction.Name)

	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	return responses.MessageResponse(ctx, http.StatusOK, "service account enabled")

}

// DisableServiceAccount @Summary List Service Accounts
// @Description Get a list of IAM service accounts
// @ID list-service-accounts
// @Produce json
// @Param project_id path string true "Project ID"
// @Success 200 {array} resp.ServiceAccountDetails
// @Router /service-accounts/{project_id} [get]
func (handler *Handler) DisableServiceAccount(ctx echo.Context) error {
	serviceAccountAction := &req.ActionOneServiceAccountRequest{}
	if err := ctx.Bind(serviceAccountAction); err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, responses.HttpErrBadRequest)
	}

	// Validate the request
	if err := serviceAccountAction.Validate(); err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	iamSession, err := handler.openIamSession()
	if err != nil {
		return responses.ErrorResponse(
			ctx,
			http.StatusInternalServerError,
			responses.HttpErrOpenedSession,
		)
	}

	err = iamSession.IamGCPService.DisableIamServiceAccount(serviceAccountAction.Name)

	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	return responses.MessageResponse(ctx, http.StatusOK, "service account disabled")

}

// GetServiceAccount @Summary Get Service Accounts
// @Description Get a list of IAM service accounts
// @ID list-service-accounts
// @Produce json
// @Param project_id path string true "Project ID"
// @Success 200 {array} resp.ServiceAccountDetails
// @Router /service-accounts/{project_id} [get]
func (handler *Handler) GetServiceAccount(ctx echo.Context) error {
	serviceAccountAction := &req.ActionOneServiceAccountRequest{}
	if err := ctx.Bind(serviceAccountAction); err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, responses.HttpErrBadRequest)
	}

	if err := serviceAccountAction.Validate(); err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	iamSession, err := handler.openIamSession()
	if err != nil {
		return responses.ErrorResponse(
			ctx,
			http.StatusInternalServerError,
			responses.HttpErrOpenedSession,
		)
	}

	serviceAcc, err := iamSession.IamGCPService.GetServiceAccount(serviceAccountAction.Name)

	// Return the response
	return responses.Response(ctx, http.StatusOK, serviceAcc)
}

// DeleteServiceAccount
// @Summary Delete Service Account
// @Description Get a list of IAM service accounts
// @ID list-service-accounts
// @Produce json
// @Param project_id path string true "Project ID"
// @Success 200 {array} resp.ServiceAccountDetails
// @Router /service-accounts/{project_id} [delete]
func (handler *Handler) DeleteServiceAccount(ctx echo.Context) error {
	serviceAccountAction := &req.ActionOneServiceAccountRequest{}
	if err := ctx.Bind(serviceAccountAction); err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, responses.HttpErrBadRequest)
	}

	// Validate the request
	if err := serviceAccountAction.Validate(); err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	iamSession, err := handler.openIamSession()
	if err != nil {
		return responses.ErrorResponse(
			ctx,
			http.StatusInternalServerError,
			responses.HttpErrOpenedSession,
		)
	}

	err = iamSession.IamGCPService.DeleteIamServiceAccount(serviceAccountAction.Name)

	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	return responses.MessageResponse(ctx, http.StatusOK, "service account deleted successfully")

}
