package iam_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"

	"gitea/pcp-inariam/inariam/pkgs/cloud/aws"
)

const (
	ROLE_NAME      = "testRole"
	USERNAME       = "ironbyte"
	NEW_USERNAME   = "ironbyte2"
	GROUP_NAME     = "testGroup"
	NEW_GROUP_NAME = "testGroup2"
	ARN            = "arn:aws:iam::782390994097:policy/TestPolicy"
	POLICY_NAME    = "TestPolicy"
	TRUST_POLICY   = `{
	  "Version": "2012-10-17",
	  "Statement": [
	    {
	      "Effect": "Allow",
	      "Principal": {"Service": "ec2.amazonaws.com"},
	      "Action": "sts:AssumeRole"
	    }
	  ]
	}`
)

func getTestCredentials() aws.Credentials {
	if err := godotenv.Load("/home/askee/.inariam/.env"); err != nil {
		panic(err)

	}
	return aws.Credentials{
		AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		Region:          os.Getenv("AWS_DEFAULT_REGION"),
	}
}

func TestAWSSession(t *testing.T) {

	credentials := getTestCredentials()
	fmt.Printf("%#v", credentials)

	awsSess, err := aws.OpenSession(&credentials)

	assert.NoError(t, err)

	awsSess.OpenIamService()

	t.Run("CheckIfRoleExists", func(t *testing.T) {
		_, err := awsSess.IamSvc.GetIamRole(ROLE_NAME)
		assert.NoError(t, err)
	})

	t.Run("CreateIAMRole", func(t *testing.T) {
		newIamRole, err := awsSess.IamSvc.CreateIAMRole(ROLE_NAME, TRUST_POLICY)
		fmt.Println(newIamRole)
		assert.NoError(t, err)
	})

	t.Run("ModifyIAMRoleTrustPolicy", func(t *testing.T) {
		_, err := awsSess.IamSvc.ModifyIAMRoleTrustPolicy(ROLE_NAME, TRUST_POLICY)
		assert.NoError(t, err)
	})

	t.Run("DeleteIAMRole", func(t *testing.T) {
		err := awsSess.IamSvc.DeleteIAMRole(ROLE_NAME)
		assert.NoError(t, err)
	})
	t.Run("CheckIfIamUserExists", func(t *testing.T) {
		_, err := awsSess.IamSvc.GetIamUser(USERNAME)
		assert.NoError(t, err)
	})

	t.Run("CreateIamUser", func(t *testing.T) {
		_, err := awsSess.IamSvc.CreateIamUser(USERNAME)
		assert.NoError(t, err)
	})

	t.Run("UpdateIamUser", func(t *testing.T) {
		_, err := awsSess.IamSvc.UpdateIamUser(USERNAME, NEW_USERNAME)
		assert.NoError(t, err)
	})

	t.Run("DeleteIamUser", func(t *testing.T) {
		err := awsSess.IamSvc.DeleteIamUser(NEW_USERNAME)
		assert.NoError(t, err)
	})

	t.Run("ListIamUsers", func(t *testing.T) {
		_, err := awsSess.IamSvc.ListIamUsers()
		assert.NoError(t, err)
	})

	t.Run("CheckIfGroupExists", func(t *testing.T) {
		_, err := awsSess.IamSvc.CheckIfGroupExists(GROUP_NAME)
		assert.NoError(t, err)
	})

	t.Run("CreateIamGroup", func(t *testing.T) {
		_, err := awsSess.IamSvc.CreateIamGroup(GROUP_NAME)
		assert.NoError(t, err)
	})

	t.Run("DeleteIamGroup", func(t *testing.T) {
		err := awsSess.IamSvc.DeleteIamGroup(GROUP_NAME)
		assert.NoError(t, err)
	})

	t.Run("ListIamGroups", func(t *testing.T) {
		_, err := awsSess.IamSvc.ListIamGroups()
		assert.NoError(t, err)
	})

	t.Run("UpdateIamGroup", func(t *testing.T) {
		_, err := awsSess.IamSvc.UpdateIamGroup(GROUP_NAME, NEW_GROUP_NAME)
		assert.NoError(t, err)
	})

	t.Run("CheckIfPolicyExists", func(t *testing.T) {
		_, err := awsSess.IamSvc.GetIamPolicy(ARN)
		assert.NoError(t, err)
	})

	t.Run("DeletePolicy", func(t *testing.T) {
		err := awsSess.IamSvc.DeleteIamPolicy(ARN)
		assert.NoError(t, err)
	})
}
