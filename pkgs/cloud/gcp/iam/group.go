// Package iam provides functionality for managing groups in Google Cloud Identity and Access Management (IAM).
package iam

import (
	"fmt"

	admin "google.golang.org/api/admin/directory/v1"

	"gitea/pcp-inariam/inariam/pkgs/log"
)

// ErrCreatingGroup is the error message for group creation failure.
const (
	ErrCreatingGroup = "error creating a group"
)

// CreateGroup creates a new group with the specified name and description.
func (adminSvc *AdminSvc) CreateGroup(groupName, description string) (*admin.Group, error) {
	newGroup := &admin.Group{
		Description: description,
		Name:        groupName,
	}

	resp, err := adminSvc.svc.Groups.Insert(newGroup).Do()
	if err != nil {
		return nil, fmt.Errorf("CreateGroup: %s : %w", ErrCreatingGroup, err)
	}

	return resp, nil
}

// ListGroups retrieves a list of all groups.
func (adminSvc *AdminSvc) ListGroups() ([]*admin.Group, error) {
	groups, err := adminSvc.svc.Groups.List().Do()
	if err != nil {
		log.Logger.Errorf("ListGroups: failed to list groups: %v", err)
		return nil, err
	}

	for _, group := range groups.Groups {
		fmt.Printf("Name: %s, Id: %s, Description: %s\n",
			group.Name, group.Id, group.Description)
	}
	return groups.Groups, nil
}

// GetGroup retrieves information about a specific group.
func (adminSvc *AdminSvc) GetGroup(groupName string) (*admin.Group, error) {
	group, err := adminSvc.svc.Groups.Get(groupName).Do()

	if err != nil {
		log.Logger.Errorf("GetGroup: failed to get group: %v", err)
		return nil, err
	}

	return group, nil
}

// UpdateGroupDescription updates the description of a group.
func (adminSvc *AdminSvc) UpdateGroupDescription(groupName, newDescription string) (*admin.Group, error) {
	group, err := adminSvc.svc.Groups.Update(groupName, &admin.Group{
		Description: newDescription,
	}).Do()

	if err != nil {
		log.Logger.Errorf("UpdateGroupDescription: failed to update group description: %v", err)
		return nil, err
	}

	return group, nil
}

// DeleteGroup deletes a group with the specified name.
func (adminSvc *AdminSvc) DeleteGroup(groupName string) error {
	err := adminSvc.svc.Groups.Delete(groupName).Do()
	if err != nil {
		return fmt.Errorf("DeleteGroup: failed to delete group: %w", err)
	}

	return nil
}
