// Package iam provides functionality for managing Identity and Access Management (IAM) roles in Google Cloud.
package iam

import (
	"errors"
	"fmt"

	"google.golang.org/api/googleapi"
	"google.golang.org/api/iam/v1"

	"gitea/pcp-inariam/inariam/pkgs/log"
)

// Error messages for IAM role-related failures.
const (
	ErrorRoleAlreadyExists          = "role already exists"
	ErrorRoleDoesNotExist           = "role does not exist"
	ErrorFailedToCheckRoleExistence = "failed to check role existence"
	ErrorFailedToCreateRole         = "failed to create role"
	ErrorFailedToGetRoleInfo        = "failed to get role information"
	ErrorFailedToDeleteRole         = "failed to delete role"
	ErrorFailedToUpdateRole         = "failed to update role"
)

// NewRole represents the information needed to create a new IAM role.
type NewRole struct {
	Id          string
	Name        string
	Title       string
	Description string
	Stage       *string
	Permissions []string
}

// GetIamRole retrieves information about an IAM role by its name.
func (iamService *IamSvc) GetIamRole(roleName string) (*iam.Role, error) {
	roleName = "projects/" + iamService.projectId + "/roles/" + roleName
	log.Logger.Infoln(roleName)
	role, err := iamService.svc.Projects.Roles.Get(roleName).Do()

	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, nil
		}
		return nil, err
	}

	log.Logger.Infof("Role Name: %s\n", role.Name)
	log.Logger.Infof("Role Title: %s\n", role.Title)
	log.Logger.Infof("Role Description: %s\n", role.Description)
	log.Logger.Infof("Role Stage: %s\n", role.Stage)

	for _, permission := range role.IncludedPermissions {
		log.Logger.Infof("Got permission: %v\n", permission)
	}
	return role, nil
}

// CreateIamRole creates a new IAM role.
func (iamService *IamSvc) CreateIamRole(newRole NewRole) (*iam.Role, error) {
	foundRole, err := iamService.GetIamRole(newRole.Name)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", ErrorFailedToCheckRoleExistence, err)
	}

	if foundRole != nil {
		return nil, fmt.Errorf("%s : %s", ErrorRoleAlreadyExists, newRole.Name)
	}

	roleStageValue := "ALPHA"
	if newRole.Stage != nil {
		roleStageValue = *newRole.Stage
	}

	request := &iam.CreateRoleRequest{
		RoleId: newRole.Id,
		Role: &iam.Role{
			Title:               newRole.Title,
			Description:         newRole.Description,
			Stage:               roleStageValue,
			IncludedPermissions: newRole.Permissions,
		},
	}

	role, err := iamService.svc.Projects.Roles.Create("projects/"+iamService.projectId, request).
		Do()
	if err != nil {
		return nil, fmt.Errorf("%s : %w", ErrorFailedToCreateRole, err)
	}

	log.Logger.Infof("Custom role created: %s\n", newRole.Name)

	return role, nil
}

// UpdateIamRole updates an existing IAM role.
func (iamService *IamSvc) UpdateIamRole(
	roleId string,
	updatedTitle string,
	updatedDescription string,
	updatedPermissions []string,
) (*iam.Role, error) {

	role, err := iamService.GetIamRole(roleId)
	if err != nil {
		return nil, err
	}

	roleId = "projects/" + iamService.projectId + "/roles/" + roleId

	if role == nil {
		return nil, fmt.Errorf("UpdateIamRole: Role not found")
	}

	role.Title = updatedTitle
	role.Description = updatedDescription
	role.IncludedPermissions = updatedPermissions

	newRole, err := iamService.svc.Projects.Roles.Patch(roleId, role).Do()
	if err != nil {

		return nil, fmt.Errorf("%s : %v", ErrorFailedToUpdateRole, err)
	}

	log.Logger.Infof("Role with title '%s' updated successfully.\n", role.Name)
	return newRole, nil
}

// DeleteIamRole deletes an IAM role by its name.
func (iamService *IamSvc) DeleteIamRole(roleName string) error {
	rlName := "projects/" + iamService.projectId + "/roles/" + roleName

	_, err := iamService.svc.Projects.Roles.Delete(rlName).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return fmt.Errorf("%s : %s", ErrorRoleDoesNotExist, roleName)
		}
		return fmt.Errorf("%s : %v", ErrorFailedToDeleteRole, err)
	}

	log.Logger.Infof("Role with title '%s' deleted successfully.\n", roleName)
	return nil
}

// ListIamRoles retrieves a list of all IAM roles in the project.
func (iamService *IamSvc) ListIamRoles() ([]*iam.Role, error) {
	response, err := iamService.svc.Projects.Roles.List("projects/" + iamService.projectId).Do()
	if err != nil {
		return nil, fmt.Errorf("%s : %v", ErrorFailedToGetRoleInfo, err)
	}
	for _, role := range response.Roles {
		log.Logger.Infof("Listing role: %v\n", role.Name)
	}
	return response.Roles, nil
}
