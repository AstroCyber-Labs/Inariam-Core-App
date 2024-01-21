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
	ErrIamGroupExists    = "error IAM group already exists"
	ErrIamGroupNotExists = "error IAM group does not exist"
	ErrIamGroupEmptyList = "error IAM group list is empty"
)

// CheckIfGroupExists checks if an IAM group exists and returns its details if found
func (IamSvc *Svc) CheckIfGroupExists(groupName string) (*iam.Group, error) {
	getGroupInput := &iam.GetGroupInput{
		GroupName: aws.String(groupName),
	}

	group, err := IamSvc.svc.GetGroup(getGroupInput)
	if err != nil {
		var iamErr awserr.Error
		if errors.As(err, &iamErr) {
			if iamErr.Code() == iam.ErrCodeNoSuchEntityException {
				return nil, errors.New(ErrIamGroupNotExists)
			}
		}
		return nil, fmt.Errorf("CheckIfGroupExists: %w", err)
	}

	return group.Group, nil
}

// ListIamGroups lists all IAM groups
func (IamSvc *Svc) ListIamGroups() ([]*iam.Group, error) {
	var groups []*iam.Group

	listGroupsInput := &iam.ListGroupsInput{
		MaxItems: aws.Int64(100), // Adjust max items according to your requirements
	}

	err := IamSvc.svc.ListGroupsPages(listGroupsInput,
		func(page *iam.ListGroupsOutput, lastPage bool) bool {
			groups = append(groups, page.Groups...)
			return !lastPage
		},
	)
	if err != nil {
		return nil, fmt.Errorf("ListIamGroups: %w\n", err)
	}

	if len(groups) == 0 {
		return nil, errors.New(ErrIamGroupEmptyList)
	}

	for _, group := range groups {
		log.Logger.Infof("Group Name: %s\n", *group.GroupName)
		log.Logger.Infof("Group ID: %s\n", *group.GroupId)
		log.Logger.Infof("Group ARN: %s\n", *group.Arn)
		log.Logger.Infof("Create Date: %s\n", group.CreateDate)
	}

	return groups, nil
}

// UpdateIamGroup updates the name of an existing IAM group and returns the updated group details
func (IamSvc *Svc) UpdateIamGroup(oldGroupName string, newGroupName string) (*iam.Group, error) {
	group, err := IamSvc.CheckIfGroupExists(oldGroupName)
	if err != nil {
		return nil, fmt.Errorf("UpdateIamGroup: %w", err)
	}

	if group == nil {
		return nil, fmt.Errorf("UpdateIamGroup: %s %w", ErrIamGroupNotExists, err)
	}

	updateGroupInput := &iam.UpdateGroupInput{
		GroupName:    aws.String(oldGroupName),
		NewGroupName: aws.String(newGroupName),
	}

	_, err = IamSvc.svc.UpdateGroup(updateGroupInput)
	if err != nil {
		return nil, fmt.Errorf("UpdateIamGroup: %w", err)
	}

	// Fetch and return updated group details
	getGroupInput := &iam.GetGroupInput{
		GroupName: aws.String(newGroupName),
	}

	getGroupOutput, err := IamSvc.svc.GetGroup(getGroupInput)
	if err != nil {
		return nil, fmt.Errorf("UpdateIamGroup: %w", err)
	}

	log.Logger.Infof("IAM group '%s' updated successfully\n", oldGroupName)
	return getGroupOutput.Group, nil
}

// CreateIamGroup creates a new IAM group and returns the group details
func (IamSvc *Svc) CreateIamGroup(groupName string) (*iam.Group, error) {
	group, err := IamSvc.CheckIfGroupExists(groupName)
	if group != nil {
		return nil, fmt.Errorf("CreateIamGroup: %s %w", ErrIamGroupExists, err)
	}

	createGroupInput := &iam.CreateGroupInput{
		GroupName: aws.String(groupName),
	}

	_, err = IamSvc.svc.CreateGroup(createGroupInput)
	if err != nil {
		return nil, fmt.Errorf("CreateIamGroup: %w ", err)
	}

	// Fetch and return group details after creation
	getGroupInput := &iam.GetGroupInput{
		GroupName: aws.String(groupName),
	}

	getGroupOutput, err := IamSvc.svc.GetGroup(getGroupInput)
	if err != nil {
		return nil, fmt.Errorf("CreateIamGroup: %w ", err)
	}

	log.Logger.Infof("IAM group '%s' created successfully\n", groupName)
	return getGroupOutput.Group, nil
}

// DeleteIamGroup deletes an existing IAM group
func (IamSvc *Svc) DeleteIamGroup(groupName string) error {
	group, err := IamSvc.CheckIfGroupExists(groupName)
	if err != nil {
		return fmt.Errorf("DeleteIamGroup: %w", err)
	}
	if group == nil {
		return fmt.Errorf("DeleteIamGroup: %s %w", ErrIamGroupNotExists, err)
	}

	deleteGroupInput := &iam.DeleteGroupInput{
		GroupName: aws.String(groupName),
	}

	_, err = IamSvc.svc.DeleteGroup(deleteGroupInput)
	if err != nil {
		return fmt.Errorf("DeleteIamGroup: %w", err)
	}

	log.Logger.Infof("IAM group '%s' deleted successfully. \n", groupName)
	return nil
}
