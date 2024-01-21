package iam_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/api/cloudresourcemanager/v1"

	"gitea/pcp-inariam/inariam/pkgs/cloud/gcp"
	"gitea/pcp-inariam/inariam/pkgs/cloud/gcp/iam"
)

const (
	CredentialsPath  = "C:\\Users\\MohamedAliWACHANI\\Downloads\\Gcp.json"
	ProjectId        = "inarium-dev"
	RoleId           = "testRoleID2"
	RoleName         = "test_gcp_test"
	RoleTitle        = "TestRoleTitle"
	UpdatedRoleTitle = "TestRole(2)"
	RoleDesc         = "Test Role Description"
	UpdatedRoleDesc  = "Test Role Description (2)"
	DisplayName      = "TestDisplayName"
	Name             = "Ironbyte"
	Email            = Name + "@" + ProjectId + ".iam.gserviceaccount.com"
)

var (
	PermList        = []string{"accessapproval.requests.approve"}
	UpdatedPermList = []string{"accessapproval.requests.approve", "accessapproval.requests.dismiss"}
	RoleX           = iam.NewRole{
		Name:        RoleName,
		Title:       RoleTitle,
		Description: RoleDesc,
		Stage:       nil,
		Permissions: PermList,
	}

	TestingPolicy = &cloudresourcemanager.Policy{
		Bindings: []*cloudresourcemanager.Binding{
			{
				Role:    "roles/editor",                       // Replace with the desired role
				Members: []string{"user:example@example.com"}, // Replace with member email
			},
		},
	}
)

func TestIamGCP(t *testing.T) {

	gSession, err := gcp.LoadCredentialsFromFile(CredentialsPath, ProjectId)
	assert.NoError(t, err)
	err = gSession.OpenIamService()
	assert.NoError(t, err)

	t.Run("RoleExists", func(t *testing.T) {
		answer, err := gSession.IamGCPService.GetIamRole(RoleName)
		fmt.Println(answer)
		assert.NoError(t, err)
	})

	t.Run("CreateRole", func(t *testing.T) {
		role, err := gSession.IamGCPService.CreateIamRole(RoleX)
		fmt.Println(role)
		assert.NoError(t, err)
	})

	t.Run("UpdateRole", func(t *testing.T) {
		role, err := gSession.IamGCPService.UpdateIamRole(RoleName, UpdatedRoleTitle, UpdatedRoleDesc, UpdatedPermList)
		fmt.Println(role)
		assert.NoError(t, err)
	})

	t.Run("GetRole", func(t *testing.T) {
		_, err := gSession.IamGCPService.GetIamRole(RoleName)
		assert.NoError(t, err)
	})

	t.Run("DeleteRole", func(t *testing.T) {
		err := gSession.IamGCPService.DeleteIamRole(RoleId)
		assert.NoError(t, err)
	})

	t.Run("ListRoles", func(t *testing.T) {
		_, err := gSession.IamGCPService.ListIamRoles()
		assert.NoError(t, err)
	})

	t.Run("CreateServiceAccount", func(t *testing.T) {
		account, err := gSession.IamGCPService.CreateIamServiceAccount(ProjectId, DisplayName, Name)
		fmt.Println(account.Email)
		assert.NoError(t, err)
	})

	t.Run(("CheckIamServiceAccountExists"), func(t *testing.T) {
		asnwer, err := gSession.IamGCPService.GetIamRole(Email)
		fmt.Println(asnwer)
		assert.NoError(t, err)
	})

	t.Run("ListIamServiceAccounts", func(t *testing.T) {
		_, err := gSession.IamGCPService.ListIamServiceAccounts()
		assert.NoError(t, err)
	})

	t.Run("DeleteIamServiceAccount", func(t *testing.T) {
		err := gSession.IamGCPService.DeleteIamServiceAccount(Email)
		assert.NoError(t, err)
	})

	t.Run("EnableServiceAccount", func(t *testing.T) {
		err := gSession.IamGCPService.EnableIamServiceAccount(Email)
		assert.NoError(t, err)
	})

	t.Run("DisableServiceAccount", func(t *testing.T) {
		err := gSession.IamGCPService.DisableIamServiceAccount(Email)
		assert.NoError(t, err)
	})

	err = gSession.OpenCrmService()
	assert.NoError(t, err)

	t.Run("SetPolicy", func(t *testing.T) {
		err := gSession.CrmGCPService.SetPolicy(TestingPolicy)
		assert.NoError(t, err)
	})

	t.Run("GetPolicy", func(t *testing.T) {
		_, err := gSession.CrmGCPService.GetIamPolicy()
		assert.NoError(t, err)
	})

	t.Run("DeletePolicy", func(t *testing.T) {
		err := gSession.CrmGCPService.DeletePolicy()
		assert.NoError(t, err)
	})
}
