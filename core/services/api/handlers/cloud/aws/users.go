package aws

import (
	"net/http"

	"github.com/labstack/echo/v4"

	req "gitea/pcp-inariam/inariam/core/services/api/requests/aws/iam"
	"gitea/pcp-inariam/inariam/core/services/api/responses"
	resp "gitea/pcp-inariam/inariam/core/services/api/responses/aws/iam"
)

// @Summary List Users
// @Description Get a list of IAM users
// @ID list-users
// @Produce json
// @Success 200 {array} resp.UserDetailResponse
// @Router /users [get]
func (awsHandler *Handler) ListUsers(ctx echo.Context) error {
	openedSession, err := awsHandler.RetrieveAwsIamSession()
	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve IAM session")
	}

	users, err := openedSession.IamSvc.ListIamUsers()
	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, resp.HttpErrMissingUserName)
	}

	var usersDetails []resp.UserDetailResponse
	for _, user := range users {
		detail := resp.UserDetailResponse{
			Username: *user.UserName,
			UserID:   *user.UserId,
			UserArn:  *user.Arn,
		}
		usersDetails = append(usersDetails, detail)
	}

	return responses.Response(ctx, http.StatusOK, usersDetails)
}

// @Summary Get User by ID
// @Description Get IAM user details by username
// @ID get-user-by-id
// @Param id path string true "Username"
// @Produce json
// @Success 200 {object} resp.UserDetailResponse
// @Router /users/{id} [get]
func (awsHandler *Handler) GetUser(ctx echo.Context) error {
	username := ctx.Param("id")

	awsSession, err := awsHandler.RetrieveAwsIamSession()
	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, resp.HttpErrMissingUserName)
	}

	userDetails, err := awsSession.IamSvc.GetIamUser(username)
	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, responses.HttpErrOpenedSession)
	}
	if userDetails == nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, resp.HttpErrMissingUserName)
	}

	return responses.Response(ctx, http.StatusOK, resp.UserDetailResponse{
		Username: *userDetails.UserName,
		UserID:   *userDetails.UserId,
		UserArn:  *userDetails.Arn,
	},
	)
}

// @Summary Create User
// @Description Create a new IAM user
// @ID create-user
// @Accept json
// @Produce json
// @Param username body string true "Username"
// @Success 201 {object} resp.UserDetailResponse
// @Router /users [post]
func (awsHandler *Handler) CreateUser(ctx echo.Context) error {
	createUserRequest := req.CreateUserRequest{}

	if err := ctx.Bind(&createUserRequest); err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	err := createUserRequest.Validate()

	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	awsSess, err := awsHandler.RetrieveAwsIamSession()
	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, responses.HttpErrOpenedSession)
	}

	createdUser, err := awsSess.IamSvc.CreateIamUser(createUserRequest.Username)
	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	return responses.Response(ctx, http.StatusCreated, resp.UserDetailResponse{
		Username: createUserRequest.Username,
		UserID:   *createdUser.UserId,
		UserArn:  *createdUser.Arn,
	},
	)
}

// @Summary Delete User
// @Description Delete an IAM user by username
// @ID delete-user
// @Param id path string true "Username"
// @Produce json
// @Success 200 {string} string "IAM user deleted successfully"
// @Router /users/{id} [delete]
func (awsHandler *Handler) DeleteUser(ctx echo.Context) error {
	deleteUserRequest := req.DeleteUserRequest{}
	if err := ctx.Bind(&deleteUserRequest); err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, resp.HttpErrMissingUserName)
	}

	err := deleteUserRequest.Validate()

	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, resp.HttpErrMissingUserName)
	}

	awsSess, err := awsHandler.RetrieveAwsIamSession()
	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, responses.HttpErrOpenedSession)
	}

	err = awsSess.IamSvc.DeleteIamUser(deleteUserRequest.Username)
	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, responses.HttpErrOpenedSession)
	}

	return responses.Response(ctx, http.StatusOK, true)
}

// @Summary Update User
// @Description Update an existing IAM user's username
// @ID update-user
// @Accept json
// @Produce json
// @Param username body string true "Current Username"
// @Param newUsername body string true "New Username"
// @Success 200 {string} string "IAM user updated successfully"
// @Router /users [put]
func (awsHandler *Handler) UpdateUser(ctx echo.Context) error {
	updateUserRequest := req.UpdateUserRequest{}
	if err := ctx.Bind(&updateUserRequest); err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, resp.HttpErrMissingUserName)
	}

	err := updateUserRequest.Validate()

	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, resp.HttpErrMissingUserName)
	}

	awsSession, err := awsHandler.RetrieveAwsIamSession()
	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, responses.HttpErrOpenedSession)
	}

	updatedUser, err := awsSession.IamSvc.UpdateIamUser(updateUserRequest.Username, updateUserRequest.NewUsername)

	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, "Error creating user")
	}

	// TODO: FIX THIS, FOR UPDATES WE RETURN THE NEW MODIFIED OBJECT
	return responses.Response(ctx, http.StatusOK, resp.UserDetailResponse{
		Username: *updatedUser.UserName,
		UserID:   *updatedUser.UserId,
		UserArn:  *updatedUser.Arn,
	})
}
