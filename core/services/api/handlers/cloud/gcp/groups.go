package gcp

import (
	"net/http"

	"github.com/labstack/echo/v4"

	req "gitea/pcp-inariam/inariam/core/services/api/requests/gcp/iam"
	"gitea/pcp-inariam/inariam/core/services/api/responses"
	resp "gitea/pcp-inariam/inariam/core/services/api/responses/gcp"
	"gitea/pcp-inariam/inariam/pkgs/log"
)

// ListGroups
// @Summary Get a list of IAM GCP groups
// @Description Retrieves a list of IAM GCP groups using the Google Cloud IAM service.
// @ID gcp-list-groups
// @Produce json
// @Success 200 {array} []resp.GroupDetails
// @Router /gcp/iam/groups [get]
func (handler *Handler) ListGroups(ctx echo.Context) error {
	gcpSession, err := handler.openAdminIamSession()
	if err != nil {
		return responses.ErrorResponse(
			ctx,
			http.StatusInternalServerError,
			resp.HttpErrOpenedSession,
		)
	}

	groups, err := gcpSession.IamAdminGCPService.ListGroups()
	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	var groupsList []resp.GroupDetails
	for _, group := range groups {
		groupsList = append(groupsList, resp.GroupDetails{
			Name:        group.Name,
			Id:          group.Id,
			Description: group.Description,
		})
	}

	return responses.Response(ctx, http.StatusOK, groupsList)
}

// GetGroup
// @Summary Get gcp-Group by Name
// @Description Get gcp-IAM-group details by group name
// @ID get-gcp-group-by-name
// @Param groupName path string true "Group Name"
// @Produce json
// @Success 200 {object} resp.GroupDetails
// @Router /gcp/iam/groups/{groupName} [get]
func (handler *Handler) GetGroup(ctx echo.Context) error {

	groupReq := req.ActionGroupRequest{}

	if err := ctx.Bind(&groupReq); err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, responses.HttpErrBadRequest)
	}

	err := groupReq.Validate()
	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, responses.HttpErrBadRequest)
	}

	gcpSession, err := handler.openAdminIamSession()
	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, "Failed to open IAM session")
	}

	group, err := gcpSession.IamAdminGCPService.GetGroup(groupReq.Name)

	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	if group == nil {
		return responses.ErrorResponse(ctx, http.StatusNotFound, "Group not found")
	}

	return responses.Response(ctx, http.StatusOK, resp.GroupDetails{
		Name:        group.Name,
		Id:          group.Id,
		Description: group.Description,
	})
}

// CreateGroup
// @Summary Create GCP Group
// @Description Create a new GCP IAM group.
// @ID create-group
// @Accept json
// @Produce json
// @Param body body iam.CreateGroupRequest true "Group details"
// @Success 200 {object} resp.GroupDetails
// @Router /gcp/iam/groups [post]
func (handler *Handler) CreateGroup(c echo.Context) error {
	createGroupRequest := req.CreateGroupRequest{}

	if err := c.Bind(&createGroupRequest); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, responses.HttpErrBadRequest)
	}

	err := createGroupRequest.Validate()

	if err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusBadRequest,
			"Make sure you have provided a valid group name",
		)
	}

	gcpSession, err := handler.openAdminIamSession()
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, responses.HttpErrOpenedSession)
	}

	newGroup, err := gcpSession.IamAdminGCPService.CreateGroup(
		createGroupRequest.Name,
		createGroupRequest.Description,
	)

	if err != nil {
		log.Logger.Errorln(err.Error())
		return responses.ErrorResponse(c, http.StatusBadRequest, "Error creating group")
	}

	return responses.Response(c, http.StatusOK, resp.GroupDetails{
		Name:        newGroup.Name,
		Id:          newGroup.Id,
		Description: newGroup.Description,
	})
}

// UpdateGroup
// @Summary Update Group
// @Description Update an existing IAM group
// @ID update-group
// @Accept json
// @Produce json
// @Param body body iam.UpdateGroupRequest true "Updated group details"
// @Success 200 {object} resp.GroupDetails
// @Router /gcp/iam/groups [put]
func (handler *Handler) UpdateGroup(c echo.Context) error {
	updateGroupRequest := req.UpdateGroupRequest{}

	if err := c.Bind(&updateGroupRequest); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, responses.HttpErrBadRequest)
	}

	err := updateGroupRequest.Validate()

	if err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusBadRequest,
			"Make sure you have provided all fields",
		)
	}

	gcpSession, err := handler.openAdminIamSession()
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, responses.HttpErrOpenedSession)
	}

	updatedGroup, err := gcpSession.IamAdminGCPService.UpdateGroupDescription(
		updateGroupRequest.Name,
		updateGroupRequest.Description,
	)

	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return responses.Response(c, http.StatusOK, resp.GroupDetails{
		Name:        updatedGroup.Name,
		Id:          updatedGroup.Id,
		Description: updateGroupRequest.Description,
	})
}

// DeleteGroup
// @Summary Delete Group
// @Description Delete an IAM group by group name
// @ID delete-group
// @Param groupName path string true "Group Name"
// @Produce json
// @Success 200 {string} string "Group deleted successfully"
// @Router /gcp/iam/groups/{groupName} [delete]
func (handler *Handler) DeleteGroup(ctx echo.Context) error {

	groupReq := req.ActionGroupRequest{}

	if err := ctx.Bind(&groupReq); err != nil {
		return responses.MessageResponse(ctx, http.StatusBadRequest, responses.HttpErrBadRequest)
	}

	err := groupReq.Validate()

	gcpSession, err := handler.openAdminIamSession()
	if err != nil {
		return responses.MessageResponse(ctx, http.StatusBadRequest, responses.HttpErrOpenedSession)
	}

	err = gcpSession.IamAdminGCPService.DeleteGroup(groupReq.Name)
	if err != nil {
		return responses.MessageResponse(ctx, http.StatusBadRequest, err.Error())
	}

	return responses.MessageResponse(ctx, http.StatusOK, "Group deleted successfully.")
}
