package iam

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/iam"

	"gitea/pcp-inariam/inariam/pkgs/log"
)

const (
	ErrIamUserExists     = "IAM User already exists"
	ErrIamRUserNotExists = "IAM User does not exist"
	ErrIamUserEmptyList  = "IAM User list is empty"
)

// GetIamUser checks if an IAM user exists and returns details if found.
func (IamSvc *Svc) GetIamUser(user string) (*iam.User, error) {
	getUserInput := &iam.GetUserInput{
		UserName: aws.String(user),
	}

	userDetails, err := IamSvc.svc.GetUser(getUserInput)
	if err != nil {
		var iamErr awserr.Error
		if errors.As(err, &iamErr) {
			if iamErr.Code() == iam.ErrCodeNoSuchEntityException {
				return nil, fmt.Errorf("CheckIfIamUserExists: %s %w", ErrIamRUserNotExists, err)
			}
		}
		return nil, fmt.Errorf("CheckIfIamUserExists: %w", err)
	}

	return userDetails.User, nil
}

// CreateIamUser creates an IAM user and returns user details.
func (IamSvc *Svc) CreateIamUser(username string) (*iam.User, error) {
	user, err := IamSvc.GetIamUser(username)
	if user != nil {
		return nil, fmt.Errorf("CreateIamUser: %s %w", ErrIamUserExists, err)
	}

	createUserInput := &iam.CreateUserInput{
		UserName: aws.String(username),
	}

	createdUser, err := IamSvc.svc.CreateUser(createUserInput)
	if err != nil {
		return nil, fmt.Errorf("CreateIamUser: %w", err)
	}

	log.Logger.Infof("IAM user '%s' created successfully\n", username)
	return createdUser.User, nil
}

// UpdateIamUser updates an IAM user and returns user details.
func (IamSvc *Svc) UpdateIamUser(oldUsername string, newUsername string) (*iam.User, error) {
	userDetails, err := IamSvc.GetIamUser(oldUsername)

	if err != nil {
		return nil, err
	}
	if userDetails == nil {
		return nil, fmt.Errorf("UpdateIamUser: %s %w", ErrIamRUserNotExists, err)
	}

	updateUserInput := &iam.UpdateUserInput{
		UserName:    aws.String(oldUsername),
		NewUserName: aws.String(newUsername),
	}

	_, err = IamSvc.svc.UpdateUser(updateUserInput)
	if err != nil {
		return nil, fmt.Errorf("UpdateIamUser: %w", err)
	}

	log.Logger.Infof("IAM user '%s' updated successfully\n", oldUsername)

	// Assuming you want to return details of the updated user

	userDetails.UserName = &newUsername

	return userDetails, nil
}

// DeleteIamUser deletes an IAM user.
func (IamSvc *Svc) DeleteIamUser(username string) error {
	userDetails, err := IamSvc.GetIamUser(username)
	if err != nil {
		return fmt.Errorf("DeleteIamUser: %w", err)
	}
	if userDetails == nil {
		return fmt.Errorf("DeleteIamUser: %s %w", ErrIamUserExists, err)
	}

	deleteUserInput := &iam.DeleteUserInput{
		UserName: aws.String(username),
	}

	_, err = IamSvc.svc.DeleteUser(deleteUserInput)
	if err != nil {
		return fmt.Errorf("DeleteIamUser: %w", err)
	}

	log.Logger.Infof("IAM user '%s' deleted successfully\n", username)
	return nil
}

// ListIamUsers lists IAM users and returns details.
func (IamSvc *Svc) ListIamUsers() ([]*iam.User, error) {
	listUsersInput := &iam.ListUsersInput{
		MaxItems: aws.Int64(10),
	}

	listUsersOutput, err := IamSvc.svc.ListUsers(listUsersInput)
	if err != nil {
		return nil, fmt.Errorf("ListIamUsers: %s %w", ErrIamUserEmptyList, err)
	}

	var users []*iam.User

	for _, user := range listUsersOutput.Users {
		// Log or process the user details as needed
		log.Logger.Infof("User Name: %s\n", *user.UserName)
		log.Logger.Infof("User ID: %s\n", *user.UserId)
		log.Logger.Infof("User ARN: %s\n", *user.Arn)
		log.Logger.Infof("Create Date: %s\n", user.CreateDate)

		// Assuming you want to collect the user details
		users = append(users, user)
	}

	return users, nil
}
