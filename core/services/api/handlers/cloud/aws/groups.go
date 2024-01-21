package aws

import (
	"net/http"

	"github.com/labstack/echo/v4"

	req "gitea/pcp-inariam/inariam/core/services/api/requests/aws/iam"
	"gitea/pcp-inariam/inariam/core/services/api/responses"
	resp "gitea/pcp-inariam/inariam/core/services/api/responses/aws/iam"
)

// ListGroups
// @Summary AWS List Groups
// @Description Get a list of AWS IAM groups
// @ID aws-list-groups
// @Produce plain
// @Success 200 {string} string "Group names separated by newline"
// @Router /aws/iam/groups [get]
func (awsHandler *Handler) ListGroups(ctx echo.Context) error {
	awsSession, err := awsHandler.RetrieveAwsIamSession()
	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusInternalServerError, responses.HttpErrOpenedSession)
	}

	groups, err := awsSession.IamSvc.ListIamGroups()
	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, responses.HttpErrServerFailed)
	}

	var groupsList []resp.Group
	for _, group := range groups {
		groupsList = append(groupsList, resp.Group{
			Arn:        *group.Arn,
			CreateDate: *group.CreateDate,
			Id:         *group.GroupId,
			Name:       *group.GroupName,
			Path:       *group.Path,
		})
	}

	return responses.Response(ctx, http.StatusOK, groupsList)
}

// GetGroup
// @Summary Get AWS Group by Name
// @Description Get AWS IAM group details by group name
// @ID get-aws-group-by-name
// @Param groupName path string true "Group Name"
// @Produce json
// @Success 200 {object} resp.Group
// @Router /aws/iam/groups/{groupName} [get]
func (awsHandler *Handler) GetGroup(c echo.Context) error {
	groupName := c.Param("id")

	awsSession, err := awsHandler.RetrieveAwsIamSession()
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, responses.HttpErrOpenedSession)
	}

	groupDetails, err := awsSession.IamSvc.CheckIfGroupExists(groupName)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return responses.Response(c, http.StatusOK, resp.Group{
		Name:       *groupDetails.GroupName,
		Id:         *groupDetails.GroupId,
		CreateDate: *groupDetails.CreateDate,
		Arn:        *groupDetails.Arn,
		Path:       groupName,
	},
	)
}

// CreateGroup
// @Summary Create AWS Group
// @Description Create a new AWS IAM group
// @ID create-aws-group
// @Accept json
// @Produce json
// @Param body body iam.CreateGroupRequest true "Group details"
// @Success 200 {object} resp.Group
// @Router /aws/iam/groups [post]
func (awsHandler *Handler) CreateGroup(ctx echo.Context) error {
	var createGrpReq = new(req.CreateGroupRequest)
	if err := ctx.Bind(&createGrpReq); err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, responses.HttpErrBadRequest)
	}

	err := createGrpReq.Validate()
	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	awsSession, err := awsHandler.RetrieveAwsIamSession()
	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, responses.HttpErrServerFailed)
	}

	groupDetails, err := awsSession.IamSvc.CreateIamGroup(createGrpReq.GroupName)
	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	return responses.Response(ctx, http.StatusOK, resp.Group{
		Name:       *groupDetails.GroupName,
		Id:         *groupDetails.GroupId,
		CreateDate: *groupDetails.CreateDate,
		Arn:        *groupDetails.Arn,
		Path:       *groupDetails.Path,
	},
	)
}

// UpdateGroup
// @Summary Update AWS Group
// @Description Update an existing AWS IAM group's name
// @ID update-aws-group
// @Accept json
// @Produce json
// @Param groupName path string true "Group Name"
// @Param  body body req.UpdateGroupRequest true "New group details"
// @Success 200 {object} resp.Group
// @Router /aws/iam/groups/{groupName} [put]
func (awsHandler *Handler) UpdateGroup(ctx echo.Context) error {
	groupName := ctx.Param("id")

	var req = new(req.UpdateGroupRequest)
	if err := ctx.Bind(&req); err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, responses.HttpErrBadRequest)
	}

	awsSession, err := awsHandler.RetrieveAwsIamSession()
	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, responses.HttpErrServerFailed)
	}

	groupDetails, err := awsSession.IamSvc.UpdateIamGroup(groupName, req.NewGroupName)
	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	return responses.Response(ctx, http.StatusOK, resp.Group{
		Name:       *groupDetails.GroupName,
		Id:         *groupDetails.GroupId,
		CreateDate: *groupDetails.CreateDate,
		Arn:        *groupDetails.Arn,
		Path:       *groupDetails.Path,
	},
	)
}

// DeleteGroup
// @Summary Delete AWS Group
// @Description Delete an AWS IAM group by group name
// @ID delete-aws-group
// @Param groupName path string true "Group Name"
// @Produce json
// @Success 200 {string} string "AWS Group deleted successfully"
// @Router /aws/iam/groups/{groupName} [delete]
func (awsHandler *Handler) DeleteGroup(ctx echo.Context) error {
	groupName := ctx.Param("id")

	awsSession, err := awsHandler.RetrieveAwsIamSession()
	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusInternalServerError, responses.HttpErrServerFailed)
	}

	err = awsSession.IamSvc.DeleteIamGroup(groupName)

	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, "Failed to delete group")
	}

	return responses.Response(ctx, http.StatusOK, true)
}
