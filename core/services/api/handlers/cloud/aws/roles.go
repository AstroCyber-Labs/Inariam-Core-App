package aws

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"

	req "gitea/pcp-inariam/inariam/core/services/api/requests/aws/iam"
	"gitea/pcp-inariam/inariam/core/services/api/responses"
	resp "gitea/pcp-inariam/inariam/core/services/api/responses/aws/iam"
	"gitea/pcp-inariam/inariam/pkgs/log"
)

// ListRoles @Summary List Roles
// @Description Get a list of IAM roles
// @ID aws-roles-list
// @Produce json
// @Success 200 {array} resp.RoleDetailResponse
// @Router /aws/roles [get]
func (awsHandler *Handler) ListRoles(c echo.Context) error {
	openedSession, err := awsHandler.RetrieveAwsIamSession()

	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, responses.HttpErrOpenedSession)
	}

	roles, err := openedSession.IamSvc.ListIamRoles()
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, responses.HttpErrOpenedSession)
	}

	var rolesList []resp.RoleDetailResponse

	for _, role := range roles {
		detail := resp.RoleDetailResponse{
			RoleName: *role.RoleName,
			RoleID:   *role.RoleId,
			RoleArn:  *role.Arn,
		}
		rolesList = append(rolesList, detail)
	}

	return responses.Response(c, http.StatusOK, rolesList)
}

// GetRole @Summary Get Role by Name
// @Description Get IAM role details by role name
// @ID aws-role-by-name
// @Param roleName path string true "Role Name"
// @Produce json
// @Success 200 {object} resp.RoleDetailResponse
// @Router /aws/roles/{roleName} [get]
func (awsHandler *Handler) GetRole(c echo.Context) error {
	roleName := c.Param("id")

	awsSession, err := awsHandler.RetrieveAwsIamSession()
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, responses.HttpErrOpenedSession)
	}

	roleDetails, err := awsSession.IamSvc.GetIamRole(roleName)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	if roleDetails == nil {
		return responses.ErrorResponse(c, http.StatusNotFound, resp.HttpErrRoleNotFound)
	}

	return responses.Response(c, http.StatusOK, resp.RoleDetailResponse{
		RoleName: *roleDetails.RoleName,
		RoleID:   *roleDetails.RoleId,
		RoleArn:  *roleDetails.Arn,
	},
	)
}

// CreateRole @Summary Create Role
// @Description Create a new IAM role
// @ID aws-role-create
// @Accept json
// @Produce json
// @Param roleName body string true "Role Name"
// @Param trustPolicy body requests.CreateRoleRequest.TrustPolicy true "Trust Policy"
// @Success 200 {object} resp.RoleDetailResponse
// @Router /aws/roles [post]
func (awsHandler *Handler) CreateRole(c echo.Context) error {

	createRoleRequest := req.CreateRoleRequest{}

	if err := c.Bind(&createRoleRequest); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, responses.HttpErrBadRequest)
	}

	err := createRoleRequest.Validate()
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, resp.HttpErrMissingRoleName)
	}

	trustPolicyJSON, err := json.Marshal(createRoleRequest.TrustPolicy)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, resp.ErrorConvertingTrustPolicyToJSON)
	}

	openedSession, err := awsHandler.RetrieveAwsIamSession()
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, responses.HttpErrOpenedSession)
	}

	createdRole, err := openedSession.IamSvc.CreateIAMRole(createRoleRequest.RoleName, string(trustPolicyJSON))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return responses.Response(c, http.StatusOK, resp.RoleDetailResponse{
		RoleName: createRoleRequest.RoleName,
		RoleArn:  *createdRole.Arn,
		RoleID:   *createdRole.RoleId,
	})
}

// UpdateRole @Summary Update Role
// @Description Update an existing IAM role's trust policy
// @ID aws-role-update
// @Accept json
// @Produce json
// @Param roleName body string true "Role Name"
// @Param trustPolicy body requests.UpdateRoleRequest.TrustPolicy true "Updated Trust Policy"
// @Success 200 {object} resp.RoleDetailResponse
// @Router /aws/roles [put]
func (awsHandler *Handler) UpdateRole(c echo.Context) error {
	updateRoleRequest := req.UpdateRoleRequest{}

	if err := c.Bind(&updateRoleRequest); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, responses.HttpErrBadRequest)
	}

	log.Logger.Infoln(updateRoleRequest)
	err := updateRoleRequest.Validate()
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, resp.HttpErrMissingRoleName)
	}

	trustPolicyJSON, err := json.Marshal(updateRoleRequest.TrustPolicy)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Error converting trust policy to JSON")
	}

	openedSession, err := awsHandler.RetrieveAwsIamSession()
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Failed to retrieve IAM session")
	}

	roleDetails, err := openedSession.IamSvc.ModifyIAMRoleTrustPolicy(updateRoleRequest.Name, string(trustPolicyJSON))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return responses.Response(c, http.StatusOK, resp.RoleDetailResponse{
		RoleName: *roleDetails.RoleName,
		RoleID:   *roleDetails.RoleId,
		RoleArn:  *roleDetails.Arn,
	})
}

// DeleteRole @Summary Delete Role
// @Description Delete an IAM role by role name
// @ID aws-role-delete
// @Param roleName path string true "Role Name"
// @Produce json
// @Success 200 {string} string "Role deleted successfully"
// @Router /aws/roles/{roleName} [delete]
func (awsHandler *Handler) DeleteRole(c echo.Context) error {
	roleName := c.Param("id")
	awsSession, err := awsHandler.RetrieveAwsIamSession()
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, responses.HttpErrOpenedSession)
	}

	err = awsSession.IamSvc.DeleteIAMRole(roleName)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return responses.Response(c, http.StatusOK, true)
}
