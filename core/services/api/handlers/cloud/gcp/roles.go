package gcp

import (
	"gitea/pcp-inariam/inariam/pkgs/log"
	"net/http"

	"github.com/labstack/echo/v4"

	req "gitea/pcp-inariam/inariam/core/services/api/requests/gcp/iam"
	"gitea/pcp-inariam/inariam/core/services/api/responses"
	resp "gitea/pcp-inariam/inariam/core/services/api/responses/gcp"
	"gitea/pcp-inariam/inariam/pkgs/cloud/gcp/iam"
)

// CreateIamRole
// @Summary Create an IAM Role
// @Description Create an IAM role for the project
// @ID gcp-create-role
// @Accept json
// @Produce json
// @Param body body req.CreateRoleRequest true "IAM role details"
// @Success 200 {object} resp.RoleResponse
// @Router /gcp/roles [post]
func (handler *Handler) CreateIamRole(ctx echo.Context) error {
	// Map incoming request to CreateRoleRequest
	createRoleRequest := &req.CreateRoleRequest{}
	if err := ctx.Bind(createRoleRequest); err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, responses.HttpErrBadRequest)
	}

	log.Logger.Infoln(createRoleRequest)

	// Validate the request
	if err := createRoleRequest.Validate(); err != nil {
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

	newRole := iam.NewRole{
		Id:          createRoleRequest.Id,
		Name:        createRoleRequest.Name,
		Title:       createRoleRequest.Title,
		Description: createRoleRequest.Description,
		Permissions: createRoleRequest.Permissions,
	}

	createdRole, err := iamSession.IamGCPService.CreateIamRole(newRole)
	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())

	}

	return responses.Response(ctx, http.StatusCreated, resp.RoleResponse{
		Name:        createdRole.Name,
		Title:       createdRole.Title,
		Description: createdRole.Description,
		Stage:       createdRole.Stage,
		Permissions: createdRole.IncludedPermissions,
	})

}

// UpdateIamRole
// @Summary Update an IAM Role
// @Description Update an IAM role for the project
// @ID gcp-update-role
// @Accept json
// @Produce json
// @Param body body req.UpdateRoleRequest true "IAM role details"
// @Success 200 {object}  resp.RoleResponse
// @Router /gcp/roles/{name} [put]
func (handler *Handler) UpdateIamRole(ctx echo.Context) error {
	updateRoleRequest := &req.UpdateRoleRequest{}
	if err := ctx.Bind(updateRoleRequest); err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, responses.HttpErrBadRequest)
	}

	// Validate the request
	if err := updateRoleRequest.Validate(); err != nil {
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

	log.Logger.Infoln(updateRoleRequest)

	updatedRole, err := iamSession.IamGCPService.UpdateIamRole(
		updateRoleRequest.Id,
		updateRoleRequest.Title,
		updateRoleRequest.Description,
		updateRoleRequest.Permissions,
	)

	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	return responses.Response(ctx, http.StatusOK, resp.RoleResponse{
		Name:        updatedRole.Name,
		Title:       updatedRole.Title,
		Description: updatedRole.Description,
		Stage:       updatedRole.Stage,
		Permissions: updatedRole.IncludedPermissions,
	})
}

// DeleteIamRole
// @Summary Delete an IAM Role
// @Description Delete an IAM role for the project
// @ID gcp-delete-role
// @Accept json
// @Produce json
// @Param body body req.ActionRoleRequest true "IAM role name"
// @Success 200 {bool}  true
// @Router /gcp/roles [delete]
func (handler *Handler) DeleteIamRole(ctx echo.Context) error {

	actionRoleRequest := &req.ActionRoleRequest{}
	if err := ctx.Bind(actionRoleRequest); err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, responses.HttpErrBadRequest)
	}

	// Validate the request
	if err := actionRoleRequest.Validate(); err != nil {
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

	err = iamSession.IamGCPService.DeleteIamRole(actionRoleRequest.Name)

	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	return responses.MessageResponse(ctx, http.StatusAccepted, "Role deleted successfully")
}

// ListRoles
// @Summary List IAM Roles
// @Description List all IAM roles the project
// @ID gcp-list-role
// @Accept json
// @Produce json
// @Success 200 {array}  []resp.RoleResponse
// @Router /gcp/roles [get]
func (handler *Handler) ListRoles(ctx echo.Context) error {

	iamSession, err := handler.openIamSession()

	if err != nil {
		return responses.ErrorResponse(
			ctx,
			http.StatusInternalServerError,
			responses.HttpErrOpenedSession,
		)
	}

	roles, err := iamSession.IamGCPService.ListIamRoles()

	rolesResp := make([]*resp.RoleResponse, len(roles))
	for i, role := range roles {
		rolesResp[i] = &resp.RoleResponse{
			Name:        role.Name,
			Title:       role.Title,
			Description: role.Description,
			Stage:       role.Stage,
			Permissions: role.IncludedPermissions,
		}
	}

	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	return responses.Response(ctx, http.StatusOK, rolesResp)
}

// GetRole
// @Summary Get an IAM Role
// @Description Get an IAM role of the project
// @ID gcp-one-role
// @Accept json
// @Produce json
// @Success 200 {object}  resp.RoleResponse
// @Router /gcp/roles/{name} [get]
func (handler *Handler) GetRole(ctx echo.Context) error {

	actionRoleRequest := &req.ActionRoleRequest{}
	if err := ctx.Bind(actionRoleRequest); err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, responses.HttpErrBadRequest)
	}

	// Validate the request
	if err := actionRoleRequest.Validate(); err != nil {
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

	role, err := iamSession.IamGCPService.GetIamRole(actionRoleRequest.Name)

	log.Logger.Infoln(role)

	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	if role == nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, resp.HttpErrRoleNotFound)
	}

	return responses.Response(ctx, http.StatusAccepted,
		resp.RoleResponse{
			Name:        role.Name,
			Title:       role.Title,
			Description: role.Description,
			Stage:       role.Stage,
			Permissions: role.IncludedPermissions,
		})
}
