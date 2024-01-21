package aws

import (
	"net/http"

	"github.com/labstack/echo/v4"

	req "gitea/pcp-inariam/inariam/core/services/api/requests/aws/iam"
	"gitea/pcp-inariam/inariam/core/services/api/responses"
	resp "gitea/pcp-inariam/inariam/core/services/api/responses/aws/iam"
	"gitea/pcp-inariam/inariam/pkgs/cloud/aws/iam"
)

// ListPolicies @Summary List Policies
// @Description Get a list of IAM policies
// @ID list-policies
// @Produce json
// @Success 200 {array} resp.PolicyDetailResponse
// @Router /aws/policies [get]
func (awsHandler *Handler) ListPolicies(c echo.Context) error {
	awsSession, err := awsHandler.RetrieveAwsIamSession()
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve IAM session")
	}

	policies, err := awsSession.IamSvc.ListIamPolicies()
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, responses.HttpErrBadRequest)
	}

	// var policiesList []resp.PolicyDetailResponse
	// for _, policy := range policies {
	// 	policiesList = append(policiesList, resp.PolicyDetailResponse{
	// 		PolicyName: *policy.PolicyName,
	// 		PolicyID:   *policy.PolicyId,
	// 	})
	// }

	return responses.Response(c, http.StatusOK, policies)
}

// GetPolicy @Summary Get Policy by Name
// @Description Get IAM policy details by policy name
// @ID aws-policy-by-name
// @Param policyName path string true "Policy Name"
// @Produce json
// @Success 200 {object} resp.PolicyDetailResponse
// @Router /aws/policies/{policyName} [get]
func (awsHandler *Handler) GetPolicy(c echo.Context) error {
	policyName := c.Param("arn")
	if policyName == "" {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Policy ARN is required")
	}

	openedSession, err := awsHandler.RetrieveAwsIamSession()
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Failed to retrieve IAM session")
	}

	policyDetails, err := openedSession.IamSvc.GetIamPolicy(policyName)

	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	if policyDetails == nil {
		return responses.ErrorResponse(c, http.StatusNotFound, "Policy not found")
	}

	return responses.Response(c, http.StatusOK, policyDetails)
}

// CreatePolicy @Summary Create Policy
// @Description Create a new IAM policy
// @ID aws-policy-create
// @Accept json
// @Produce json
// @Param body body req.CreatePolicyRequest true "Policy details"
// @Success 200 {object} resp.PolicyDetailResponse
// @Router /aws/policies [post]
func (awsHandler *Handler) CreatePolicy(c echo.Context) error {
	createPolicyRequest := req.CreatePolicyRequest{}

	if err := c.Bind(&createPolicyRequest); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, responses.HttpErrBadRequest)
	}

	err := createPolicyRequest.Validate()

	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, resp.ErrorMissingPolicyARN)
	}

	// TODO TO CHANGE ALL THE CODE BELOW IT'S BLACK MAGIC

	awsSession, err := awsHandler.RetrieveAwsIamSession()
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, responses.HttpErrOpenedSession)
	}

	statementsEntries := []iam.StatementEntry{}

	for _, statementEntry := range createPolicyRequest.PolicyDocument.Statements {
		statementsEntries = append(statementsEntries, iam.StatementEntry{
			Effect:   statementEntry.Effect,
			Action:   statementEntry.Action,
			Resource: statementEntry.Resource,
		})
	}

	awsPolicy, err := awsSession.IamSvc.CreateIamPolicy(createPolicyRequest.PolicyName, createPolicyRequest.Description, iam.PolicyDocument{
		Version:   createPolicyRequest.PolicyDocument.Version,
		Statement: statementsEntries,
	})

	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return responses.Response(c, http.StatusOK, awsPolicy)
}

// UpdatePolicy @Summary Update Policy
// @Description Update an existing IAM policy
// @ID aws-policy-update
// @Accept json
// @Produce json
// @Param body body req.UpdatePolicyRequest true "Updated policy details"
// @Success 200 {object} resp.PolicyDetailResponse
// @Router /aws/policies [put]
func (awsHandler *Handler) UpdatePolicy(c echo.Context) error {

	updatePolicyRequest := req.UpdatePolicyRequest{}

	if err := c.Bind(&updatePolicyRequest); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, responses.HttpErrBadRequest)
	}

	err := updatePolicyRequest.Validate()

	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, resp.ErrorMissingPolicyARN)
	}

	awsSession, err := awsHandler.RetrieveAwsIamSession()
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, responses.HttpErrOpenedSession)
	}

	statementsEntries := []iam.StatementEntry{}

	for _, statementEntry := range updatePolicyRequest.PolicyDocument.Statements {

		statementsEntries = append(statementsEntries, iam.StatementEntry{
			Effect:   statementEntry.Effect,
			Action:   statementEntry.Action,
			Resource: statementEntry.Resource,
		})
	}

	updatedPolicy, err := awsSession.IamSvc.UpdateIamPolicy(updatePolicyRequest.PolicyARN, iam.PolicyDocument{
		Version:   updatePolicyRequest.PolicyDocument.Version,
		Statement: statementsEntries,
	})

	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return responses.Response(c, http.StatusOK, updatedPolicy)
}

// DeletePolicy @Summary Delete Policy
// @Description Delete an IAM policy by policy name
// @ID aws-delete-policy
// @Param policyName path string true "Policy Name"
// @Produce json
// @Success 200 {string} string "Policy deleted successfully"
// @Router /aws/policies/{policyName} [delete]
func (awsHandler *Handler) DeletePolicy(c echo.Context) error {
	// TODO
	policyARN := c.Param("arn")

	awsSession, err := awsHandler.RetrieveAwsIamSession()
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, responses.HttpErrOpenedSession)
	}

	err = awsSession.IamSvc.DeleteIamPolicy(policyARN)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return responses.Response(c, http.StatusOK, "Policy deleted successfully.")
}
