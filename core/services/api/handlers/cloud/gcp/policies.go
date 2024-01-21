package gcp

import (
	"net/http"

	req "gitea/pcp-inariam/inariam/core/services/api/requests/gcp/iam"
	"gitea/pcp-inariam/inariam/core/services/api/responses"
	"github.com/labstack/echo/v4"
	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
)

// SetPolicy
// @Summary Set IAM Policy
// @Description Set IAM policy for a project
// @ID set-iam-policy
// @Accept json
// @Produce json
// @Param projectID path string true "Project ID"
// @Param body body req.SetIamPolicyRequest true "IAM policy details"
// @Success 200 {string} string "IAM policy set successfully"
// @Router /gcp/iam/policy [put]
func (handler *Handler) SetPolicy(ctx echo.Context) error {

	setIamPolicyRequest := req.SetIamPolicyRequest{}
	if err := ctx.Bind(&setIamPolicyRequest); err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, responses.HttpErrBadRequest)
	}

	err := setIamPolicyRequest.Validate()
	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	crmService, err := handler.openCrmSession()
	if err != nil {
		return responses.ErrorResponse(
			ctx,
			http.StatusInternalServerError,
			responses.HttpErrOpenedSession,
		)
	}

	policy := cloudresourcemanager.Policy{
		Version:  setIamPolicyRequest.Policy.Version,
		Etag:     setIamPolicyRequest.Policy.Etag,
		Bindings: make([]*cloudresourcemanager.Binding, len(setIamPolicyRequest.Policy.Bindings)),
	}

	// Map Bindings
	for i, binding := range setIamPolicyRequest.Policy.Bindings {
		policy.Bindings[i] = &cloudresourcemanager.Binding{
			Role:    binding.Role,
			Members: binding.Members,
		}
	}

	err = crmService.CrmGCPService.SetPolicy(&policy)
	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	return responses.Response(ctx, http.StatusOK, "IAM policy set successfully")
}

// GetPolicy
// @Summary Get IAM Policy
// @Description Get IAM policy for a project
// @ID gcp-one-policy
// @Param projectID path string true "Project ID"
// @Produce json
// @Success 200 {object} cloudresourcemanager.Policy
// @Router /gcp/iam/policy [get]
func (handler *Handler) GetPolicy(c echo.Context) error {

	crmService, err := handler.openCrmSession()

	if err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusInternalServerError,
			responses.HttpErrOpenedSession,
		)
	}

	policy, err := crmService.CrmGCPService.GetIamPolicy()
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return responses.Response(c, http.StatusOK, policy)
}

// DeletePolicy
// @Summary Delete IAM Policy
// @Description Delete IAM policy for a project
// @ID delete-iam-policy
// @Param projectID path string true "Project ID"
// @Produce json
// @Success 200 {string} string "IAM policy deleted successfully"
// @Router /gcp/iam/policy [delete]
func (handler *Handler) DeletePolicy(c echo.Context) error {

	crmService, err := handler.openCrmSession()
	if err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusInternalServerError,
			responses.HttpErrOpenedSession,
		)
	}

	err = crmService.CrmGCPService.DeletePolicy()
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return responses.Response(c, http.StatusOK, true)
}
