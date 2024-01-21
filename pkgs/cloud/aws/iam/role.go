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
	ErrIamRoleExists    = "IAM role already exists"
	ErrIamRoleNotExists = "IAM role does not exist"
)

// GetIamRole checks if an IAM role exists and returns details if found.
func (IamSvc *Svc) GetIamRole(roleName string) (*iam.Role, error) {
	getRoleInput := &iam.GetRoleInput{
		RoleName: aws.String(roleName),
	}

	role, err := IamSvc.svc.GetRole(getRoleInput)
	if err != nil {
		var iamErr awserr.Error
		if errors.As(err, &iamErr) {
			if iamErr.Code() == iam.ErrCodeNoSuchEntityException {
				return nil, nil
			}
		}

		return nil, fmt.Errorf("CheckIfIamRoleExists: %w", err)
	}

	return role.Role, nil
}

// CreateIAMRole creates an IAM role and returns role ARN.
func (IamSvc *Svc) CreateIAMRole(roleName string, trustPolicy string) (*iam.Role, error) {
	roleDetails, _ := IamSvc.GetIamRole(roleName)
	if roleDetails != nil {
		return nil, errors.New(ErrIamRoleExists)
	}

	createRoleInput := &iam.CreateRoleInput{
		RoleName:                 aws.String(roleName),
		AssumeRolePolicyDocument: aws.String(trustPolicy),
	}
	createRoleOutput, err := IamSvc.svc.CreateRole(createRoleInput)
	if err != nil {
		return nil, fmt.Errorf("CreateIAMRole: %w", err)
	}

	log.Logger.Infof("IAM role '%s' created successfully\n", roleName)

	return createRoleOutput.Role, nil
}

// ModifyIAMRoleTrustPolicy modifies an IAM role trust policy and returns role details.
func (IamSvc *Svc) ModifyIAMRoleTrustPolicy(roleName string, newTrustPolicy string) (*iam.Role, error) {
	roleDetails, _ := IamSvc.GetIamRole(roleName)

	if roleDetails == nil {
		return nil, errors.New(ErrIamRoleNotExists)
	}

	updateAssumeRolePolicyInput := &iam.UpdateAssumeRolePolicyInput{

		PolicyDocument: aws.String(newTrustPolicy),
		RoleName:       aws.String(roleName),
	}
	_, err := IamSvc.svc.UpdateAssumeRolePolicy(updateAssumeRolePolicyInput)

	if err != nil {
		return nil, fmt.Errorf("ModifyIAMRoleTrustPolicy: %w", err)
	}

	log.Logger.Infof("IAM role '%s' trust policy modified successfully\n", roleName)
	return roleDetails, nil
}

// DeleteIAMRole deletes an IAM role.
func (IamSvc *Svc) DeleteIAMRole(roleName string) error {
	roleDetails, _ := IamSvc.GetIamRole(roleName)
	if roleDetails == nil {
		return errors.New(ErrIamRoleNotExists)
	}

	_, err := IamSvc.svc.DeleteRole(&iam.DeleteRoleInput{
		RoleName: aws.String(roleName),
	})
	if err != nil {
		return fmt.Errorf("DeleteIAMRole: %w", err)
	}

	log.Logger.Infof("IAM role '%s' deleted successfully\n", roleName)
	return nil
}

// ListIamRoles lists all IAM roles.
func (IamSvc *Svc) ListIamRoles() ([]*iam.Role, error) {
	listRolesInput := &iam.ListRolesInput{}

	result, err := IamSvc.svc.ListRoles(listRolesInput)
	if err != nil {
		return nil, fmt.Errorf("ListRoles: %w", err)
	}

	return result.Roles, nil
}
